package calendar

import (
	"fmt"
	"net/http"
	"time"

	"github.com/J4yTr1n1ty/meal-planner/pkg/boot"
	"github.com/J4yTr1n1ty/meal-planner/pkg/models"
	"github.com/J4yTr1n1ty/meal-planner/pkg/web/htmx"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) CalendarPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		http.ServeFile(w, r, "static/calendar.html")
	}
}

func (h *Handler) GetCalendarData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		currentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

		monthParam := r.URL.Query().Get("month")
		var firstOfMonth time.Time
		if monthParam != "" {
			parsed, err := time.Parse("2006-01", monthParam)
			if err != nil {
				htmx.RenderError(w, http.StatusBadRequest, "Invalid month format, expected YYYY-MM")
				return
			}
			firstOfMonth = time.Date(parsed.Year(), parsed.Month(), 1, 0, 0, 0, 0, now.Location())
		} else {
			firstOfMonth = currentMonth
		}

		firstOfNext := firstOfMonth.AddDate(0, 1, 0)

		var mealPlans []models.MealPlan
		boot.DB.Joins("FamilyMember").Joins("Meal").
			Where("date >= ? AND date < ?", firstOfMonth, firstOfNext).
			Order("date").
			Find(&mealPlans)

		// Group meals by day-of-month
		mealsByDay := make(map[int][]htmx.CalendarMeal)
		for _, mp := range mealPlans {
			day := mp.Date.Day()
			mealsByDay[day] = append(mealsByDay[day], htmx.CalendarMeal{
				Meal:         mp.Meal.Name,
				FamilyMember: mp.FamilyMember.Name,
			})
		}

		// Build weeks: Mon=0 ... Sun=6
		// time.Weekday: Sun=0, Mon=1 ... Sat=6
		startWeekday := int(firstOfMonth.Weekday())
		// Convert to Mon-based (Mon=0, ..., Sun=6)
		startOffset := (startWeekday + 6) % 7

		daysInMonth := firstOfNext.AddDate(0, 0, -1).Day()

		var weeks [][]htmx.CalendarDay
		var week []htmx.CalendarDay

		// Padding cells before the first day
		for i := 0; i < startOffset; i++ {
			week = append(week, htmx.CalendarDay{})
		}

		for day := 1; day <= daysInMonth; day++ {
			date := time.Date(firstOfMonth.Year(), firstOfMonth.Month(), day, 0, 0, 0, 0, now.Location())
			cell := htmx.CalendarDay{
				DayNum: day,
				Date:   date,
				Meals:  mealsByDay[day],
			}
			week = append(week, cell)
			if len(week) == 7 {
				weeks = append(weeks, week)
				week = nil
			}
		}

		// Padding cells after the last day
		if len(week) > 0 {
			for len(week) < 7 {
				week = append(week, htmx.CalendarDay{})
			}
			weeks = append(weeks, week)
		}

		prevMonth := firstOfMonth.AddDate(0, -1, 0).Format("2006-01")
		var nextMonth string
		if firstOfMonth.Before(currentMonth) || firstOfMonth.Equal(currentMonth) {
			// Only show next if we're not already at the current month
			if firstOfMonth.Before(currentMonth) {
				nextMonth = firstOfNext.Format("2006-01")
			}
		}

		data := htmx.CalendarData{
			MonthName: fmt.Sprintf("%s %d", firstOfMonth.Month().String(), firstOfMonth.Year()),
			PrevMonth: prevMonth,
			NextMonth: nextMonth,
			Weeks:     weeks,
		}

		if err := htmx.RenderCalendar(w, data); err != nil {
			htmx.RenderError(w, http.StatusInternalServerError, err.Error())
		}
	}
}
