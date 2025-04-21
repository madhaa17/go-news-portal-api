package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"news-app/config"
	"news-app/internal/adapter/repository"
	"news-app/internal/core/domain/entity"
	"news-app/lib/auth"
	"news-app/lib/conv"

	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
)

var (
	err  error
	code string
)

type AuthService interface {
	GetUserByEmail(ctx context.Context, req entity.LoginRequest) (*entity.AccessToken, error)
}

type authService struct {
	authRepository repository.AuthRepository
	cfg            *config.Config
	jwtToken       auth.Jwt
}

// GetUserByEmail implements AuthService.
func (a *authService) GetUserByEmail(ctx context.Context, req entity.LoginRequest) (*entity.AccessToken, error) {
	result, err := a.authRepository.GetUserByEmail(ctx, req)
	if err != nil {
		code = "[SERVICE] GetUserByEmail - 1"
		log.Errorw(code, err)
		return nil, err
	}

	if checkPass := conv.CheckPasswordHash(req.Password, result.Password); !checkPass {
		code = "[SERVICE] GetUserByEmail - 2"
		err = errors.New("email or password is incorrect")
		log.Errorw(code, err)
		return nil, err
	}

	jwtData := &entity.JwtData{
		UserID: float64(result.ID),
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
			ID:        fmt.Sprint(result.ID),
		},
	}

	accessToken, expireAt, err := a.jwtToken.GenerateToken(jwtData)
	if err != nil {
		code = "[SERVICE] GetUserByEmail - 3"
		log.Errorw(code, err)
		return nil, err
	}

	res := entity.AccessToken{
		Token:    accessToken,
		ExpireAt: expireAt,
	}

	return &res, nil
}

func NewAuthService(authRepository repository.AuthRepository, cfg *config.Config, jwtToken auth.Jwt) AuthService {
	return &authService{authRepository: authRepository, cfg: cfg, jwtToken: jwtToken}
}
