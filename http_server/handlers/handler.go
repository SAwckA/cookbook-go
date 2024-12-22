package handlers

import (
	_ "cookbook/docs"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
	v  *validator.Validate
}

func New(db *gorm.DB) *Handler {
	var v = validator.New()
	return &Handler{
		db: db,
		v:  v,
	}
}

func (h *Handler) parseBodyAndValidate(ctx *fiber.Ctx, out interface{}) error {

	if err := ctx.BodyParser(&out); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return h.validate(out)
}

func (h *Handler) parseQueryAndValidate(ctx *fiber.Ctx, out interface{}) error {
	if err := ctx.QueryParser(out); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return h.validate(out)
}

func (h *Handler) validate(s interface{}) error {
	if err := h.v.Struct(s); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return nil
}

func (h *Handler) parsePathID(ctx *fiber.Ctx, key string) (uint, error) {
	id, err := ctx.ParamsInt(key)
	if err != nil {
		return 0, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if id <= 0 {
		return 0, fiber.NewError(fiber.StatusBadRequest, "ID must be greater than 0")
	}
	return uint(id), err
}

type OffsetLimitRequest struct {
	Offset uint `query:"offset" validate:""`
	Limit  uint `query:"limit" validate:"max=100"`
}
