package middleware

import (
	"strings"

	"news-app/internal/adapter/handler/response"
	"news-app/lib/auth"

	"news-app/config"

	"github.com/gofiber/fiber/v2"
)

type Middleware interface {
	CheckToken() fiber.Handler
}

type Options struct {
	authJwt auth.Jwt
}

// CheckToken implements Middleware.
func (o *Options) CheckToken() fiber.Handler {
	var errorResponse response.ErrorResponseDefault
	return func(c *fiber.Ctx) error {
		authHandler := c.Get("Authorization")
		if authHandler == "" {
			errorResponse.Meta.Status = false
			errorResponse.Meta.Message = "Authorization header is missing"
			return c.Status(fiber.StatusUnauthorized).JSON(errorResponse)
		}

		token := strings.Split(authHandler, "Bearer ")[1]
		claims, err := o.authJwt.VeryfyToken(token)
		if err != nil {
			errorResponse.Meta.Status = false
			errorResponse.Meta.Message = "Invalid token"
			return c.Status(fiber.StatusUnauthorized).JSON(errorResponse)
		}

		c.Locals("user", claims)

		return c.Next()
	}
}

func NewMiddleware(cfg *config.Config) Middleware {
	opt := new(Options)
	opt.authJwt = auth.NewJwt(cfg)

	return opt
}
