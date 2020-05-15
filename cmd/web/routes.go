package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable)
	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Post("/addToCart", dynamicMiddleware.ThenFunc(app.addToCart))
	mux.Get("/emptyCart", dynamicMiddleware.ThenFunc(app.emptyCart))
	mux.Get("/admin", dynamicMiddleware.ThenFunc(app.dashboard))
	mux.Get("/admin/createProduct", dynamicMiddleware.ThenFunc(app.showCreateProduct))
	mux.Post("/admin/createProduct", dynamicMiddleware.ThenFunc(app.createProduct))
	mux.Get("/admin/deleteProduct", dynamicMiddleware.ThenFunc(app.deleteProduct))
	mux.Get("/buy", dynamicMiddleware.ThenFunc(app.showBuy))
	mux.Post("/buy", dynamicMiddleware.ThenFunc(app.buy))
	mux.Get("/about", dynamicMiddleware.ThenFunc(app.about))
	mux.Get("/design", dynamicMiddleware.ThenFunc(app.design))

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	return standardMiddleware.Then(mux)
}
