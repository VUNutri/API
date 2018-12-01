package menu

import (
	"app/modules/db"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"time"
	"encoding/json"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Day struct {
	Count   int      `json:"dayCount"`
	Recipes []Recipe `json:"meals"`
}

type Menu struct {
	Days int `json:"days"`
	Meals int `json:"meals"`
	Calories int `json:"calories"`
	Block []int `json:"blockedIngredients"`
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
	router.Post("/getMenu", getMenu)
	return router
}

func getMenu(w http.ResponseWriter, r *http.Request) {
	var menu Menu
	
	json.NewDecoder(r.Body).Decode(&menu)

	if !checkIfValid(menu) {
		http.Error(w, "Bad request", 400)
		return
	}

	breakfast, err := getRecipes(0, menu.Calories / 4)
	if err != nil {
		http.Error(w, "Bad request", 400)
		return
	}

	mainMeal, err := getRecipes(1, menu.Calories / 2)
	if err != nil {
		http.Error(w, "Bad request", 400)
		return
	}

	days := []Day{}

	if menu.Meals == 1 {
		for i := 0; i < menu.Days; i++ {
			day := Day{}
			day.Count = i + 1
			day.Recipes = append(day.Recipes, returnRand(mainMeal))
			days = append(days, day)
		}
	} else if menu.Meals == 2 {
		for i := 0; i < menu.Days; i++ {
			day := Day{}
			day.Count = i + 1
			day.Recipes = append(day.Recipes, returnRand(breakfast))
			day.Recipes = append(day.Recipes, returnRand(mainMeal))
			days = append(days, day)
		}
	}
	render.JSON(w, r, days)
	return
}

func getRecipes(cat int, calories int) ([]Recipe, error) {
	db := db.InitDB()
	result, err := db.Query("SELECT * FROM recipes WHERE category = ? AND calories <= ?", cat, calories)
	if err != nil {
		return nil, err
	}

	recipes := []Recipe{}
	for result.Next() {
		var recipe Recipe
		err := result.Scan(&recipe.ID, &recipe.Title, &recipe.Category, &recipe.Time, &recipe.Image, &recipe.Instructions, &recipe.Calories, &recipe.Carbs, &recipe.Proteins)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, recipe)
	}

	for idx, recipe := range recipes {
		result, er := db.Query("SELECT products.id, products.title, ingredients.value, products.calories, products.proteins, products.carbs FROM ingredients LEFT JOIN products ON ingredients.productId = products.id WHERE ingredients.recipeId = ?", recipe.ID)
		if er != nil {
			return nil, err
		}
		for result.Next() {
			var product Product
			err := result.Scan(&product.ID, &product.Title, &product.Value, &product.Calories, &product.Proteins, &product.Carbs)
			if err != nil {
				return nil, err
			}
			recipes[idx].Products = append(recipes[idx].Products, product)
		}
	}
	defer db.Close()

	sort.SliceStable(recipes, func(i, j int) bool {
		return recipes[i].Calories < recipes[j].Calories
	})

	return recipes, nil
}

func returnRand(recipes []Recipe) Recipe {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(recipes))
	log.Println(len(recipes))
	return recipes[n]
}

func checkIfValid(menu Menu) bool {
	if menu.Days == 0 {
		return false
	}
	if menu.Meals == 0 {
		return false
	}
	if menu.Calories == 0 {
		return false
	}
	return true
}