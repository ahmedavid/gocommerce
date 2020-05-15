package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/ahmedavid/gocommerce/pkg/models"
)

func (app *application) emptyCart(w http.ResponseWriter, r *http.Request) {
	app.session.Destroy(r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) addToCart(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
		return
	}

	ID := r.PostForm.Get("id")
	IDNum, err := strconv.Atoi(ID)
	CatID := r.PostForm.Get("cat_id")
	CatIDNum, err := strconv.Atoi(CatID)
	Name := r.PostForm.Get("name")
	Price := r.PostForm.Get("price")
	PriceNum, err := strconv.ParseFloat(Price, 64)
	Stock := r.PostForm.Get("stock")
	StockNum, err := strconv.Atoi(Stock)
	ImgURL := r.PostForm.Get("img_url")
	Description := r.PostForm.Get("description")

	var productsInCart []models.Product
	cartCookie := app.session.GetBytes(r, "cart")
	if cartCookie != nil {
		err = json.Unmarshal(cartCookie, &productsInCart)
		if err != nil {
			app.serverError(w, err)
			return
		}
	}
	productsInCart = append(productsInCart, models.Product{
		ID:          IDNum,
		CatID:       CatIDNum,
		Name:        Name,
		Price:       PriceNum,
		Stock:       StockNum,
		ImgURL:      ImgURL,
		Description: Description,
	})
	cookieBytes, err := json.Marshal(productsInCart)
	if err != nil {
		app.serverError(w, err)
	}
	app.session.Put(r, "cart", cookieBytes)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	catID, err := strconv.Atoi(r.URL.Query().Get("cat_id"))
	categories, err := app.categories.GetAll()
	if err != nil {
		fmt.Println("CATEGORY DATABASE ERROR")
	}

	products, err := app.products.GetByCategory(catID)
	if err != nil {
		fmt.Println("PRODUCT DATABASE ERROR", err)
	}
	data, err := app.getShoppingCart(r)

	app.render(w, r, "home.page.html", &templateData{
		Categories:   categories,
		Products:     products,
		ShoppingCart: data,
	})
}

func (app *application) showBuy(w http.ResponseWriter, r *http.Request) {
	data, err := app.getShoppingCart(r)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, r, "buy.page.html", &templateData{
		ShoppingCart: data,
	})
}

func (app *application) buy(w http.ResponseWriter, r *http.Request) {
	// Process Transaction
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) design(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/design.page.html",
		"./ui/html/base.layout.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/about.page.html",
		"./ui/html/base.layout.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}
