package web

import (
	"fmt"
	"html/template"
	"net/http"
)

type Renderer struct {
	pages     map[string]*template.Template
	fragments *template.Template
}

func NewRenderer() (*Renderer, error) {
	shared, err := template.ParseFS(
		assets,
		"templates/layouts/*.html",
		"templates/components/*.html",
		"templates/partials/*.html",
	)

	if err != nil {
		return nil, fmt.Errorf("parse shard templates: %w", err)
	}

	pages := map[string]string{
		"home":     "templates/pages/home.html",
		"play":     "templates/pages/play.html",
		"library":  "templates/pages/library.html",
		"engine":   "templates/pages/engines.html",
		"game":     "templates/pages/game.html",
		"settings": "templates/pages/settings.html",
	}

	parsedPages := make(map[string]*template.Template, len(pages))

	for name, pattern := range pages {
		pageTemplates, err := shared.Clone()
		if err != nil {
			return nil, fmt.Errorf("clone templates for page %q: %w", name, err)
		}

		pageTemplates, err = pageTemplates.ParseFS(assets, pattern)
		if err != nil {
			return nil, fmt.Errorf("parse page template %q: %w", name, err)
		}

		parsedPages[name] = pageTemplates
	}

	return &Renderer{
		pages:     parsedPages,
		fragments: shared,
	}, nil
}

func (r *Renderer) RenderPage(w http.ResponseWriter, page, templateName string, data any) {
	tmpl, ok := r.pages[page]
	if !ok {
		http.Error(w, "page template not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := tmpl.ExecuteTemplate(w, templateName, data); err != nil {
		http.Error(w, "template render error", http.StatusInternalServerError)
	}
}

func (r *Renderer) RenderFragment(w http.ResponseWriter, templateName string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := r.fragments.ExecuteTemplate(w, templateName, data); err != nil {
		http.Error(w, "template render error", http.StatusInternalServerError)
	}
}
