package handlers

import (
	"cookbook/database/models"
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NewUserHandlers(app *fiber.App, handler *Handler) *fiber.App {
	router := app.Group("user")

	router.Get("/:id", handler.getUser)
	router.Get("/", handler.user, handler.getAllUsers)
	router.Delete("/:id", handler.user, handler.deleteUser)

	return app
}

func (h *Handler) getAllUsers(ctx *fiber.Ctx) error {
	var user = ctx.Locals(USER_KEY).(*models.User)
	if user.Role.Name != "admin" {
		return fiber.NewError(fiber.StatusForbidden, "Get all users can only admin")
	}
	var schema = new(OffsetLimitRequest)

	if err := h.parseQueryAndValidate(ctx, schema); err != nil {
		return err
	}

	if schema.Limit == 0 {
		schema.Limit = 10
	}

	var users = make([]*models.User, 0)

	h.db.Offset(int(schema.Offset)).Limit(int(schema.Limit)).Find(&users)

	return ctx.JSON(users)
}

func (h *Handler) getUser(ctx *fiber.Ctx) error {
	id, err := h.parsePathID(ctx, "id")
	if err != nil {
		return err
	}

	var user = new(models.User)

	if err := h.db.First(user, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	return ctx.JSON(user)
}

func (h *Handler) deleteUser(ctx *fiber.Ctx) error {
	id, err := h.parsePathID(ctx, "id")
	if err != nil {
		return err
	}

	var user_actor = ctx.Locals(USER_KEY).(*models.User)
	if user_actor.Role.Name != "admin" || user_actor.ID == id {
		return fiber.NewError(fiber.StatusForbidden, "Delete users can admin or self user")
	}

	var user = new(models.User)

	if err := h.db.First(user, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	h.db.Delete(user)

	return ctx.SendStatus(fiber.StatusNoContent)
}
