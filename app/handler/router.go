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
	writeProduct := app.Group("/product-service").Use(middleware.Auth(cfg.Jwt.SecretKey))

	writeProduct.Post("/products", writeProductHandler.Create)
	writeProduct.Patch("/products/:id", writeProductHandler.SetActiveStatus)
}
