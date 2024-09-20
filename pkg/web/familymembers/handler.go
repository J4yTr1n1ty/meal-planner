package familymembers

import (
	"encoding/json"
	"net/http"

	"github.com/J4yTr1n1ty/meal-planner/pkg/boot"
	"github.com/J4yTr1n1ty/meal-planner/pkg/models"
)

func NewHandler() *Handler {
	return &Handler{}
}

type Handler struct{}

func (h *Handler) GetFamilyMembers(w http.ResponseWriter, r *http.Request) {
	familyMembers := []models.FamilyMember{}
	boot.DB.Find(&familyMembers)
	names := []string{}
	for _, familyMember := range familyMembers {
		names = append(names, familyMember.Name)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(names)
}
