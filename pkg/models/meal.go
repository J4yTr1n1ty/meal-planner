package models

import "gorm.io/gorm"

type Meal struct {
	gorm.Model
	Name string
}
