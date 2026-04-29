package handler

import (
	"context"
	"html/template"
	"net/http"

	"github.com/techies/orders-api/helpers"
)

// ShowLandingPage renders the HTML landing page
func ShowLandingPage(w http.ResponseWriter, r *http.Request) {
	// Fetch books from repository (assume global or injected BookRepo)
	var books []interface{}
	if repo, ok := r.Context().Value("bookRepo").(interface {
		List(context.Context) ([]interface{}, error)
	}); ok {
		var err error
		books, err = repo.List(r.Context())
		if err != nil {
			helpers.WriteJSONError(w, "Failed to fetch books", http.StatusInternalServerError)
			return
		}
	}

	tmpl, err := template.ParseFiles(
		"views/layouts/RootLayout.html",
		"views/pages/index.html",
	)
	if err != nil {
		helpers.WriteJSONError(w, "Failed to load template", http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"Books": books,
	}
	if err := tmpl.ExecuteTemplate(w, "RootLayout", data); err != nil {
		helpers.WriteJSONError(w, "Failed to render template", http.StatusInternalServerError)
	}
}
