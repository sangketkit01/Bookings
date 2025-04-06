package models

import "github.com/sangketkit01/bookings/internal/forms"

// TemplateData holds data to send to the front-end
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
	Form      *forms.Form
	IsAuthenticated int
}
