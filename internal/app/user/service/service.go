package service

import (
	"context"
	"time"

	"github.com/artfoxe6/quick-gin/internal/app/core/apperr"
	"github.com/artfoxe6/quick-gin/internal/app/core/repository/builder"
	"github.com/artfoxe6/quick-gin/internal/app/core/request"
	"github.com/artfoxe6/quick-gin/internal/app/user/dto"
	"github.com/artfoxe6/quick-gin/internal/app/user/model"
	"github.com/artfoxe6/quick-gin/internal/pkg/kit"
	"github.com/artfoxe6/quick-gin/internal/pkg/mailer"
	"github.com/artfoxe6/quick-gin/internal/pkg/token"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	FindOne(map[string]any, ...*builder.Builder) *model.User
	Create(*model.User) error
	Update(*model.User) error
	Delete(uint) error
	Get(uint, ...*builder.Builder) (*model.User, error)
	ListWithCount(int, int, ...*builder.Builder) ([]model.User, int64, error)
	GetByEmail(string) *model.User
	Count(...*builder.Builder) int64
}

type CodeRepository interface {
	FindOne(map[string]any, ...*builder.Builder) *model.Code
	Create(*model.Code) error
}

type UserService interface {
	Create(r *dto.UserUpsert) error
	Update(r *dto.UserUpsert) error
	Delete(id uint) error
	Detail(id uint) (any, error)
	List(r *request.NormalSearch) (any, int64, error)
	Login(data *dto.UserLogin) (string, error)
	SuperUserToken(email, password string) (string, error)
	Register(data *dto.UserCreate) (string, error)
	SendCode(data *dto.Code) error
	UpdatePassword(data *dto.UpdatePassword) error
	GetByID(id uint) (*model.User, error)
}

type userService struct {
	userRepo UserRepository
	codeRepo CodeRepository
}

func NewUserService(userRepo UserRepository, codeRepo CodeRepository) UserService {
	return &userService{userRepo: userRepo, codeRepo: codeRepo}
}

func (s *userService) Create(r *dto.UserUpsert) error {
	user := model.User{
		Avatar:   *r.Avatar,
		Email:    *r.Email,
		Name:     *r.Name,
		Password: *r.Password,
		Role:     *r.Role,
	}
	if one := s.userRepo.FindOne(map[string]any{"email": user.Email}); one.ID != 0 {
		return apperr.BadRequest("email exists")
	}
	if err := s.userRepo.Create(&user); err != nil {
		return apperr.Internal(err)
	}
	return nil
}

func (s *userService) Update(r *dto.UserUpsert) error {
	user, err := s.userRepo.Get(*r.Id)
	if err != nil {
		return apperr.Internal(err)
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
		hashPassword, hashErr := bcrypt.GenerateFromPassword([]byte(*r.Password), bcrypt.DefaultCost)
		if hashErr != nil {
			return apperr.Internal(hashErr)
		}
		user.Password = string(hashPassword)
	}
	if one := s.userRepo.FindOne(map[string]any{"email": user.Email}); one.ID != 0 && one.ID != user.ID {
		return apperr.BadRequest("email exists")
	}
	if err = s.userRepo.Update(user); err != nil {
		return apperr.Internal(err)
	}
	return nil
}

func (s *userService) Delete(id uint) error {
	if err := s.userRepo.Delete(id); err != nil {
		return apperr.Internal(err)
	}
	return nil
}

func (s *userService) Detail(id uint) (any, error) {
	user, err := s.userRepo.Get(id)
	if err != nil {
		return nil, apperr.Internal(err)
	}
	return user.ToMap(), nil
}

func (s *userService) List(r *request.NormalSearch) (any, int64, error) {
	b := builder.New()
	if r.Keyword != nil {
		b.Where("name like ? or email like ?", "%"+*r.Keyword+"%", "%"+*r.Keyword+"%")
	}
	orderSet := map[int]string{
		0: "id desc",
		1: "id asc",
	}
	b.Order(orderSet[r.Sort])
	users, total, err := s.userRepo.ListWithCount(r.Offset(), r.Limit, b)
	if err != nil {
		return nil, 0, apperr.Internal(err)
	}
	list := make([]map[string]any, 0, len(users))
	for _, v := range users {
		list = append(list, v.ToMap())
	}
	return list, total, nil
}

