package models

import (
	"time"

	"gorm.io/gorm"
)

type (
	Recepie struct {
		ID          uint           `gorm:"primarykey" json:"id"`
		CreatedAt   time.Time      `json:"created_at"`
		UpdatedAt   time.Time      `json:"-"`
		DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
		Name        string         `json:"name"`
		TimeToCook  uint           `json:"time_to_cook" gorm:"default=0"`
		Ingredients []Ingredient   `json:"-"`
		Steps       []Step         `json:"-"`
		AuthorID    uint           `json:"author"`
	}

	Ingredient struct {
		ID         uint           `gorm:"primarykey" json:"id"`
		CreatedAt  time.Time      `json:"created_at"`
		UpdatedAt  time.Time      `json:"updated_at"`
		DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
		Name       string         `json:"name"`
		RecepieID  uint           `json:"-"`
		Amount     uint           `json:"amount"`
		AmountType string         `json:"amount_type"`
	}

	Step struct {
		ID          uint           `gorm:"primarykey" json:"id"`
		CreatedAt   time.Time      `json:"created_at"`
		UpdatedAt   time.Time      `json:"updated_at"`
		DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
		Name        string         `json:"name"`
		TimeToCook  uint           `json:"time_to_cook" gorm:"default:0"`
		Description string         `json:"description"`
		StepOrder   int            `json:"step_order" gorm:"default:0"`
		RecepieID   uint           `json:"-"`
	}
)
