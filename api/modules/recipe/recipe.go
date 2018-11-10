package recipe

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"../db"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Day struct {
	Count   int      `json:"dayCount"`
	Recipes []Recipe `json:"meals"`
}

type Product struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Value    int    `json:"value"`
	Calories int    `json:"calories"`
	Carbs    int    `json:"carbs"`
	Proteins int    `json:"proteins"`
}

type Recipe struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Category     int       `json:"category"`
	Time         int       `json:"time"`
	Image        string    `json:"image"`
	Instructions string    `json:"instructions"`
	Calories     int       `json:"calories"`
	Carbs        int       `json:"carbs"`
	Proteins     int       `json:"proteins"`
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
	router.Get("/getById/{recipeId}", getRecipeById)
	router.Get("/getMenu/{daysCount}/{mealsCount}/{caloriesCount}", getMenu)
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

	res, err := query.Exec(recipe.Title, recipe.Category, recipe.Time, recipe.Image, recipe.Instructions)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	recipeID, err := res.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	recipe.ID = int(recipeID)

	err = createRecipeIngredients(&recipe)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err = sumRecipeNutrition(&recipe)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err = updateRecipeNutrition(&recipe)
	if err != nil {
		http.Error(w, err.Error(), 400)
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

	recipes := []Recipe{}
	for result.Next() {
		var recipe Recipe
		err := result.Scan(&recipe.ID, &recipe.Title, &recipe.Category, &recipe.Time, &recipe.Image, &recipe.Instructions, &recipe.Calories, &recipe.Carbs, &recipe.Proteins)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		recipes = append(recipes, recipe)
	}

	for idx, recipe := range recipes {
		result, er := db.Query("SELECT products.id, products.title, ingredients.value, products.calories, products.proteins, products.carbs FROM ingredients LEFT JOIN products ON ingredients.productId = products.id WHERE ingredients.recipeId = ?", recipe.ID)
		if er != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		for result.Next() {
			var product Product
			err := result.Scan(&product.ID, &product.Title, &product.Value, &product.Calories, &product.Proteins, &product.Carbs)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			recipes[idx].Products = append(recipes[idx].Products, product)
		}
	}
	render.JSON(w, r, recipes)
}

func getMenu(w http.ResponseWriter, r *http.Request) {
	daysCount, err := strconv.Atoi(chi.URLParam(r, "daysCount"))
	mealsCount, err := strconv.Atoi(chi.URLParam(r, "mealsCount"))
	caloriesCount, err := strconv.Atoi(chi.URLParam(r, "caloriesCount"))
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	db := db.InitDB()
	result, err := db.Query("SELECT * FROM recipes")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	recipes := []Recipe{}
	for result.Next() {
		var recipe Recipe
		err := result.Scan(&recipe.ID, &recipe.Title, &recipe.Category, &recipe.Time, &recipe.Image, &recipe.Instructions, &recipe.Calories, &recipe.Carbs, &recipe.Proteins)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		recipes = append(recipes, recipe)
	}

	for idx, recipe := range recipes {
		result, er := db.Query("SELECT products.id, products.title, ingredients.value, products.calories, products.proteins, products.carbs FROM ingredients LEFT JOIN products ON ingredients.productId = products.id WHERE ingredients.recipeId = ?", recipe.ID)
		if er != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		for result.Next() {
			var product Product
			err := result.Scan(&product.ID, &product.Title, &product.Value, &product.Calories, &product.Proteins, &product.Carbs)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			recipes[idx].Products = append(recipes[idx].Products, product)
		}
	}
	log.Println(daysCount, caloriesCount, mealsCount)

	days := []Day{}

	for i := 0; i < daysCount; i++ {
		day := Day{}
		calories := 0
		for meal := 0; meal < mealsCount; meal++ {
			for _, recipe := range recipes {
				if recipe.Calories+calories <= caloriesCount {
					calories += recipe.Calories
					day.Count = i + 1
					day.Recipes = append(day.Recipes, recipe)
					break
				}
			}
		}
		days = append(days, day)
	}

	render.JSON(w, r, days)
}

func getRecipeById(w http.ResponseWriter, r *http.Request) {
	recipeID := chi.URLParam(r, "recipeId")

	db := db.InitDB()
	result, err := db.Query("SELECT * FROM recipes WHERE id = ?", recipeID)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	recipe := Recipe{}
	for result.Next() {
		err = result.Scan(&recipe.ID, &recipe.Title, &recipe.Category, &recipe.Time, &recipe.Image, &recipe.Instructions, &recipe.Calories, &recipe.Carbs, &recipe.Proteins)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}

	result, err = db.Query("SELECT products.id, products.title, ingredients.value, products.calories, products.proteins, products.carbs FROM ingredients LEFT JOIN products ON ingredients.productId = products.id WHERE ingredients.recipeId = ?", recipe.ID)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	for result.Next() {
		var product Product
		err := result.Scan(&product.ID, &product.Title, &product.Value, &product.Calories, &product.Proteins, &product.Carbs)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		recipe.Products = append(recipe.Products, product)
	}
	render.JSON(w, r, recipe)
}

func updateRecipeNutrition(recipe *Recipe) (err error) {
	db := db.InitDB()
	query, err := db.Prepare("UPDATE recipes SET calories = ?, carbs = ?, proteins = ? WHERE id = ?")
	if err == nil {
		_, err = query.Exec(recipe.Calories, recipe.Carbs, recipe.Proteins, recipe.ID)
	}
	defer db.Close()
	return err
}

func sumRecipeNutrition(recipe *Recipe) (err error) {
	db := db.InitDB()
	result, err := db.Query("SELECT products.title, ingredients.value, products.calories, products.proteins, products.carbs FROM ingredients LEFT JOIN products ON ingredients.productId = products.id WHERE ingredients.recipeId = ?", recipe.ID)
	defer db.Close()
	if err == nil {
		for result.Next() {
			var product Product
			err := result.Scan(&product.Title, &product.Value, &product.Calories, &product.Proteins, &product.Carbs)
			if err != nil {
				return err
			}
			recipe.Calories += product.Calories
			recipe.Carbs += product.Carbs
			recipe.Proteins += product.Proteins
		}
	}
	return err
}

func createRecipeIngredients(recipe *Recipe) (err error) {
	db := db.InitDB()
	for _, product := range recipe.Products {
		query, err := db.Prepare("INSERT INTO ingredients(recipeid, productid, value) VALUES(?,?,?)")
		if err == nil {
			_, er := query.Exec(recipe.ID, product.ID, product.Value)
			defer db.Close()
			if er != nil {
				return err
			}
		}
	}
	defer db.Close()
	return err
}
