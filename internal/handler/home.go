package handler

import (
	"html/template"
	"net/http"
	"time"
)

var templates *template.Template

// InitTemplates parses all templates from the templates directory.
func InitTemplates() error {
	var err error
	templates, err = template.ParseGlob("templates/*.html")
	if err != nil {
		return err
	}
	// Parse partials
	templates, err = templates.ParseGlob("templates/partials/*.html")
	return err
}

// Home redirects to today's day view.
func Home(w http.ResponseWriter, r *http.Request) {
	today := time.Now().Format("2006-01-02")
	http.Redirect(w, r, "/day/"+today, http.StatusFound)
}
