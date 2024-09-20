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
    <tr hx-delete="/mealplans/{{ .ID }}" hx-target="#mealplans" hx-confirm="Are you sure you want to delete this meal plan?">
      <td>{{ .RelativeTime }} ({{ .Date.Format "02.Jan.2006" }})</td>
      <td>{{ .Name }}</td>
      <td>{{ .Meal }}</td>
    </tr>
    {{ end }}
  </tbody>
</table>
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
