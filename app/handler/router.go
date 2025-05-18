package handler

import (
	"product-service/app/middleware"
	"product-service/config"

	"github.com/gofiber/fiber/v2"
)

func SetupRouter(app *fiber.App, readProductHandler *productReadHandler, writeProductHandler *productWriteHandler, cfg *config.Config) {
	// Setup routes
	productGroup := app.Group("/product-service")

	productGroup.Get("/products/:id", readProductHandler.GetByID)
	productGroup.Get("/products", readProductHandler.GetListByQuery)

	// write product routes
	internalGroup := app.Group("/internal/product-service")
	internalGroup.Use(middleware.AuthInternal(cfg))

	internalGroup.Post("/products", writeProductHandler.Create)
	internalGroup.Put("/products/:id", writeProductHandler.Update)
	internalGroup.Patch("/products/:id", writeProductHandler.SetActiveStatus)
}
