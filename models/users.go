package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	User     string
	Password string
}
