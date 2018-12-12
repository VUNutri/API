package main

import (
	"app/modules/auth"
	"app/modules/menu"
	"app/modules/product"
	"app/modules/recipe"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func routes() *chi.Mux {
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	router := chi.NewRouter()
	router.Use(
		cors.Handler,
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
		r.Mount("/api/auth", auth.Routes())
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
