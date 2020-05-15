package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

func (app *application) dashboard(w http.ResponseWriter, r *http.Request) {
	catID, err := strconv.Atoi(r.URL.Query().Get("cat_id"))
	categories, err := app.categories.GetAll()
	if err != nil {
		fmt.Println("CATEGORY DATABASE ERROR")
	}

	products, err := app.products.GetByCategory(catID)
	if err != nil {
		fmt.Println("PRODUCT DATABASE ERROR", err)
	}

	flash := app.session.PopString(r, "flash")

	app.render(w, r, "dashboard.page.html", &templateData{
		Categories: categories,
		Products:   products,
		Flash:      flash,
	})
}

func (app *application) showCreateProduct(w http.ResponseWriter, r *http.Request) {
	categories, err := app.categories.GetAll()
	if err != nil {
		app.serverError(w, err)
	}
	app.render(w, r, "createProduct.page.html", &templateData{
		Categories: categories,
	})
}

func (app *application) deleteProduct(w http.ResponseWriter, r *http.Request) {
	productID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		app.badRequest(w)
		return
	}
	rowsAffected, err := app.products.DeleteProduct(productID)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if rowsAffected == 1 {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	} else {
		app.serverError(w, fmt.Errorf("SQL ERROR"))
		return
	}
}

func (app *application) createProduct(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	errors := make(map[string]string)
	cat_id := r.PostForm.Get("cat_id")
	name := r.PostForm.Get("name")
	description := r.PostForm.Get("description")
	price := r.PostForm.Get("price")
	stock := r.PostForm.Get("stock")

	fmt.Println(name, price, stock)

	// Validate cat_id
	catIDNum, err := strconv.Atoi(cat_id)
	if err != nil {
		app.serverError(w, err)
	}
	// Validate Price
	if strings.TrimSpace(price) == "" {
		errors["price"] = "This field cannot be blank"
	}
	priceNum, err := strconv.ParseFloat(price, 64)
	if err != nil {
		errors["price"] = "This field should be numeric"
	}

	// Validate stock

	if strings.TrimSpace(stock) == "" {
		errors["stock"] = "This field cannot be blank"
	}
	stockNum, err := strconv.Atoi(stock)
	if err != nil {
		errors["stock"] = "This field should be numeric"
	}

	// Validate Name
	if strings.TrimSpace(name) == "" {
		errors["name"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(name) > 100 {
		errors["name"] = "This field is too long (maximum is 100 characters)"
	}
	// Validate Description
	if strings.TrimSpace(description) == "" {
		errors["description"] = "This field cannot be blank"
	}

	imgUrl, err := app.uploadFile(w, r)
	if err != nil {
		if err == http.ErrMissingFile {
			errors["image"] = "File is empty"
		}
	}

	if len(errors) > 0 {
		fmt.Println(errors)
		fmt.Println(r.PostForm)
		app.render(w, r, "createProduct.page.html", &templateData{
			FormData:   r.PostForm,
			FormErrors: errors,
		})
		return
	}

	_, err = app.products.CreateProduct(catIDNum, name, description, imgUrl, priceNum, stockNum)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.session.Put(r, "flash", "Product created successfully!")
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
