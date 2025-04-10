package services

import (
	"errors"
	"github.com/artfoxe6/quick-gin/internal/app/models"
	"github.com/artfoxe6/quick-gin/internal/app/repositories"
	"github.com/artfoxe6/quick-gin/internal/app/repositories/builder"
	"github.com/artfoxe6/quick-gin/internal/app/request"
	"github.com/artfoxe6/quick-gin/internal/pkg/kit"
	"github.com/artfoxe6/quick-gin/internal/pkg/mailer"
	"github.com/artfoxe6/quick-gin/internal/pkg/token"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserService struct {
	repository *repositories.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		repository: repositories.NewUserRepository(),
	}
}

func (s *UserService) Create(r *request.UserUpsert) error {
	user := models.User{
		Avatar:   *r.Avatar,
		Email:    *r.Email,
		Name:     *r.Name,
		Password: *r.Password,
		Role:     *r.Role,
	}
	if one := s.repository.FindOne(map[string]any{"email": user.Email}); one.ID != 0 {
		return errors.New("email exists")
	}
	if err := s.repository.Create(&user); err != nil {
		return err
	}
	return nil
}

func (s *UserService) Update(r *request.UserUpsert) error {
	user, err := s.repository.Get(*r.Id)
	if err != nil {
		return err
	}
	if r.Name != nil {
		user.Name = *r.Name
	}
	if r.Email != nil {
		user.Email = *r.Email
	}
	if r.Role != nil {
		user.Role = *r.Role
	}
	if r.Avatar != nil {
		user.Avatar = *r.Avatar
	}
	if r.Password != nil {
		hashPassword, _ := bcrypt.GenerateFromPassword([]byte(*r.Password), bcrypt.DefaultCost)
		user.Password = string(hashPassword)
	}
	if one := s.repository.FindOne(map[string]any{"email": user.Email}); one.ID != 0 && one.ID != user.ID {
		return errors.New("email exists")
	}
	if err = s.repository.Update(user); err != nil {
		return err
	}
	return nil
}

func (s *UserService) Delete(id uint) error {
	return s.repository.Delete(id)
}
func (s *UserService) Detail(id uint) (any, error) {
	user, err := s.repository.Get(id)
	if err != nil {
		return nil, err
	}
	return user.ToMap(), nil
}
func (s *UserService) List(r *request.NormalSearch) (any, int64, error) {
	b := builder.New()
	if r.Keyword != nil {
		b.Where("name like ? or email like ?",
			"%"+*r.Keyword+"%", "%"+*r.Keyword+"%")
	}
	orderSet := map[int]string{
		0: "id desc",
		1: "id asc",
	}
	b.Order(orderSet[r.Sort])
	users, total, err := s.repository.ListWithCount(r.Offset, r.Limit, b)
	if err != nil {
		return nil, 0, err
	}
	list := make([]map[string]any, 0, len(users))
	for _, v := range users {
		list = append(list, v.ToMap())
	}
	return list, total, nil
}

func (s *UserService) Login(data *request.UserLogin) (string, error) {

	user := s.repository.GetByEmail(data.Email)
	if user == nil {
		return "", errors.New("email not exists")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		return "", errors.New("password is wrong")
	}
	if data.Role != "" && user.Role != data.Role {
		user.Role = data.Role
	}
	return token.Generate(user.TokenData())
}

func (s *UserService) SuperUserToken(email, password string) (string, error) {
	user := s.repository.GetByEmail(email)
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if user == nil {
		user = &models.User{
			Name:     "Admin",
			Password: string(hashPassword),
			Email:    email,
		}
		err := s.repository.Create(user)
		if err != nil {
			return "", err
		}
	}

	if err := bcrypt.CompareHashAndPassword(hashPassword, []byte(password)); err != nil {
		user.Password = string(hashPassword)
		_ = s.repository.Update(user)
	}
	return token.Generate(map[string]any{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  "admin",
	})
}

func (s *UserService) Register(data *request.UserCreate) (string, error) {

	if data.Code != "" {
		codeRepository := repositories.NewCodeRepository()
		code := codeRepository.FindOne(map[string]any{"email": data.Email})
		if code.Code != data.Code {
			return "", errors.New("code is wrong")
		}
	}
	user := s.repository.GetByEmail(data.Email)
	if user != nil {
		return "", errors.New("email exists")
	}

	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	user = &models.User{
		Name:     data.Name,
		Password: string(hashPassword),
		Email:    data.Email,
		Avatar:   data.Avatar,
	}
	err := s.repository.Create(user)
	if err != nil {
		return "", err
	}

	return token.Generate(user.TokenData())
}

func (s *UserService) SendCode(data *request.Code) error {
	if data.Type == 1 {
		user := s.repository.FindOne(map[string]any{"email": data.Email})
		if user != nil {
			return errors.New("user not found")
		}
	}
	if data.Email != "" {
		b := builder.New()
		b.Eq("email", data.Email).Gt("created_at", time.Now().Format("20060102"))
		if todayCount := s.repository.Count(b); todayCount > 20 {
			return errors.New("The number of sending times for the day has reached the limit")
		}
		code := kit.GenCode(6)
		resp, err := mailer.New(mailer.Template["code"], map[string]any{"code": code}).SendTo("", data.Email)
		if err != nil {
			return err
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return errors.New(resp.Body)
		}
		codeRepository := repositories.NewCodeRepository()
		return codeRepository.Create(&models.Code{
			Email: data.Email,
			Type:  data.Type,
			Code:  code,
		})
	}
	return nil
}

func (s *UserService) UpdatePassword(data *request.UpdatePassword) error {

	if data.Code != "" {
		codeRepository := repositories.NewCodeRepository()
		code := codeRepository.FindOne(map[string]any{"email": data.Email})
		if code.Code != data.Code {
			return errors.New("code is wrong")
		}
	}
	user := s.repository.FindOne(map[string]any{"email": data.Email})
	if user == nil {
		return errors.New("user not found")
	}

	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	user.Password = string(hashPassword)
	return s.repository.Update(user)

}
