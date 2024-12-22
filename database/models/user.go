package models

import (
	"time"

	"gorm.io/gorm"
)

type (
	User struct {
		ID        uint           `gorm:"primarykey" json:"id"`
		CreatedAt time.Time      `json:"-"`
		UpdatedAt time.Time      `json:"-"`
		DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
		Username  string         `gorm:"unique" json:"username"`
		Password  string         `json:"-"`
		Sessions  []Session      `json:"-"`
		Recepies  []Recepie      `json:"-" gorm:"foreignKey:AuthorID;not null"`
		RoleID    uint           `json:"-"`
		Role      Role           `json:"-"`
	}

	Session struct {
		ID     uint   `gorm:"primarykey"`
		Sid    string `gorm:"index"`
		UserID uint
	}

	Role struct {
		ID    uint   `gorm:"primarykey" json:"-"`
		Name  string `json:"name"`
		Users []User `json:"-"`
	}
)
