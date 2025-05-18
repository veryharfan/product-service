package middleware

import (
	"product-service/app/domain"
	"product-service/app/handler/response"
	"product-service/config"
	"product-service/pkg"

	"github.com/gofiber/fiber/v2"
)

func AuthInternal(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the auth header from the request
		authHeader := c.Get(string(pkg.AuthInternalHeaderKey))
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(response.Error(domain.ErrUnauthorized))
		}
		// Check if the auth header is valid (you can implement your own logic here)
		if authHeader != cfg.InternalAuthHeader {
			return c.Status(fiber.StatusUnauthorized).JSON(response.Error(domain.ErrUnauthorized))
		}

		return c.Next()
	}
}
