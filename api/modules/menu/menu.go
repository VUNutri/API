package menu

import (
	"app/modules/db"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Day struct {
	Count   int      `json:"dayCount"`
	Recipes []Recipe `json:"meals"`
}

type DayMenu struct {
	Meals    int      `json:"meals"`
	Calories int      `json:"calories"`
	Time     int      `json:"time"`
	Block    []string `json:"blockedIngredients"`
	Count    int      `json:"dayCount"`
}

type SingleMenu struct {
	Category int      `json:"category"`
	Calories int      `json:"calories"`
	Time     int      `json:"time"`
	Block    []string `json:"blockedIngredients"`
}

type Menu struct {
	Days     int      `json:"days"`
	Meals    int      `json:"meals"`
	Calories int      `json:"calories"`
	Time     int      `json:"time"`
	Block    []string `json:"blockedIngredients"`
}

type Product struct {
	ID       int     `json:"id"`
	Title    string  `json:"title"`
	Value    float64 `json:"value"`
	Size     string  `json:"size"`
	Calories int     `json:"calories"`
	Carbs    int     `json:"carbs"`
	Proteins int     `json:"proteins"`
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
	RecipeID  int     `json:"recipeId"`
	ProductID int     `json:"productId"`
	Value     float64 `json:"value"`
}

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/getMenu", getMenu)
	router.Post("/getDailyMenu", getDailyMenu)
	router.Post("/getDayOneMenu", getOneDayMenu)
	return router
}

func getMenu(w http.ResponseWriter, r *http.Request) {
	var menu Menu

	json.NewDecoder(r.Body).Decode(&menu)

	if !checkIfValid(menu) {
		http.Error(w, "Ivalid request", 400)
		return
	}

	breakfast, err := getRecipes(menu.Block, 1, menu.Calories/4, menu.Time)
	if err != nil {
		http.Error(w, "Bad request", 400)
		return
	}

	mainMeal, err := getRecipes(menu.Block, 2, menu.Calories/2, menu.Time)
	if err != nil {
		http.Error(w, "Bad request", 400)
		return
	}

	snacks, err := getRecipes(menu.Block, 3, menu.Calories/2, menu.Time)
	if err != nil {
		http.Error(w, "Bad request", 400)
		return
	}

	days := []Day{}

	if menu.Meals == 1 {
		if len(mainMeal) < 1 {
			http.Error(w, "Bad requirements", 400)
			return
		}
		for i := 0; i < menu.Days; i++ {
			day := Day{}
			day.Count = i + 1
			day.Recipes = append(day.Recipes, returnRand(mainMeal))
			days = append(days, day)
		}
	} else if menu.Meals == 2 {
		if len(mainMeal) < 1 {
			http.Error(w, "Bad requirements", 400)
			return
		} else if len(breakfast) < 1 {
			http.Error(w, "Bad requirements", 400)
			return
		}
		for i := 0; i < menu.Days; i++ {
			day := Day{}
			day.Count = i + 1
			day.Recipes = append(day.Recipes, returnRand(breakfast))
			day.Recipes = append(day.Recipes, returnRand(mainMeal))
			days = append(days, day)
		}
	} else if menu.Meals == 3 {
		if len(mainMeal) < 1 {
			http.Error(w, "Bad requirements", 400)
			return
		} else if len(breakfast) < 1 {
			http.Error(w, "Bad requirements", 400)
			return
		}
		for i := 0; i < menu.Days; i++ {
			day := Day{}
			day.Count = i + 1
			day.Recipes = append(day.Recipes, returnRand(breakfast))
			day.Recipes = append(day.Recipes, returnRand(mainMeal))
			day.Recipes = append(day.Recipes, returnRand(mainMeal))
			days = append(days, day)
		}
	} else if menu.Meals > 3 {
		if len(mainMeal) < 1 {
			http.Error(w, "Bad requirements", 400)
			return
		} else if len(breakfast) < 1 {
			http.Error(w, "Bad requirements", 400)
			return
		} else if len(snacks) < 1 {
			http.Error(w, "Bad requirements", 400)
			return
		}
		for i := 0; i < menu.Days; i++ {
			day := Day{}
			day.Count = i + 1
			day.Recipes = append(day.Recipes, returnRand(breakfast))
			day.Recipes = append(day.Recipes, returnRand(mainMeal))
			day.Recipes = append(day.Recipes, returnRand(mainMeal))
			for i := 0; i < menu.Meals-3; i++ {
				day.Recipes = append(day.Recipes, returnRand(snacks))
			}
			days = append(days, day)
		}

	}
	render.JSON(w, r, days)
	return
}

