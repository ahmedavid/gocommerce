package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/ahmedavid/gocommerce/pkg/models"
)

type ShoppingCartData struct {
	Products []models.Product
	Total    int
}

type templateData struct {
	Categories   []models.Category
	Products     []models.Product
	Product      *models.Product
	FormData     url.Values
	FormErrors   map[string]string
	Flash        string
	ShoppingCart *ShoppingCartData
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.html"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.html"))
		if err != nil {
			return nil, err
		}
		// ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.html"))
		// if err != nil {
		// 	return nil, err
		// }
		cache[name] = ts
	}

	return cache, nil
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}
	err := ts.Execute(w, td)
	if err != nil {
		app.serverError(w, err)
	}
}
