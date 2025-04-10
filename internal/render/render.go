package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/justinas/nosurf"
	"github.com/sangketkit01/bookings/internal/config"
	"github.com/sangketkit01/bookings/internal/models"
)

var app *config.AppConfig
var functions = template.FuncMap{
	"humanDate" : HumanDate,
	"formatDate" : FormatDate,
	"iterate" : Iterate,
	"add" : Add,
}
var pathToTemplates = "../../templates"

// NewRenderer sets the config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

// HumanDate returns time in YYYY-MM-DD format
func HumanDate(t time.Time) string{
	return t.Format("2006-01-02")
}

// FormatDate nothing
func FormatDate(t time.Time, f string) string{
	return t.Format(f)
}


// Iterate returns a slice of int starting at 1, goint to count
func Iterate(count int) []int{
	var i int 
	var items []int
	for i = 0 ; i < count ; i++ {
		items = append(items, i)
	}

	return items
}

// Add returns
func Add(a,b int) int {
	return a + b
}

// AddDefaultData adds default data to our app
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(),"flash")
	td.Warning = app.Session.PopString(r.Context(),"warning")
	td.Error = app.Session.PopString(r.Context(),"error")
	td.CSRFToken = nosurf.Token(r)
	
	if app.Session.Exists(r.Context(),"user_id"){
		td.IsAuthenticated = 1
	}

	return td
}

// Template renders the template file to show in the browser
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {

	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// get request template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Println("template not found")
		return errors.New("Can't get template from cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	if err := t.Execute(buf, td); err != nil {
		log.Fatalln(err)
	}

	// render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}

	return nil

}

// CreateTemplateCache nothing
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl",pathToTemplates))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl",pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl",pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
