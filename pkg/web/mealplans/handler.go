package mealplans

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"time"

	"github.com/J4yTr1n1ty/meal-planner/pkg/boot"
	"github.com/J4yTr1n1ty/meal-planner/pkg/config"
	"github.com/J4yTr1n1ty/meal-planner/pkg/models"
	"github.com/J4yTr1n1ty/meal-planner/pkg/web/htmx"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func GetMealPlanData() ([]models.MealPlan, []htmx.MealTableData) {
	mealplans := []models.MealPlan{}
	boot.DB.Joins("FamilyMember").Joins("Meal").Order("date").Find(&mealplans)

	filtered_mealplans := []models.MealPlan{}
	for _, mealplan := range mealplans {
		if time.Since(mealplan.Date).Hours() < 24 {
			filtered_mealplans = append(filtered_mealplans, mealplan)
		}
	}

	tableData := []htmx.MealTableData{}

	for _, mealplan := range filtered_mealplans {
		tableData = append(tableData, htmx.MealTableData{
			ID:           mealplan.ID,
			RelativeTime: mealplan.GetDaysFromNow(),
			Date:         mealplan.Date,
			Name:         mealplan.FamilyMember.Name,
			Meal:         mealplan.Meal.Name,
		})
	}

	return filtered_mealplans, tableData
}

func (h *Handler) GetMealPlans() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mealplans, tableData := GetMealPlanData()

		accept_header := r.Header["Accept"]

		if slices.Contains(accept_header, "application/json") {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mealplans)
		} else {
			w.Header().Set("Content-Type", "text/html")

			if err := htmx.RenderMealTable(w, tableData); err != nil {
				htmx.RenderError(w, http.StatusInternalServerError, err.Error())
			}
		}
	}
}

func (h *Handler) UpdateMealPlan() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			htmx.RenderError(w, http.StatusBadRequest, "Missing id")
			return
		}

		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		cooking := r.Form.Get("cooking")
		date := r.Form.Get("date")
		meal := r.Form.Get("meal")

		if cooking == "" || date == "" || meal == "" {
			htmx.RenderError(w, http.StatusBadRequest, "Missing data")
			return
		}

		family_member := models.FamilyMember{}
		log.Println("Searching for Person: ", cooking)
		boot.DB.Where("lower(name) = lower(?)", cooking).First(&family_member)
		if family_member.ID == 0 {
			log.Println("Creating new person: ", cooking)
			family_member = models.FamilyMember{Name: cooking}
			boot.DB.Create(&family_member)
		}

		mealFromDB := models.Meal{}
		log.Println("Searching for Meal: ", meal)
		boot.DB.Where("lower(name) = lower(?)", meal).First(&mealFromDB)
		if mealFromDB.ID == 0 {
			log.Println("Creating new meal: ", meal)
			mealFromDB = models.Meal{Name: meal}
			boot.DB.Create(&mealFromDB)
		}

		dateTime, err := time.Parse("2006-01-02", date)
		if err != nil {
			htmx.RenderError(w, http.StatusBadRequest, "Invalid date")
			return
		}

		mealPlan := models.MealPlan{
			FamilyMember: family_member,
			Meal:         mealFromDB,
			Date:         dateTime,
		}

		result := boot.DB.Where("id = ?", id).Updates(&mealPlan)
		if config.IsDebug() {
			log.Println("Updated rows: ", result.RowsAffected)
		}
		if result.Error != nil {
			htmx.RenderError(w, http.StatusInternalServerError, "Unable to update mealplan in Database")
		}

		htmx.Redirect(w, r, "/")
	}
}

func (h *Handler) CreateMealPlan() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		cooking := r.Form.Get("cooking")
		date := r.Form.Get("date")
		meal := r.Form.Get("meal")

		family_member := models.FamilyMember{}
		log.Println("Searching for Person: ", cooking)
		boot.DB.Where("lower(name) = lower(?)", cooking).First(&family_member)
		if family_member.ID == 0 {
			log.Println("Creating new person: ", cooking)
			family_member = models.FamilyMember{Name: cooking}
			boot.DB.Create(&family_member)
		}

		mealFromDB := models.Meal{}
		log.Println("Searching for Meal: ", meal)
		boot.DB.Where("lower(name) = lower(?)", meal).First(&mealFromDB)
		if mealFromDB.ID == 0 {
			log.Println("Creating new meal: ", meal)
			mealFromDB = models.Meal{Name: meal}
			boot.DB.Create(&mealFromDB)
		}

		dateTime, err := time.Parse("2006-01-02", date)
		if err != nil {
			htmx.RenderError(w, http.StatusBadRequest, "Invalid date")
			return
		}

		plan := models.MealPlan{
			Date:         dateTime,
			FamilyMember: family_member,
			Meal:         mealFromDB,
		}

		result := boot.DB.Create(&plan)

		if result.Error != nil {
			htmx.RenderError(w, http.StatusInternalServerError, result.Error.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		htmx.RenderSuccess(w, "Meal plan created successfully")
	}
}

func (h *Handler) DeleteMealPlan() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		if id == "" {
			htmx.RenderError(w, http.StatusBadRequest, "Invalid ID")
			return
		}

		tx := boot.DB.Delete(&models.MealPlan{}, id)

		if tx.Error != nil {
			htmx.RenderError(w, http.StatusInternalServerError, tx.Error.Error())
			return
		}

		mealplans, tableData := GetMealPlanData()

		accept_header := r.Header["Accept"]

		if slices.Contains(accept_header, "application/json") {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mealplans)
			return
		} else {
			w.Header().Set("Content-Type", "text/html")
			htmx.RenderMealTable(w, tableData)
			return
		}
	}
}
