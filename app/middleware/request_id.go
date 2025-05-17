package middleware

import (
	"product-service/pkg/ctxutil"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RequestIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		reqID := c.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}
		c.Locals(ctxutil.RequestIDKey, reqID)
		return c.Next()
	}
}
