package handler

import (
	"news-app/internal/adapter/handler/request"
	"news-app/internal/adapter/handler/response"
	"news-app/internal/core/domain/entity"
	"news-app/internal/core/service"
	validatorLib "news-app/lib/validator"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

var (
	err         error
	code        string
	errResponse response.ErrorResponseDefault
	validate    = validator.New()
)

type AuthHandler interface {
	Login(c *fiber.Ctx) error
}

type authHandler struct {
	authService service.AuthService
}

// Login implements AuthHandler.
func (a *authHandler) Login(c *fiber.Ctx) error {
	req := request.LoginRequest{}
	res := response.SuccessAuthResponse{}

	if err := c.BodyParser(&req); err != nil {
		code = "[HANDLER] Login - 1"
		log.Errorw(code, err)
		errResponse.Meta.Status = false
		errResponse.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errResponse)
	}

	if err := validatorLib.ValidateStruct(req); err != nil {
		code = "[HANDLER] Login - 2"
		log.Errorw(code, err)
		errResponse.Meta.Status = false
		errResponse.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errResponse)
	}

	reqLogin := entity.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	result, err := a.authService.GetUserByEmail(c.Context(), reqLogin)
	if err != nil {
		code = "[HANDLER] Login - 3"
		errResponse.Meta.Status = false
		errResponse.Meta.Message = err.Error()

		if err.Error() == "email or password is incorrect" {
			return c.Status(fiber.StatusUnauthorized).JSON(errResponse)
		}

		return c.Status(fiber.StatusInternalServerError).JSON(errResponse)
	}

	res.Meta.Status = true
	res.Meta.Message = "Login success"
	res.AccessToken = result.Token
	res.ExpiredAt = result.ExpireAt

	return c.JSON(res)
}

func NewAuthHandler(authService service.AuthService) AuthHandler {
	return &authHandler{
		authService: authService,
	}
}
