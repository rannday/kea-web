package handlers

import (
	"embed"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/rannday/kea-web/internal/utils"
)

//go:embed templates/*.html
var templatesFS embed.FS

// PageData holds dynamic content and asset filenames
type PageData struct {
	Title     string
	CSSBundle string
	JSBundle  string
	Data      map[string]interface{}
}

// Cached asset filenames (set once at startup)
var (
	cssBundle string
	jsBundle  string
)

// SetBundledAssets allows `main.go` to pass the filenames once generated
func SetBundledAssets(css, js string) {
	cssBundle = filepath.Base(css)
	jsBundle = filepath.Base(js)
}

// RenderTemplate loads layout.html + specific content template
func RenderTemplate(w http.ResponseWriter, tmpl string, data PageData) {
	layout := "templates/layout.html"
	content := "templates/" + tmpl + ".html"

	utils.Debug("Parsing templates: layout=%s, content=%s", layout, content)

	// Inject hashed bundle filenames
	data.CSSBundle = cssBundle
	data.JSBundle = jsBundle

	t, err := template.ParseFS(templatesFS, layout, content)
	if err != nil {
		utils.Error("Failed to parse templates: %v", err)
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	if err := t.ExecuteTemplate(w, "layout", data); err != nil {
		utils.Error("Failed to execute template: %v", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
