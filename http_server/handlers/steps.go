package handlers

import (
	"cookbook/database/models"
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NewStepsHandler(app *fiber.App, handler *Handler) *fiber.App {
	router := app.Group("recepie/:recepie_id/step")

	router.Get("/", handler.recepie, handler.getAllSteps)
	router.Get("/:id", handler.recepie, handler.getStep)
	router.Post("/", handler.user, handler.recepie, handler.createStep)
	router.Put("/:id", handler.user, handler.recepie, handler.updateStep)
	router.Delete("/:id", handler.user, handler.recepie, handler.deleteStep)
	router.Patch("/", handler.user, handler.recepie, handler.changeStepOrder)

	return app
}

// @Summary		Get all Steps
// @Tags			Recepie/Step
// @Description	get all Recepie Steps
// @Accept			json
// @Produce		json
// @Param			offset		query		int	false	"offset"
// @Param			limit		query		int	false	"limit"
// @Param			recepie_id	path		int	true	"Recepie ID"
// @Success		200			{object}	[]models.Step
// @Failure		400			{object}	string
// @Router			/recepie/{recepie_id}/step [get]
func (h *Handler) getAllSteps(ctx *fiber.Ctx) error {
	var schema = new(OffsetLimitRequest)
	var recepie = ctx.Locals(RECEPIE_KEY).(*models.Recepie)

	if err := h.parseQueryAndValidate(ctx, schema); err != nil {
		return err
	}
	if schema.Limit == 0 {
		schema.Limit = 10
	}

	var steps = make([]*models.Step, 0)

	h.db.Offset(int(schema.Offset)).
		Limit(int(schema.Limit)).
		Where("recepie_id = $1", recepie.ID).
		Order("step_order DESC").
		Find(&steps)

	return ctx.JSON(steps)
}

// @Summary		Get Recepie Step
// @Tags			Recepie/Step
// @Description	get Recepie Step
// @Accept			json
// @Produce		json
// @Param			recepie_id	path		int	true	"Recepie ID"
// @Param			step_id		path		int	true	"Recepie ID"
// @Success		200			{object}	models.Step
// @Failure		400			{object}	string
// @Failure		404			{object}	string
// @Router			/recepie/{recepie_id}/step/{step_id} [get]
func (h *Handler) getStep(ctx *fiber.Ctx) error {
	var recepie = ctx.Locals(RECEPIE_KEY).(*models.Recepie)

	id, err := h.parsePathID(ctx, "id")
	if err != nil {
		return err
	}
	var step = new(models.Step)

	if err := h.db.Where("recepie_id = $1", recepie.ID).First(step, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	return ctx.JSON(step)
}

// @Summary		Create Recepie Step
// @Tags			Recepie/Step
// @Description	get Recepie Step
// @Accept			json
// @Produce		json
// @Param			recepie_id	path		int			true	"Recepie ID"
// @Param			step		body		createStep	true	"Recepie ID"
// @Success		200			{object}	models.Step
// @Failure		400			{object}	string
// @Failure		401			{object}	string
// @Failure		404			{object}	string
// @Router			/recepie/{recepie_id}/step [post]
func (h *Handler) createStep(ctx *fiber.Ctx) error {
	var schema = new(createStep)
	var user = ctx.Locals(USER_KEY).(*models.User)
	var recepie = ctx.Locals(RECEPIE_KEY).(*models.Recepie)
	if user.ID != recepie.AuthorID || user.Role.Name != "admin" {
		return fiber.NewError(fiber.StatusForbidden, "Add recepie steps can only author or admin")
	}

	if err := h.parseBodyAndValidate(ctx, schema); err != nil {
		return err
	}

	step := models.Step{
		Name:        schema.Name,
		Description: schema.Description,
		TimeToCook:  schema.TimeToCook,
		StepOrder:   schema.StepOrder,
		RecepieID:   recepie.ID,
	}

	h.db.Create(&step)

	return ctx.JSON(step)
}

// @Summary		Update Recepie Step
// @Tags			Recepie/Step
// @Description	update Recepie Step
// @Accept			json
// @Produce		json
// @Param			recepie_id	path		int			true	"Recepie ID"
// @Param			step_id		path		int			true	"Step ID"
// @Param			step		body		updateStep	true	"Recepie ID"
// @Success		200			{object}	models.Step
// @Failure		400			{object}	string
// @Failure		401			{object}	string
// @Failure		403			{object}	string
// @Failure		404			{object}	string
// @Router			/recepie/{recepie_id}/step/{step_id} [put]
func (h *Handler) updateStep(ctx *fiber.Ctx) error {
	var user = ctx.Locals(USER_KEY).(*models.User)
	var recepie = ctx.Locals(RECEPIE_KEY).(*models.Recepie)

	if user.ID != recepie.AuthorID || user.Role.Name != "admin" {
		return fiber.NewError(fiber.StatusForbidden, "Update recepie steps can only author or admin")
	}

	id, err := h.parsePathID(ctx, "id")
	if err != nil {
		return err
	}

	var schema = new(updateStep)
	if err := h.parseBodyAndValidate(ctx, schema); err != nil {
		return err
	}

	var step = new(models.Step)
	if err := h.db.Where("steps.recepie_id", recepie.ID).First(step, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	step.Name = schema.Name
	step.TimeToCook = schema.TimeToCook
	step.Description = schema.Description
	step.StepOrder = schema.StepOrder

	h.db.Save(&step)

	return ctx.JSON(step)
}

// @Summary		Delete Recepie Step
// @Tags			Recepie/Step
// @Description	delete Recepie Step
// @Accept			json
// @Produce		json
// @Param			recepie_id	path		int	true	"Recepie ID"
// @Param			step_id		path		int	true	"Step ID"
// @Success		204			{object}	nil
// @Failure		400			{object}	string
// @Failure		401			{object}	string
// @Failure		403			{object}	string
// @Failure		404			{object}	string
// @Router			/recepie/{recepie_id}/step/{step_id} [delete]
func (h *Handler) deleteStep(ctx *fiber.Ctx) error {
	var user = ctx.Locals(USER_KEY).(*models.User)
	var recepie = ctx.Locals(RECEPIE_KEY).(*models.Recepie)
	if user.ID != recepie.AuthorID || user.Role.Name != "admin" {
		return fiber.NewError(fiber.StatusForbidden, "Delete recepie steps can only author or admin")
	}

	id, err := h.parsePathID(ctx, "id")
	if err != nil {
		return err
	}

	var step = new(models.Step)

	if err := h.db.Where("steps.recepie_id", recepie.ID).First(step, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	h.db.Delete(step)

	return ctx.SendStatus(fiber.StatusNoContent)
}

// @Summary		Change Recepie Step order
// @Tags			Recepie/Step
// @Description	Change order
// @Accept			json
// @Produce		json
// @Param			recepie_id	path		int				true	"Recepie ID"
// @Param			new_order	body		changeStepOrder	true	"New Order"
// @Success		204			{object}	nil
// @Failure		400			{object}	string
// @Failure		401			{object}	string
// @Failure		403			{object}	string
// @Failure		404			{object}	string
// @Router			/recepie/{recepie_id}/step [put]
func (h *Handler) changeStepOrder(ctx *fiber.Ctx) error {
	var user = ctx.Locals(USER_KEY).(*models.User)
	var recepie = ctx.Locals(RECEPIE_KEY).(*models.Recepie)

	if user.ID != recepie.AuthorID || user.Role.Name != "admin" {
		return fiber.NewError(fiber.StatusForbidden, "Change recepie steps order can only author or admin")
	}

	var schema = make([]changeStepOrder, 0)

	if err := ctx.BodyParser(&schema); err != nil {
		return err
	}

	var steps = []*models.Step{}

	var step_ids = make([]uint, len(schema))
	for _, v := range schema {
		step_ids = append(step_ids, v.StepID)
	}

	if err := h.db.Where("steps.recepie_id", recepie.ID).Where("id IN ?", step_ids).Find(&steps).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	for _, req := range schema {
		for _, step := range steps {
			if step.ID == req.StepID {
				step.StepOrder = req.StepOrder
			}
		}
	}

	h.db.Save(&steps)

	return ctx.JSON(steps)
}
