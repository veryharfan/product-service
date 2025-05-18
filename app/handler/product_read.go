package handler

import (
	"log/slog"
	"product-service/app/domain"
	"product-service/app/handler/response"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type productReadHandler struct {
	productUsecase domain.ProductReadUsecase
	validator      *validator.Validate
}

func NewProductReadHandler(productUsecase domain.ProductReadUsecase, validator *validator.Validate) *productReadHandler {
	return &productReadHandler{productUsecase, validator}
}

func (h *productReadHandler) GetByID(c *fiber.Ctx) error {
	idstr := c.Params("id")
	if idstr == "" {
		slog.ErrorContext(c.Context(), "[productReadHandler] GetByID", "params", "product ID is empty")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(domain.ErrBadRequest))
	}

	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil || id <= 0 {
		slog.ErrorContext(c.Context(), "[productReadHandler] GetByID", "params:"+idstr, err)
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(domain.ErrBadRequest))
	}

	product, err := h.productUsecase.GetByID(c.Context(), id)
	if err != nil {
		slog.ErrorContext(c.Context(), "[productReadHandler] GetByID", "usecase", err)
		status, response := response.FromError(err)
		return c.Status(status).JSON(response)
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(product))
}

func (h *productReadHandler) GetListByQuery(c *fiber.Ctx) error {
	var query domain.ProductQuery
	if err := c.QueryParser(&query); err != nil {
		slog.ErrorContext(c.Context(), "[productReadHandler] GetListByQuery", "query", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(domain.ErrBadRequest))
	}

	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 {
		query.Limit = 10
	}
	if query.Limit > 20 {
		query.Limit = 20
	}
	if query.SortBy == "" || (query.SortBy != "created_at" && query.SortBy != "price") {
		query.SortBy = "created_at"
	}
	if query.SortOrder == "" || query.SortOrder != "asc" {
		query.SortOrder = "desc"
	}

	products, err := h.productUsecase.GetListByQuery(c.Context(), query)
	if err != nil {
		slog.ErrorContext(c.Context(), "[productReadHandler] GetListByQuery", "usecase", err)
		status, response := response.FromError(err)
		return c.Status(status).JSON(response)
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(products))
}
