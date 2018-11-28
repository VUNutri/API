package menu

import (
	"app/modules/db"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"time"

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
	router.Get("/getMenu/{daysCount}/{mealsCount}/{caloriesCount}", getMenu)
	return router
}

func getMenu(w http.ResponseWriter, r *http.Request) {
	daysCount, err := strconv.Atoi(chi.URLParam(r, "daysCount"))
	mealsCount, err := strconv.Atoi(chi.URLParam(r, "mealsCount"))
	caloriesCount, err := strconv.Atoi(chi.URLParam(r, "caloriesCount"))
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	breakfast, err := getRecipes(0)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	log.Print(breakfast)

	mainMeal, err := getRecipes(1)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	days := []Day{}

	if mealsCount == 1 {
		for i := 0; i < daysCount; i++ {
			day := Day{}
			day.Count = i + 1
			day.Recipes = append(day.Recipes, returnRand(prepRecipes(mainMeal, caloriesCount)))
			days = append(days, day)
		}
	} else if mealsCount == 2 {
		for i := 0; i < daysCount; i++ {
			day := Day{}
			day.Count = i + 1
			day.Recipes = append(day.Recipes, returnRand(prepRecipes(breakfast, caloriesCount)))
			day.Recipes = append(day.Recipes, returnRand(prepRecipes(mainMeal, caloriesCount)))
			days = append(days, day)
		}
	}
	render.JSON(w, r, days)
}

func getRecipes(cat int) ([]Recipe, error) {
	db := db.InitDB()
	result, err := db.Query("SELECT * FROM recipes WHERE category = ?", cat)
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

func prepRecipes(recipes []Recipe, cal int) []Recipe {
	i := 0
	for recipes[i].Calories <= cal {
		log.Println(len(recipes), i)
		i++
		if i == len(recipes) {
			break
		}
	}
	return recipes[0:i]
}
