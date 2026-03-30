package htmx

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

var (
	ErrParsingTemplate   = fmt.Errorf("error parsing template")
	ErrExecutingTemplate = fmt.Errorf("error executing template")
)

var ErrorTemplate = `
<div class="alert alert-danger">
  <p>{{ .Message }}</p>
</div>
`

var SuccessTemplate = `
<div class="alert alert-success">
  <p>{{ .Message }}</p>
</div>
`

var InfoTemplate = `
<div class="alert alert-info">
  <p>{{ .Message }}</p>
</div>
`

var EmptyTableTemplate = `
<table class="table table-hover">
  <thead>
    <tr>
      <th scope="col">Relative Time</th>
      <th scope="col">Date</th>
      <th scope="col">Name</th>
    </tr>
  </thead>
  <tbody>
  </tbody>
</table>
`

type MealTableData struct {
	ID           uint
	RelativeTime string
	Date         time.Time
	Name         string
	Meal         string
}

var MealTableTemplate = `
<table class="table table-hover">
  <thead>
    <tr>
      <th scope="col">Date</th>
      <th scope="col">Name</th>
      <th scope="col">Meal</th>
    </tr>
  </thead>
  <tbody>
    {{ range . }}
    <tr onclick="window.location='/editmeal/{{ .ID }}'" class="mealplanentry">
      <td>{{ .RelativeTime }} ({{ .Date.Format "02.Jan.2006" }})</td>
      <td>{{ .Name }}</td>
      <td>{{ .Meal }}</td>
    </tr>
    {{ end }}
  </tbody>
</table>
`

type CalendarMeal struct {
	Meal         string
	FamilyMember string
}

type CalendarDay struct {
	DayNum int // 0 = padding cell
	Date   time.Time
	Meals  []CalendarMeal
}

type CalendarData struct {
	MonthName string        // e.g. "March 2026"
	PrevMonth string        // e.g. "2026-02"
	NextMonth string        // e.g. "2026-04", empty string if current month
	Weeks     [][]CalendarDay // rows of 7 days, Mon–Sun
}

var CalendarTemplate = `
<div id="calendar-grid" class="mt-3">
  <div class="d-flex justify-content-between align-items-center mb-2">
    <button class="btn btn-outline-secondary btn-sm"
            hx-get="/calendardata?month={{ .PrevMonth }}"
            hx-target="#calendar-grid" hx-swap="outerHTML">← Prev</button>
    <strong>{{ .MonthName }}</strong>
    {{ if .NextMonth }}
    <button class="btn btn-outline-secondary btn-sm"
            hx-get="/calendardata?month={{ .NextMonth }}"
            hx-target="#calendar-grid" hx-swap="outerHTML">Next →</button>
    {{ else }}
    <span class="btn btn-outline-secondary btn-sm disabled">Next →</span>
    {{ end }}
  </div>
  <div class="calendar-month">
    <div class="calendar-header d-flex text-center fw-bold mb-1">
      <div style="width:calc(100%/7);flex:0 0 calc(100%/7)">Mon</div>
      <div style="width:calc(100%/7);flex:0 0 calc(100%/7)">Tue</div>
      <div style="width:calc(100%/7);flex:0 0 calc(100%/7)">Wed</div>
      <div style="width:calc(100%/7);flex:0 0 calc(100%/7)">Thu</div>
      <div style="width:calc(100%/7);flex:0 0 calc(100%/7)">Fri</div>
      <div style="width:calc(100%/7);flex:0 0 calc(100%/7)">Sat</div>
      <div style="width:calc(100%/7);flex:0 0 calc(100%/7)">Sun</div>
    </div>
    {{ range .Weeks }}
    <div class="calendar-week d-flex">
      {{ range . }}
      <div class="calendar-cell border p-1 overflow-hidden" style="width:calc(100%/7);flex:0 0 calc(100%/7);min-height:80px">
        {{ if .DayNum }}<small class="text-muted">{{ .DayNum }}</small>{{ end }}
        {{ range .Meals }}
        <div class="badge bg-success text-wrap w-100 mb-1" style="font-size:0.7rem">
          {{ .Meal }}<br><small>{{ .FamilyMember }}</small>
        </div>
        {{ end }}
      </div>
      {{ end }}
    </div>
    {{ end }}
  </div>
</div>
`

func RenderError(w http.ResponseWriter, httpStatus int, message string) error {
	w.WriteHeader(httpStatus)

	data := struct {
		Message string
	}{
		Message: message,
	}

	value, err := renderTemplate(ErrorTemplate, data)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrExecutingTemplate, err)
	}

	w.Write([]byte(value))

	return nil
}

func RenderSuccess(w http.ResponseWriter, message string) error {
	w.WriteHeader(http.StatusOK)

	data := struct {
		Message string
	}{
		Message: message,
	}

	value, err := renderTemplate(SuccessTemplate, data)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrExecutingTemplate, err)
	}

	w.Write([]byte(value))

	return nil
}

func RenderInfo(w http.ResponseWriter, message string) error {
	w.WriteHeader(http.StatusOK)

	data := struct {
		Message string
	}{
		Message: message,
	}

	value, err := renderTemplate(InfoTemplate, data)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrExecutingTemplate, err)
	}

	w.Write([]byte(value))

	return nil
}

func RenderMealTable(w http.ResponseWriter, data []MealTableData) error {
	w.WriteHeader(http.StatusOK)

	if len(data) == 0 {
		value, err := renderTemplate(EmptyTableTemplate, []string{})
		if err != nil {
			return fmt.Errorf("%w: %v", ErrExecutingTemplate, err)
		}

		w.Write([]byte(value))
		return nil
	}

	value, err := renderTemplate(MealTableTemplate, data)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrExecutingTemplate, err)
	}

	w.Write([]byte(value))

	return nil
}

func RenderCalendar(w http.ResponseWriter, data CalendarData) error {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	value, err := renderTemplate(CalendarTemplate, data)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrExecutingTemplate, err)
	}

	w.Write([]byte(value))
	return nil
}

func Redirect(w http.ResponseWriter, r *http.Request, path string) {
	w.Header().Set("HX-Redirect", path)
	w.WriteHeader(http.StatusFound)
}

func renderTemplate(templateStr string, data interface{}) (string, error) {
	tmpl, err := template.New("template").Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrParsingTemplate, err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrExecutingTemplate, err)
	}

	return buf.String(), nil
}
