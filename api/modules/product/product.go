package product

import (
	"app/modules/db"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Product struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Size     string `json:"size"`
	Calories int    `json:"calories"`
	Carbs    int    `json:"carbs"`
	Proteins int    `json:"proteins"`
}

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/create", createProduct)
	router.Get("/getAll", getAllProducts)
	router.Get("/getById/{productId}", getProductByID)
	return router
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	//session, _ := auth.Store.Get(r, "cookie")
	//if auth, ok := session.Values["auth"].(bool); !ok || !auth {
	//http.Error(w, "Forbidden", http.StatusForbidden)
	//return
	//}

	var product Product
	json.NewDecoder(r.Body).Decode(&product)

	if !checkIfValid(product) {
		http.Error(w, "Bad request", 400)
		return
	}

	db := db.InitDB()

	query, err := db.Prepare("INSERT INTO products(title, calories, carbs, proteins, size) VALUES(?,?,?,?,?)")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	_, er := query.Exec(product.Title, product.Calories, product.Carbs, product.Proteins, product.Size)
	defer db.Close()
	if er != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	render.JSON(w, r, "Product was created")
}

func getAllProducts(w http.ResponseWriter, r *http.Request) {
	db := db.InitDB()
	result, err := db.Query("SELECT * FROM products")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	defer db.Close()

	products := []Product{}
	for result.Next() {
		var product Product
		err := result.Scan(&product.ID, &product.Title, &product.Calories, &product.Carbs, &product.Proteins, &product.Size)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		products = append(products, product)
	}
	render.JSON(w, r, products)
}

func getProductByID(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "productId")
	db := db.InitDB()
	result, err := db.Query("SELECT * FROM products WHERE id = ?", productID)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	defer db.Close()

	product := Product{}
	for result.Next() {
		err := result.Scan(&product.ID, &product.Title, &product.Calories, &product.Carbs, &product.Proteins, &product.Size)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
	render.JSON(w, r, product)
}

func checkIfValid(p Product) bool {
	if len(p.Title) < 4 {
		return false
	}
	if len(p.Size) < 1 {
		return false
	}
	if p.Calories == 0 {
		return false
	}
	if p.Carbs == 0 {
		return false
	}
	if p.Proteins == 0 {
		return false
	}
	return true
}