func getDailyMenu(w http.ResponseWriter, r *http.Request) {
	var menu DayMenu
	json.NewDecoder(r.Body).Decode(&menu)

	if !checkIfValidDMenu(menu) {
		http.Error(w, "Ivalid requests", 400)
		return
	}

	breakfast, err := getRecipes(menu.Block, 1, menu.Calories/4, menu.Time)
	if err != nil {
		http.Error(w, "Bad request", 400)
		return
	}

	mainMeal, err := getRecipes(menu.Block, 2, menu.Calories/2, menu.Time)
	if err != nil {
		http.Error(w, "Bad request", 400)
		return
	}

	snacks, err := getRecipes(menu.Block, 3, menu.Calories/2, menu.Time)
	if err != nil {
		http.Error(w, "Bad request", 400)
		return
	}

	day := Day{}
	day.Count = menu.Count

	if menu.Meals == 1 {
		if len(mainMeal) < 1 {
			http.Error(w, "Bad requirements", 400)
			return
		}
		day.Recipes = append(day.Recipes, returnRand(mainMeal))
	} else if menu.Meals == 2 {
		if len(mainMeal) < 1 {
			http.Error(w, "Bad requirements", 400)
			return
		} else if len(breakfast) < 1 {
			http.Error(w, "Bad requirements", 400)
			return
		}
		day.Recipes = append(day.Recipes, returnRand(breakfast))
		day.Recipes = append(day.Recipes, returnRand(mainMeal))
	} else if menu.Meals == 3 {
		if len(mainMeal) < 1 {
			http.Error(w, "Bad requirements", 400)
			return
		} else if len(breakfast) < 1 {
			http.Error(w, "Bad requirements", 400)
			return
		}
		day.Recipes = append(day.Recipes, returnRand(breakfast))
		day.Recipes = append(day.Recipes, returnRand(mainMeal))
		day.Recipes = append(day.Recipes, returnRand(mainMeal))
	} else if menu.Meals > 3 {
		if len(mainMeal) < 1 {
			http.Error(w, "Bad requirements", 400)
			return
		} else if len(breakfast) < 1 {
			http.Error(w, "Bad requirements", 400)
			return
		} else if len(snacks) < 1 {
			http.Error(w, "Bad requirements", 400)
			return
		}
		day.Recipes = append(day.Recipes, returnRand(breakfast))
		day.Recipes = append(day.Recipes, returnRand(mainMeal))
		day.Recipes = append(day.Recipes, returnRand(mainMeal))
		for i := 0; i < menu.Meals-3; i++ {
			day.Recipes = append(day.Recipes, returnRand(snacks))
		}
	}
	render.JSON(w, r, day)
	return
}

func getOneDayMenu(w http.ResponseWriter, r *http.Request) {
	var menu SingleMenu
	json.NewDecoder(r.Body).Decode(&menu)

	if !checkIfValidDOneMenu(menu) {
		http.Error(w, "Ivalid requests", 400)
		return
	}

	if menu.Category == 1 {
		breakfast, err := getRecipes(menu.Block, 1, menu.Calories/4, menu.Time)
		if err != nil {
			http.Error(w, "Bad request", 400)
			return
		}
		render.JSON(w, r, returnRand(breakfast))
	} else if menu.Category == 2 {
		mainMeal, err := getRecipes(menu.Block, 2, menu.Calories/2, menu.Time)
		if err != nil {
			http.Error(w, "Bad request", 400)
			return
		}
		render.JSON(w, r, returnRand(mainMeal))
	} else if menu.Category == 3 {
		snacks, err := getRecipes(menu.Block, 3, menu.Calories/2, menu.Time)
		if err != nil {
			http.Error(w, "Bad request", 400)
			return
		}
		render.JSON(w, r, returnRand(snacks))
	}

	http.Error(w, "Bad request", 400)
	return
}

func getRecipes(products []string, cat int, calories int, time int) ([]Recipe, error) {
	db := db.InitDB()

	ids := strings.Join(products, "','")
	if len(ids) < 1 {
		ids = "0"
	}

	quer := fmt.Sprintf("SELECT A.id, A.title, A.category, A.time, A.image, A.instructions, A.calories, A.carbs, A.proteins FROM recipes A INNER JOIN (SELECT recipeId FROM ingredients GROUP BY recipeId HAVING SUM(CASE WHEN productId IN (%s) THEN 1 ELSE 0 END) = 0) B ON A.id = B.recipeId AND A.category = %d AND A.calories <= %d AND A.time <= %d", ids, cat, calories, time)

	result, err := db.Query(quer)
	if err != nil {
		log.Println(err)
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
		result, er := db.Query("SELECT products.id, products.title, ingredients.value, products.calories, products.proteins, products.carbs, products.size FROM ingredients LEFT JOIN products ON ingredients.productId = products.id WHERE ingredients.recipeId = ?", recipe.ID)
		if er != nil {
			return nil, err
		}
		for result.Next() {
			var product Product
			err := result.Scan(&product.ID, &product.Title, &product.Value, &product.Calories, &product.Proteins, &product.Carbs, &product.Size)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			recipes[idx].Products = append(recipes[idx].Products, product)
		}
	}
	defer db.Close()

	return recipes, nil
}

func returnRand(recipes []Recipe) Recipe {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(recipes))
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
	if menu.Time == 0 {
		return false
	}
	return true
}

func checkIfValidDMenu(menu DayMenu) bool {
	if menu.Meals == 0 {
		return false
	}
	if menu.Calories == 0 {
		return false
	}
	if menu.Time == 0 {
		return false
	}
	if menu.Count == 0 {
		return false
	}
	return true
}

func checkIfValidDOneMenu(menu SingleMenu) bool {
	if menu.Calories == 0 {
		return false
	}
	if menu.Time == 0 {
		return false
	}
	if menu.Category == 0 {
		return false
	}
	return true
}
