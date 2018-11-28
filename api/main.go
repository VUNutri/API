package main

import (
	"app/modules/menu"
	"app/modules/product"
	"app/modules/recipe"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.RedirectSlashes,
		middleware.Recoverer,
	)

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/api/products", product.Routes())
		r.Mount("/api/recipes", recipe.Routes())
		r.Mount("/api/menu", menu.Routes())
	})

	return router
}

func main() {
	router := routes()

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Println(err)
	}
}
