package web

import (
	"net/http"

	"github.com/J4yTr1n1ty/meal-planner/pkg/web/familymembers"
	"github.com/J4yTr1n1ty/meal-planner/pkg/web/homepage"
	"github.com/J4yTr1n1ty/meal-planner/pkg/web/mealplans"
	"github.com/J4yTr1n1ty/meal-planner/pkg/web/meals"
)

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	homepageHandler := homepage.NewHandler()
	mux.HandleFunc("GET /", homepageHandler.Homepage)
	mux.HandleFunc("GET /addmeal", homepageHandler.AddMealPage)

	mealPlanHandler := mealplans.NewHandler()
	mux.HandleFunc("GET /mealplans", mealPlanHandler.GetMealPlans)
	mux.HandleFunc("POST /mealplans", mealPlanHandler.CreateMealPlan)
	mux.HandleFunc("DELETE /mealplans/{id}", mealPlanHandler.DeleteMealPlan)

	mealHandler := meals.NewHandler()
	mux.HandleFunc("GET /meals", mealHandler.GetMeals)

	familyMemberHandler := familymembers.NewHandler()
	mux.HandleFunc("GET /familymembers", familyMemberHandler.GetFamilyMembers)

	return mux
}
