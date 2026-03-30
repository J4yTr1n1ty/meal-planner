package suggestions

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/J4yTr1n1ty/meal-planner/pkg/boot"
	"github.com/J4yTr1n1ty/meal-planner/pkg/models"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) GetSuggestions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Exclude meals already planned today or in the next 5 days
		now := time.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		fiveDaysFromNow := today.AddDate(0, 0, 5)

		var plannedMealIDs []uint
		boot.DB.Model(&models.MealPlan{}).
			Where("date >= ? AND date <= ?", today, fiveDaysFromNow).
			Pluck("meal_id", &plannedMealIDs)

		// Get all distinct meal names not in the planned set
		query := boot.DB.Model(&models.Meal{}).Select("name")
		if len(plannedMealIDs) > 0 {
			query = query.Where("id NOT IN (?)", plannedMealIDs)
		}

		var names []string
		query.Pluck("name", &names)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)

		if len(names) > 0 {
			w.Write([]byte(names[rand.Intn(len(names))]))
		}
	}
}
