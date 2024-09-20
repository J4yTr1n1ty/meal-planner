package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type MealPlan struct {
	gorm.Model
	Meal           Meal `gorm:"foreignKey:MealID"`
	MealID         uint
	FamilyMember   FamilyMember `gorm:"foreignKey:FamilyMemberID"`
	FamilyMemberID uint
	Date           time.Time
}

func (m *MealPlan) GetDaysFromNow() string {
	t := m.Date
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	given := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	diff := given.Sub(today)
	days := int(diff.Hours() / 24)

	switch {
	case days == 0:
		return "Today"
	case days == 1:
		return "Tomorrow"
	case days > 1:
		return fmt.Sprintf("In %d days", days)
	default:
		return "Invalid date"
	}
}
