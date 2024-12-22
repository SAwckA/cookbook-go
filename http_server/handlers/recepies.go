package handlers

import (
	"cookbook/database/models"
	_ "cookbook/docs"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

const RECEPIE_KEY = "recepie"

func NewRecepieHandler(app *fiber.App, handler *Handler) *fiber.App {
	router := app.Group("recepie")

	router.Get("/", handler.allRecepies)
	router.Get("/:recepie_id", handler.recepie, handler.getRecepie)
	router.Post("/", handler.user, handler.createRecepie)
	router.Put("/:recepie_id", handler.recepie, handler.user, handler.updateRecepie)
	router.Delete("/:recepie_id", handler.recepie, handler.user, handler.deleteRecepie)

	return app
}

func (h *Handler) recepie(ctx *fiber.Ctx) error {

	id, err := h.parsePathID(ctx, "recepie_id")
	if err != nil {
		return err
	}

	var recepie = new(models.Recepie)
	if err := h.db.First(&recepie, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusNotFound, "recepie not found")
	}

	ctx.Locals(RECEPIE_KEY, recepie)
	log.Debug("pass recepie: ", recepie)
	return ctx.Next()
}

//	@Summary		Show an Recepie
//	@Tags			Recepie
//	@Description	get by ID
//	@ID				get-by-int
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Recepie ID"
//	@Success		200	{object}	models.Recepie
//	@Failure		400	{object}	string
//	@Failure		404	{object}	string
//	@Router			/recepie/{id} [get]
func (h *Handler) getRecepie(ctx *fiber.Ctx) error {
	var recepie = ctx.Locals(RECEPIE_KEY)
	return ctx.JSON(recepie)
}

//	@Summary	Create an Recepie
//
//	@Security	CookieSID
//
//	@Tags		Recepie
//	@Accept		json
//	@Produce	json
//	@Param		recepie	body		recepieCreate	true	"Recepie name"
//	@Success	200		{object}	models.Recepie
//	@Failure	400		{object}	string
//	@Failure	401		{object}	string
//	@Router		/recepie [post]
func (h *Handler) createRecepie(ctx *fiber.Ctx) error {
	var user = ctx.Locals(USER_KEY).(*models.User)
	var schema = new(recepieCreate)

	if err := h.parseBodyAndValidate(ctx, schema); err != nil {
		return err
	}

	var recepie = models.Recepie{
		Name:       schema.Name,
		TimeToCook: schema.TimeToCook,
		AuthorID:   user.ID,
	}

	h.db.Create(&recepie)

	return ctx.JSON(recepie)
}

//	@Summary	Update an Recepie
//
//	@Security	CookieSID
//
//	@Tags		Recepie
//	@Accept		json
//	@Produce	json
//	@Param		recepie	body		recepieEdit	true	"Recepie name"
//	@Param		id		path		int			true	"Recepie ID"
//	@Success	200		{object}	models.Recepie
//	@Failure	400		{object}	string
//	@Failure	401		{boject}	string
//	@Failure	403		{boject}	string
//	@Router		/recepie/{id} [put]
func (h *Handler) updateRecepie(ctx *fiber.Ctx) error {

	var schema = new(recepieEdit)
	var recepie = ctx.Locals(RECEPIE_KEY).(*models.Recepie)
	var user = ctx.Locals(USER_KEY).(*models.User)

	if recepie.AuthorID != user.ID || user.Role.Name == "admin" {
		return fiber.NewError(fiber.StatusForbidden, "Update recepie can only author or admin")
	}

	if err := h.parseBodyAndValidate(ctx, schema); err != nil {
		return err
	}

	recepie.Name = schema.Name
	recepie.TimeToCook = schema.TimeToCook

	h.db.Save(recepie)

	return ctx.JSON(recepie)
}

//	@Summary		Delete an Recepie
//
//	@Security		CookieSID
//
//	@Description	delete by ID
//	@Tags			Recepie
//	@ID				delete-by-int
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Recepie ID"
//	@Success		204	{object}	nil
//	@Failure		400	{object}	string
//	@Failure		401	{boject}	string
//	@Failure		403	{boject}	string
//	@Failure		404	{object}	string
//	@Router			/recepie/{id} [delete]
func (h *Handler) deleteRecepie(ctx *fiber.Ctx) error {
	var recepie = ctx.Locals(RECEPIE_KEY).(*models.Recepie)
	var user = ctx.Locals(RECEPIE_KEY).(*models.User)

	if user.ID != recepie.AuthorID || user.Role.Name == "admin" {
		return fiber.NewError(fiber.StatusForbidden, "Delete recepie can only author or admin")
	}

	h.db.Delete(recepie)

	return ctx.SendStatus(fiber.StatusNoContent)
}

//	@Summary		Get all Recepies
//	@Tags			Recepie
//	@Description	get all Recepies
//	@Accept			json
//	@Produce		json
//	@Param			offset	query		int	false	"offset"
//	@Param			limit	query		int	false	"limit"
//	@Success		200		{object}	[]models.Recepie
//	@Failure		400		{object}	string
//	@Router			/recepie [get]
func (h *Handler) allRecepies(ctx *fiber.Ctx) error {
	var schema = new(OffsetLimitRequest)

	if err := h.parseQueryAndValidate(ctx, schema); err != nil {
		return err
	}
	if schema.Limit == 0 {
		schema.Limit = 10
	}

	var recepies = make([]*models.Recepie, 0)

	h.db.Offset(int(schema.Offset)).Limit(int(schema.Limit)).Find(&recepies)

	return ctx.JSON(recepies)
}
