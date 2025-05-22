package handler

import (
	"log/slog"
	"product-service/app/domain"
	"product-service/app/handler/response"
	"product-service/pkg/ctxutil"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type productWriteHandler struct {
	productUsecase domain.ProductWriteUsecase
	validator      *validator.Validate
}

func NewProductWriteHandler(productUsecase domain.ProductWriteUsecase, validator *validator.Validate) *productWriteHandler {
	return &productWriteHandler{productUsecase, validator}
}

func (h *productWriteHandler) Create(c *fiber.Ctx) error {
	var product domain.CreateProductRequest
	if err := c.BodyParser(&product); err != nil {
		slog.ErrorContext(c.Context(), "[productWriteHandler] Create", "body", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(domain.ErrBadRequest))
	}

	if err := h.validator.Struct(product); err != nil {
		slog.ErrorContext(c.Context(), "[productWriteHandler] Create", "validation", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(domain.ErrBadRequest))
	}

	shopID, err := ctxutil.GetShopIDCtx(c.Context())
	if err != nil {
		slog.ErrorContext(c.Context(), "[productHandler] Create", "getShopIDCtx", err)
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error(domain.ErrUnauthorized))
	}

	res, err := h.productUsecase.Create(c.Context(), shopID, &product)
	if err != nil {
		slog.ErrorContext(c.Context(), "[productWriteHandler] Create", "usecase", err)
		status, response := response.FromError(err)
		return c.Status(status).JSON(response)
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res))
}

func (h *productWriteHandler) Update(c *fiber.Ctx) error {
	idstr := c.Params("id")
	if idstr == "" {
		slog.ErrorContext(c.Context(), "[productWriteHandler] Update", "params", "product ID is empty")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(domain.ErrBadRequest))
	}

	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil || id <= 0 {
		slog.ErrorContext(c.Context(), "[productWriteHandler] Update", "params:"+idstr, err)
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(domain.ErrBadRequest))
	}

	var product domain.UpdateProductRequest
	if err := c.BodyParser(&product); err != nil {
		slog.ErrorContext(c.Context(), "[productWriteHandler] Update", "body", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(domain.ErrBadRequest))
	}

	if err := h.validator.Struct(product); err != nil {
		slog.ErrorContext(c.Context(), "[productWriteHandler] Update", "validation", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(domain.ErrBadRequest))
	}

	res, err := h.productUsecase.Update(c.Context(), id, &product)
	if err != nil {
		slog.ErrorContext(c.Context(), "[productWriteHandler] Update", "usecase", err)
		status, response := response.FromError(err)
		return c.Status(status).JSON(response)
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res))
}

func (h *productWriteHandler) SetActiveStatus(c *fiber.Ctx) error {
	idstr := c.Params("id")
	if idstr == "" {
		slog.ErrorContext(c.Context(), "[productWriteHandler] SetActiveStatus", "params", "product ID is empty")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(domain.ErrBadRequest))
	}

	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil || id <= 0 {
		slog.ErrorContext(c.Context(), "[productWriteHandler] SetActiveStatus", "params:"+idstr, err)
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(domain.ErrBadRequest))
	}

	var req domain.SetActiveStatusRequest
	if err := c.BodyParser(&req); err != nil {
		slog.ErrorContext(c.Context(), "[productWriteHandler] SetActiveStatus", "body", err)
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(domain.ErrBadRequest))
	}

	err = h.productUsecase.SetActiveStatus(c.Context(), id, req.Active)
	if err != nil {
		slog.ErrorContext(c.Context(), "[productWriteHandler] SetActiveStatus", "usecase", err)
		status, response := response.FromError(err)
		return c.Status(status).JSON(response)
	}

	return c.Status(fiber.StatusOK).JSON(response.Success[any](nil))
}
