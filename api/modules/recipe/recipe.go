package recipe

import (
	"encoding/json"
	"net/http"

	"../db"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Recipe struct {
	ID       int
	Title    string `json:"title"`
	Category string `json:"category"`
	Time     int    `json:"time"`
	Image    string `json:"image"`
}

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/create", createRecipe)
}

func createRecipe(w http.ResponseWriter, r *http.Request) {
	var recipe Recipe
	json.NewDecoder(r.Body).Decode(&recipe)
	db := db.InitDB()

	query, err := db.Prepare("INSERT INTO recipes(title, category, time, image) VALUES(?,?,?,?)")
	if err != nil {
		http.Error(w, "Can not create recipe", 400)
		return
	}

	_, er := query.Exec(recipe.Title, recipe.Category, recipe.Time, recipe.Image)
	defer db.Close()
	if er != nil {
		http.Error(w, "Can not create recipe", 400)
		return
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
