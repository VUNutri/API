package product

import (
	"encoding/json"
	"log"
	"net/http"

	"../db"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Product struct {
	ID       int
	Title    string `json:"title"`
	Calories int    `json:"calories"`
	Carbs    int    `json:"carbs"`
	Proteins int    `json:"proteins"`
}

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/create", createProduct)
	router.Get("/getAll", getAllProducts)
	return router
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	json.NewDecoder(r.Body).Decode(&product)
	db := db.InitDB()

	log.Printf(product.Title)

	query, err := db.Prepare("INSERT INTO products(title, calories, carbs, proteins) VALUES(?,?,?,?)")
	if err != nil {
		http.Error(w, "Can not create product", 400)
		return
	}

	_, er := query.Exec(product.Title, product.Calories, product.Carbs, product.Proteins)
	defer db.Close()
	if er != nil {
		http.Error(w, "Can not create product", 400)
		return
	}

	render.JSON(w, r, "Post was created")
}

func getAllProducts(w http.ResponseWriter, r *http.Request) {
	db := db.InitDB()
	result, err := db.Query("SELECT * FROM products")
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	defer db.Close()

	products := []Product{}
	for result.Next() {
		var product Product
		result.Scan(&product.ID, &product.Title, &product.Calories, &product.Carbs, &product.Proteins)
		products = append(products, product)
	}
	render.JSON(w, r, products)
}
