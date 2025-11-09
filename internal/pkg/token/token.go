package token

import (
	"errors"
	"github.com/artfoxe6/quick-gin/internal/app/core/config"
	"github.com/artfoxe6/quick-gin/internal/app/user/model"
	"github.com/artfoxe6/quick-gin/internal/app/user/repo"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func Generate(data map[string]any) (string, error) {
	//添加过期时间
	data["exp"] = time.Now().Add(time.Hour * time.Duration(config.Jwt.Exp)).Unix()
	data["refresh"] = time.Now().Add(time.Hour * time.Duration(config.Jwt.RefreshExp)).Unix()
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(data))
	return t.SignedString([]byte(config.Jwt.Secret))
}

func Parse(token string) (map[string]interface{}, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.Jwt.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func Refresh(token string) (string, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.Jwt.Secret), nil
	}, jwt.WithoutClaimsValidation())
	if err != nil {
		return "", err
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok || !t.Valid {
		return "", errors.New("invalid")
	}
	expireTime, _ := claims.GetExpirationTime()
	if expireTime.Add(time.Hour*time.Duration(config.Jwt.RefreshExp)).Unix() < time.Now().Unix() {
		return "", errors.New("refresh expire")
	}
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(config.Jwt.Exp)).Unix()
	newT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return newT.SignedString([]byte(config.Jwt.Secret))
}

func GetUserByToken(token string) (*model.User, error) {
	data, err := Parse(token)
	if err != nil {
		return nil, err
	}
	userId := uint(data["id"].(float64))
	userRepository := repo.NewUserRepository()
	return userRepository.Get(userId)
}
