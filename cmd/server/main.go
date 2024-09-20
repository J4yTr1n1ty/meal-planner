package main

import (
	"log"
	"net/http"

	"github.com/J4yTr1n1ty/meal-planner/pkg/boot"
	"github.com/J4yTr1n1ty/meal-planner/pkg/models"
	"github.com/J4yTr1n1ty/meal-planner/pkg/web"
	"github.com/J4yTr1n1ty/meal-planner/pkg/web/middleware"
)

func main() {
	log.Println("Starting Meal Planner...")

	boot.LoadEnv()
	boot.ConnectDB()
	err := boot.DB.AutoMigrate(&models.FamilyMember{}, &models.Meal{}, &models.MealPlan{})
	if err != nil {
		log.Fatal(err)
	}

	router := web.SetupRoutes()

	stack := middleware.CreateStack(middleware.Logging)

	server := http.Server{
		Addr:    ":" + boot.Environment.GetEnv("PORT"),
		Handler: stack(router),
	}

	log.Println("Listening on port :" + boot.Environment.GetEnv("PORT"))
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Error starting server: " + err.Error())
	}
}
