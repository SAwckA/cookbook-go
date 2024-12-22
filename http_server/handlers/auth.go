package handlers

import (
	"cookbook/database/models"
	_ "cookbook/docs"
	"cookbook/utils"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const USER_KEY = "user"

func NewAuthHandlers(app *fiber.App, handler *Handler) *fiber.App {
	router := app.Group("auth")

	router.Post("/register", handler.register)
	router.Post("/login", handler.login)
	router.Post("/logout", handler.user, handler.logout)

	return app
}

func (h *Handler) user(ctx *fiber.Ctx) error {
	var user = new(models.User)
	var sid = ctx.Cookies(utils.SESSION_COOKIE_KEY, "")

	if sid == "" {
		// Не делаем запрос в бд
		return fiber.NewError(fiber.StatusUnauthorized, "sid not specified")
	}
	if err := h.db.Joins("join sessions on sessions.user_id=users.id").
		Where("sessions.sid = $1", sid).Preload("Role").First(user).Error; err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid sid")
	}

	ctx.Locals(USER_KEY, user)

	return ctx.Next()
}

// @Summary		Register an User
// @Tags			Auth
// @Description	register
// @ID				register-user
// @Accept			json
// @Produce		json
// @Param			id	body		Register	true	"New user"
// @Success		200	{object}	models.Recepie
// @Failure		400	{object}	string
// @Failure		409	{object}	string
// @Router			/auth/register [post]
func (h *Handler) register(ctx *fiber.Ctx) error {
	var schema = new(Register)

	if err := h.parseBodyAndValidate(ctx, schema); err != nil {
		return err
	}

	var user = &models.User{
		Username: schema.Username,
		Password: utils.HashPassword(schema.Password),
	}

	if err := h.db.Create(user).Error; err != nil {
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			return fiber.NewError(fiber.StatusConflict, "Username already exists")
		}
		return err
	}

	return ctx.JSON(user)
}

// @Summary		Login an User
// @Tags			Auth
// @Description	Login
// @ID				login-user
// @Accept			json
// @Produce		json
// @Param			id	body		Login	true	"User credentials"
// @Success		200	{object}	models.User
// @Header			200	{string}	Set-Cookie	"sid="
// @Failure		400	{object}	string
// @Failure		404	{object}	string
// @Router			/auth/login [post]
func (h *Handler) login(ctx *fiber.Ctx) error {
	var schema = new(Login)

	if err := h.parseBodyAndValidate(ctx, schema); err != nil {
		return err
	}

	var user = new(models.User)
	if err := h.db.Where("username = $1", schema.Username).First(user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	if ok := utils.CheckPassword(schema.Password, user.Password); !ok {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	var sid = utils.CreateSID()

	var session = &models.Session{
		Sid:    sid,
		UserID: user.ID,
	}

	if err := h.db.Create(session).Error; err != nil {
		return err
	}

	ctx.Cookie(&fiber.Cookie{
		Name:  utils.SESSION_COOKIE_KEY,
		Value: sid,
	})

	return ctx.JSON(user)
}

// @Summary		Logout an User
// @Security		CookieSID
// @Tags			Auth
// @Description	Logout, delete session
// @ID				logout-user
// @Accept			json
// @Produce		json
// @Success		200	{object}	nil
// @Header			200	{string}	Set-Cookie	"sid=; Max-Age=0;"
// @Failure		400	{object}	string
// @Failure		401	{object}	string
// @Router			/auth/logout [post]
func (h *Handler) logout(ctx *fiber.Ctx) error {
	if err := h.db.Where("sid = $1", ctx.Cookies(utils.SESSION_COOKIE_KEY)).Delete(&models.Session{}).Error; err != nil {
		return err
	}

	ctx.ClearCookie(utils.SESSION_COOKIE_KEY)

	return ctx.SendStatus(fiber.StatusOK)
}
