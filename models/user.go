package models

import (
	"github.com/jinzhu/gorm"
)

// User ... user model
type User struct {
	gorm.Model
	Name     string
	Email    string `gorm:"type:varchar(100);unique_index"`
	Gender   string `json:"Gender"`
	Password string `json:"Password"`
}