func (s *userService) Login(data *dto.UserLogin) (string, error) {
	user := s.userRepo.GetByEmail(data.Email)
	if user == nil {
		return "", apperr.BadRequest("email not exists")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		return "", apperr.BadRequest("password is wrong")
	}
	if data.Role != "" && user.Role != data.Role {
		user.Role = data.Role
	}
	tokenStr, err := token.Generate(user.TokenData())
	if err != nil {
		return "", apperr.Internal(err)
	}
	return tokenStr, nil
}

func (s *userService) SuperUserToken(email, password string) (string, error) {
	user := s.userRepo.GetByEmail(email)
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", apperr.Internal(err)
	}
	if user == nil {
		user = &model.User{
			Name:     "Admin",
			Password: string(hashPassword),
			Email:    email,
		}
		if err := s.userRepo.Create(user); err != nil {
			return "", apperr.Internal(err)
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		user.Password = string(hashPassword)
		if updateErr := s.userRepo.Update(user); updateErr != nil {
			return "", apperr.Internal(updateErr)
		}
	}
	tokenStr, err := token.Generate(map[string]any{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  "admin",
	})
	if err != nil {
		return "", apperr.Internal(err)
	}
	return tokenStr, nil
}

func (s *userService) Register(data *dto.UserCreate) (string, error) {
	if data.Code != "" && s.codeRepo != nil {
		code := s.codeRepo.FindOne(map[string]any{"email": data.Email})
		if code.Code != data.Code {
			return "", apperr.BadRequest("code is wrong")
		}
	}
	user := s.userRepo.GetByEmail(data.Email)
	if user != nil {
		return "", apperr.BadRequest("email exists")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", apperr.Internal(err)
	}
	user = &model.User{
		Name:     data.Name,
		Password: string(hashPassword),
		Email:    data.Email,
		Avatar:   data.Avatar,
	}
	if err := s.userRepo.Create(user); err != nil {
		return "", apperr.Internal(err)
	}

	tokenStr, err := token.Generate(user.TokenData())
	if err != nil {
		return "", apperr.Internal(err)
	}
	return tokenStr, nil
}

func (s *userService) SendCode(data *dto.Code) error {
	if data.Type == 1 {
		user := s.userRepo.FindOne(map[string]any{"email": data.Email})
		if user != nil {
			return apperr.BadRequest("user not found")
		}
	}
	if data.Email == "" {
		return nil
	}
	b := builder.New()
	b.Eq("email", data.Email).Gt("created_at", time.Now().Format("20060102"))
	if todayCount := s.userRepo.Count(b); todayCount > 20 {
		return apperr.BadRequest("The number of sending times for the day has reached the limit")
	}
	code := kit.GenCode(6)
	if err := mailer.New(mailer.Template["code"], map[string]any{"code": code}).SendTo(context.Background(), "", data.Email); err != nil {
		return apperr.Internal(err)
	}
	if s.codeRepo == nil {
		return apperr.Internal(nil)
	}
	if err := s.codeRepo.Create(&model.Code{
		Email: data.Email,
		Type:  data.Type,
		Code:  code,
	}); err != nil {
		return apperr.Internal(err)
	}
	return nil
}

func (s *userService) UpdatePassword(data *dto.UpdatePassword) error {
	if data.Code != "" && s.codeRepo != nil {
		code := s.codeRepo.FindOne(map[string]any{"email": data.Email})
		if code.Code != data.Code {
			return apperr.BadRequest("code is wrong")
		}
	}
	user := s.userRepo.FindOne(map[string]any{"email": data.Email})
	if user == nil {
		return apperr.BadRequest("user not found")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return apperr.Internal(err)
	}
	user.Password = string(hashPassword)
	if err := s.userRepo.Update(user); err != nil {
		return apperr.Internal(err)
	}
	return nil
}

func (s *userService) GetByID(id uint) (*model.User, error) {
	user, err := s.userRepo.Get(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
