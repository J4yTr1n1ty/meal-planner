package meals

import (
	"encoding/json"
	"net/http"

	"github.com/J4yTr1n1ty/meal-planner/pkg/boot"
	"github.com/J4yTr1n1ty/meal-planner/pkg/models"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) GetMeals(w http.ResponseWriter, r *http.Request) {
	meals := []models.Meal{}
	boot.DB.Find(&meals)

	names := []string{}
	for _, meal := range meals {
		names = append(names, meal.Name)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(names)
}
