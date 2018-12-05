package auth

import (
	"app/modules/db"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/gorilla/sessions"
)

type Admin struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
}

var Store = sessions.NewCookieStore([]byte("Tempkey-must-be-changed-in-prod"))

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/create", createAdmin)
	router.Post("/login", loginAdmin)
	return router
}

func createAdmin(w http.ResponseWriter, r *http.Request) {
	var admin Admin
	json.NewDecoder(r.Body).Decode(&admin)
	if !checkIfValid(admin) {
		http.Error(w, "Bad request", 400)
		return
	}
	db := db.InitDB()
	query, err := db.Prepare("INSERT INTO admins(name, pass) VALUES(?,?)")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	_, er := query.Exec(admin.Name, admin.Pass)
	defer db.Close()
	if er != nil {
		http.Error(w, "Username exists", 400)
		return
	}
	render.JSON(w, r, "Admin account was created")
}

func loginAdmin(w http.ResponseWriter, r *http.Request) {
	var admin Admin
	json.NewDecoder(r.Body).Decode(&admin)
	if !checkIfValid(admin) {
		http.Error(w, "Bad request", 400)
		return
	}
	db := db.InitDB()
	var ifExists bool
	err := db.QueryRow("SELECT IF(COUNT(*),'true','false') FROM admins WHERE name = ? AND pass = ?", admin.Name, admin.Pass).Scan(&ifExists)
	defer db.Close()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if ifExists {
		session, _ := Store.Get(r, "cookie")
		session.Values["auth"] = true
		session.Values["name"] = admin.Name
		session.Options = &sessions.Options{
			Path:   "/v1/api",
			MaxAge: 900}
		session.Save(r, w)
	}
	if !ifExists {
		render.JSON(w, r, "Incorrect login credencials")
	}

}

func checkIfValid(a Admin) bool {
	if len(a.Pass) < 15 {
		return false
	}
	if len(a.Name) == 0 {
		return false
	}
	return true
}
