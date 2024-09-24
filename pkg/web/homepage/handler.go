package homepage

import (
	"net/http"
	"time"

	"github.com/J4yTr1n1ty/meal-planner/pkg/boot"
	"github.com/J4yTr1n1ty/meal-planner/pkg/models"
	"github.com/J4yTr1n1ty/meal-planner/pkg/web/templates"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

type EditData struct {
	ID           uint
	FamilyMember string
	Meal         string
	Date         time.Time
}

func (h *Handler) Homepage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		http.ServeFile(w, r, "static/index.html")
	}
}

func (h *Handler) AddMealPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		http.ServeFile(w, r, "static/addmeal.html")
	}
}

func (h *Handler) EditMealPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		mealPlan := models.MealPlan{}
		result := boot.DB.Joins("FamilyMember").Joins("Meal").Find(&mealPlan, id)

		if result.Error != nil {
			http.Error(w, "Unable to retrive MealPlan from Database", http.StatusBadRequest)
		}

		editData := EditData{
			ID:           mealPlan.ID,
			FamilyMember: mealPlan.FamilyMember.Name,
			Meal:         mealPlan.Meal.Name,
			Date:         mealPlan.Date,
		}

		editTemplate := templates.NewTemplate()

		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		editTemplate.Render(w, "editmealplan", editData)
	}
}
