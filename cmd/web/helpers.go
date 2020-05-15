package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime/debug"

	"github.com/ahmedavid/gocommerce/pkg/models"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) badRequest(w http.ResponseWriter) {
	app.clientError(w, http.StatusBadRequest)
}

func (app *application) uploadFile(w http.ResponseWriter, r *http.Request) (fileName string, err error) {

	file, handler, err := r.FormFile("image")
	if err != nil {
		return "", err
	}

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	defer file.Close()

	// u := uuid.NewV4()

	tempFile, err := ioutil.TempFile("ui/static/uploads", "upload-*.png")
	if err != nil {
		return "", err
	}
	defer tempFile.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	tempFile.Write(fileBytes)

	return filepath.Base(tempFile.Name()), nil
}

func (app *application) getShoppingCart(r *http.Request) (*ShoppingCartData, error) {

	var productsInCart []models.Product
	cartCookie := app.session.GetBytes(r, "cart")
	if cartCookie == nil {
		productsInCart = make([]models.Product, 0)
	} else {
		err := json.Unmarshal(cartCookie, &productsInCart)
		if err != nil {
			return nil, err
		}
	}

	// fmt.Println("BYTES: ", cartCookie)
	// fmt.Println("CART: ", productsInCart)

	pMap := reduceProducts(productsInCart)

	total := 0
	for _, p := range pMap {
		total += p.Stock * int(p.Price)
	}

	return &ShoppingCartData{
		Products: pMap,
		Total:    total,
	}, nil
}

func reduceProducts(products []models.Product) []models.Product {
	productMap := make(map[models.Product]int)

	for _, product := range products {
		if productMap[product] == 0 {
			productMap[product] = 1
		} else {
			productMap[product] = productMap[product] + 1
		}
	}

	productsToReturn := make([]models.Product, 0)

	for k, v := range productMap {
		k.Stock = v
		productsToReturn = append(productsToReturn, k)
	}

	return productsToReturn
}
