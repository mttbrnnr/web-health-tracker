package handler

import (
	"html/template"
	"io/fs"
	"net/http"
	"time"
)

var templates *template.Template

// InitTemplates parses all templates from the embedded filesystem.
func InitTemplates(fsys fs.FS) error {
	var err error
	templates, err = template.ParseFS(fsys, "templates/*.html", "templates/partials/*.html")
	return err
}

// Home redirects to today's day view.
func Home(w http.ResponseWriter, r *http.Request) {
	today := time.Now().Format("2006-01-02")
	http.Redirect(w, r, "/day/"+today, http.StatusFound)
}
