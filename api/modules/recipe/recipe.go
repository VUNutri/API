package recipe

import (
	"encoding/json"
	"net/http"

	"../db"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Product struct {
	ID    int `json:"id"`
	Value int `json:"value"`
}

type Recipe struct {
	ID           int
	Title        string    `json:"title"`
	Category     int       `json:"category"`
	Time         int       `json:"time"`
	Image        string    `json:"image"`
	Instructions string    `json:"instructions"`
	Products     []Product `json:"products"`
}

type Ingredients struct {
	ID        int
	RecipeID  int `json:"recipeId"`
	ProductID int `json:"productId"`
	Value     int `json:"value"`
}

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/create", createRecipe)
	return router
}

func createRecipe(w http.ResponseWriter, r *http.Request) {
	var recipe Recipe
	json.NewDecoder(r.Body).Decode(&recipe)
	db := db.InitDB()

	query, err := db.Prepare("INSERT INTO recipes(title, category, time, image, instructions) VALUES(?,?,?,?,?)")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	res, er := query.Exec(recipe.Title, recipe.Category, recipe.Time, recipe.Image, recipe.Instructions)
	if er != nil {
		http.Error(w, "Can not create recipe", 400)
		return
	}

	recipeID, err := res.LastInsertId()
	if err != nil {
		http.Error(w, "Can not create recipe", 400)
		return
	}

	for _, product := range recipe.Products {
		query, err = db.Prepare("INSERT INTO ingredients(recipeid, productid, value) VALUES(?,?,?)")
		_, er = query.Exec(recipeID, product.ID, product.Value)
		if er != nil {
			http.Error(w, "Can not create recipe", 400)
			return
		}
	}

	render.JSON(w, r, "Recipe was created")
}

func getAllRecipes(w http.ResponseWriter, r *http.Request) {
	db := db.InitDB()
	result, err := db.Query("SELECT * FROM recipes")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	defer db.Close()

	recipes := []Recipe{}
	for result.Next() {
		var recipe Recipe
		err := result.Scan(&recipe.ID, &recipe.Title, &recipe.Category, &recipe.Time, &recipe.Image)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		recipes = append(recipes, recipe)
	}
	render.JSON(w, r, recipes)
}
