package models

import "gorm.io/gorm"

type FamilyMember struct {
	gorm.Model
	Name string
}
