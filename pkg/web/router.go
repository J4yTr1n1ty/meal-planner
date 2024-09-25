package web

import (
	"net/http"

	"github.com/J4yTr1n1ty/meal-planner/pkg/web/authentication"
	"github.com/J4yTr1n1ty/meal-planner/pkg/web/familymembers"
	"github.com/J4yTr1n1ty/meal-planner/pkg/web/homepage"
	"github.com/J4yTr1n1ty/meal-planner/pkg/web/mealplans"
	"github.com/J4yTr1n1ty/meal-planner/pkg/web/meals"
	"github.com/J4yTr1n1ty/meal-planner/pkg/web/middleware"
)

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("GET /static/", middleware.LoginRequired(http.StripPrefix("/static/", fs)))

	homepageHandler := homepage.NewHandler()
	mux.Handle("GET /favicon.ico", homepageHandler.Favicon())
	mux.Handle("GET /", middleware.LoginRequired(homepageHandler.Homepage()))
	mux.Handle("GET /addmeal", middleware.LoginRequired(homepageHandler.AddMealPage()))
	mux.Handle("GET /editmeal/{id}", middleware.LoginRequired(homepageHandler.EditMealPage()))

	authenticationHandler := authentication.NewHandler()
	mux.Handle("GET /login", authenticationHandler.Login())
	mux.Handle("POST /login", authenticationHandler.LoginPost())

	mealPlanHandler := mealplans.NewHandler()
	mux.Handle("GET /mealplans", middleware.LoginRequired(mealPlanHandler.GetMealPlans()))
	mux.Handle("POST /mealplans", middleware.LoginRequired(mealPlanHandler.CreateMealPlan()))
	mux.Handle("PUT /mealplans/{id}", middleware.LoginRequired(mealPlanHandler.UpdateMealPlan()))
	mux.Handle("DELETE /mealplans/{id}", middleware.LoginRequired(mealPlanHandler.DeleteMealPlan()))

	mealHandler := meals.NewHandler()
	mux.Handle("GET /meals", middleware.LoginRequired(mealHandler.GetMeals()))

	familyMemberHandler := familymembers.NewHandler()
	mux.Handle("GET /familymembers", middleware.LoginRequired(familyMemberHandler.GetFamilyMembers()))

	return mux
}
