package main

import (
	"database/sql"
	"fmt"
	backend "forum/backend"
	"html/template"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	var fileName = "./frontend/template/home/index.html"
	t, err := template.ParseFiles(fileName)
	if err != nil {
		fmt.Println("Erreur pendant le parsing", err)
		return
	}

	t.Execute(w, nil)
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("frontend/template/home/registration/index.html")
		if err != nil {
			http.Error(w, "Erreur lors du chargement de la page", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	} else if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Erreur lors du traitement du formulaire", http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")



		if username == "" || email == "" || password == "" {
			http.Error(w, "Tous les champs doivent être remplis", http.StatusBadRequest)
			return
		}


		err = backend.InsertUser(db, username, email, password)
		if err != nil {
			http.Error(w, "Erreur lors de l'insertion en base de données", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

var db *sql.DB
func main() {
	db = backend.InitDB()
	if db == nil {
		fmt.Println("Impossible de démarrer l'application sans base de données")
		return
	}
	defer db.Close()


	fs := http.FileServer(http.Dir("./frontend/public/"))
	http.HandleFunc("/", home)
	http.HandleFunc("/register", register)
	http.Handle("/public/", http.StripPrefix("/public/", fs))
	http.Handle("frontend/public/js", http.StripPrefix("frontend/public/js", fs))
	http.ListenAndServe("", nil)
}