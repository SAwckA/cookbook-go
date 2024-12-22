package handlers

import (
	"cookbook/database/models"
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NewIngredientsHandler(app *fiber.App, handler *Handler) *fiber.App {
	var router = app.Group("recepie/:recepie_id/ingredient")

	router.Get("/", handler.recepie, handler.getAllIngredients)
	router.Get("/:id", handler.recepie, handler.getIngredient)
	router.Post("/", handler.user, handler.recepie, handler.createIngredient)
	router.Put("/:id", handler.user, handler.recepie, handler.updateIngredient)
	router.Delete("/:id", handler.user, handler.recepie, handler.deleteIngredient)

	return app
}

// @Summary		Get all Ingredients
// @Tags			Recepie/Ingredient
// @Description	get all Recepie Ingredients
// @Accept			json
// @Produce		json
// @Param			offset		query		int	false	"offset"
// @Param			limit		query		int	false	"limit"
// @Param			recepie_id	path		int	true	"Recepie ID"
// @Success		200			{object}	[]models.Ingredient
// @Failure		400			{object}	string
// @Router			/recepie/{recepie_id}/ingredient [get]
func (h *Handler) getAllIngredients(ctx *fiber.Ctx) error {
	var schema = new(OffsetLimitRequest)
	var recepie = ctx.Locals(RECEPIE_KEY).(*models.Recepie)

	if err := h.parseQueryAndValidate(ctx, schema); err != nil {
		return err
	}
	if schema.Limit == 0 {
		schema.Limit = 10
	}

	var ingredients = []*models.Ingredient{}

	h.db.Offset(int(schema.Offset)).Limit(int(schema.Limit)).Where("recepie_id = $1", recepie.ID).Find(&ingredients)

	return ctx.JSON(ingredients)
}

// @Summary		Get Ingredient
// @Tags			Recepie/Ingredient
// @Description	get Recepie Ingredient by ID
// @Accept			json
// @Produce		json
// @Param			recepie_id		path		int	true	"Recepie ID"
// @Param			ingredient_id	path		int	true	"Ingredient ID"
// @Success		200				{object}	models.Ingredient
// @Failure		400				{object}	string
// @Failure		404				{object}	string
// @Router			/recepie/{recepie_id}/ingredient/{ingredient_id} [get]
func (h *Handler) getIngredient(ctx *fiber.Ctx) error {
	id, err := h.parsePathID(ctx, "id")
	var recepie = ctx.Locals(RECEPIE_KEY).(*models.Recepie)

	if err != nil {
		return err
	}

	var ingredient = new(models.Ingredient)

	if err := h.db.Where("ingredients.recepie_id = $1", recepie.ID).First(ingredient, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.ErrNotFound.Code, err.Error())
	}

	return ctx.JSON(ingredient)
}

// @Summary		Create Ingredient
// @Tags			Recepie/Ingredient
// @Description	create Recepie Ingredient
// @Accept			json
// @Produce		json
// @Param			recepie_id	path		int					true	"Recepie ID"
// @Param			ingredient	body		createIngredient	true	"Ingredient"
// @Success		200			{object}	models.Ingredient
// @Failure		400			{object}	string
// @Failure		401			{object}	string
// @Failure		404			{object}	string
// @Router			/recepie/{recepie_id}/ingredient/ [post]
func (h *Handler) createIngredient(ctx *fiber.Ctx) error {
	var schema = new(createIngredient)
	var user = ctx.Locals(USER_KEY).(*models.User)
	var recepie = ctx.Locals(RECEPIE_KEY).(*models.Recepie)
	if user.ID != recepie.AuthorID || user.Role.Name != "admin" {
		return fiber.NewError(fiber.StatusForbidden, "Create recepie ingredients can only author or admin")
	}

	if err := h.parseBodyAndValidate(ctx, schema); err != nil {
		return err
	}

	var ingredient = models.Ingredient{
		Name:       schema.Name,
		Amount:     schema.Amount,
		AmountType: schema.AmountType,
		RecepieID:  recepie.ID,
	}

	if err := h.db.Create(&ingredient).Error; err != nil {
		return err
	}

	return ctx.JSON(ingredient)
}

// @Summary		Update Ingredient
// @Tags			Recepie/Ingredient
// @Description	update Recepie Ingredient
// @Accept			json
// @Produce		json
// @Param			recepie_id		path		int					true	"Recepie ID"
// @Param			ingredient_id	path		int					true	"Recepie ID"
// @Param			ingredient		body		updateIngredient	true	"Ingredient"
// @Success		200				{object}	models.Ingredient
// @Failure		400				{object}	string
// @Failure		401				{object}	string
// @Failure		403				{object}	string
// @Failure		404				{object}	string
// @Router			/recepie/{recepie_id}/ingredient/{ingredient_id} [put]
func (h *Handler) updateIngredient(ctx *fiber.Ctx) error {
	var user = ctx.Locals(USER_KEY).(*models.User)
	var recepie = ctx.Locals(RECEPIE_KEY).(*models.Recepie)
	if user.ID != recepie.AuthorID || user.Role.Name != "admin" {
		return fiber.NewError(fiber.StatusForbidden, "Update recepie ingredients can only author or admin")
	}

	id, err := h.parsePathID(ctx, "id")
	if err != nil {
		return err
	}

	var schema = new(updateIngredient)
	if err := h.parseBodyAndValidate(ctx, schema); err != nil {
		return err
	}

	var ingredient = new(models.Ingredient)

	if err := h.db.Where("ingredients.recepie_id = $1", recepie.ID).First(ingredient, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	ingredient.Name = schema.Name
	ingredient.Amount = schema.Amount
	ingredient.AmountType = schema.AmountType

	h.db.Save(&ingredient)

	return ctx.JSON(ingredient)
}

// @Summary		Delete Ingredient
// @Tags			Recepie/Ingredient
// @Description	delete Recepie Ingredient
// @Accept			json
// @Produce		json
// @Param			recepie_id		path		int	true	"Recepie ID"
// @Param			ingredient_id	path		int	true	"Recepie ID"
// @Success		200				{object}	models.Ingredient
// @Failure		400				{object}	string
// @Failure		401				{object}	string
// @Failure		403				{object}	string
// @Failure		404				{object}	string
// @Router			/recepie/{recepie_id}/ingredient/{ingredient_id} [delete]
func (h *Handler) deleteIngredient(ctx *fiber.Ctx) error {
	var user = ctx.Locals(USER_KEY).(*models.User)
	var recepie = ctx.Locals(RECEPIE_KEY).(*models.Recepie)
	if user.ID != recepie.AuthorID || user.Role.Name != "admin" {
		return fiber.NewError(fiber.StatusForbidden, "Delete recepie ingredients can only author or admin")
	}

	id, err := h.parsePathID(ctx, "id")
	if err != nil {
		return err
	}

	var ingredient = new(models.Ingredient)

	if err := h.db.Where("ingredients.recepie_id = $1", recepie.ID).First(ingredient, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.ErrNotFound.Code, err.Error())
	}

	h.db.Delete(ingredient)

	return ctx.SendStatus(fiber.StatusNoContent)
}
