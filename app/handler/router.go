package handler

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRouter(app *fiber.App, handler *productReadHandler) {
	// Setup user routes
	userGroup := app.Group("/products")

	userGroup.Post("/:id", handler.GetByID)
	userGroup.Post("/", handler.GetListByQuery)
}
