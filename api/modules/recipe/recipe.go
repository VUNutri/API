package recipe

import (
	"encoding/json"
	"log"
	"net/http"

	"../db"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Product struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Value    int    `json:"value"`
	Calories int    `json:"calories"`
	Carbs    int    `json:"carbs"`
	Proteins int    `json:"proteins"`
}

type Recipe struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Category     int    `json:"category"`
	Time         int    `json:"time"`
	Image        string `json:"image"`
	Instructions string `json:"instructions"`
	Calories     int
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
	router.Get("/getAll", getAllRecipes)
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
	defer db.Close()

	render.JSON(w, r, "Recipe was created")
}

func getAllRecipes(w http.ResponseWriter, r *http.Request) {
	db := db.InitDB()
	result, err := db.Query("SELECT * FROM recipes")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	recipes := []Recipe{}
	for result.Next() {
		var recipe Recipe
		err := result.Scan(&recipe.ID, &recipe.Title, &recipe.Category, &recipe.Time, &recipe.Image, &recipe.Instructions)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		recipes = append(recipes, recipe)
	}

	for idx, recipe := range recipes {
		result, er := db.Query("SELECT products.title, ingredients.value, products.calories, products.proteins, products.carbs FROM ingredients LEFT JOIN products ON ingredients.productId = products.id WHERE ingredients.recipeId = ?", recipe.ID)
		if er != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		for result.Next() {
			var product Product
			err := result.Scan(&product.Title, &product.Value, &product.Calories, &product.Proteins, &product.Carbs)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			recipes[idx].Products = append(recipes[idx].Products, product)
			log.Print(recipes[idx].Products)
		}
	}
	render.JSON(w, r, recipes)
}
