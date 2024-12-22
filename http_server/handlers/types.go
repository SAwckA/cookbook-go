package handlers

// ------- Ingrediets -------
type (
	createIngredient struct {
		Name       string `json:"name" validate:"required"`
		Amount     uint   `json:"amount" validate:"required"`
		AmountType string `json:"amount_type"`
	}

	updateIngredient struct {
		createIngredient
	}
)

// ------- Recepies -------
type (
	recepieEdit struct {
		recepieCreate
	}

	recepieCreate struct {
		Name       string `json:"name" validate:"required"`
		TimeToCook uint   `json:"time_to_cook,omitempty" validate:"min=0"`
	}
)

// ------- Step -------
type (
	createStep struct {
		Name        string `json:"name" validate:"required"`
		TimeToCook  uint   `json:"time_to_cook,omitempty" validate:"min=0"`
		Description string `json:"description" validate:"required"`
		StepOrder   int    `json:"step_order,omitempty"`
	}

	updateStep struct {
		createStep
	}

	changeStepOrder struct {
		StepID    uint `json:"step_id"`
		StepOrder int  `json:"step_order"`
	}
)

// ------- Auth -------

type (
	Register struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	Login struct {
		Register
	}
)
