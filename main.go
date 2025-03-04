package main

import (
	"database/sql"
	"fmt"
	backend "forum/backend"
	"html/template"
	"net/http"
	"golang.org/x/crypto/bcrypt"
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
		
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Erreur lors du hashage du mot de passe", http.StatusInternalServerError)
			return
		}

		err = backend.InsertUser(db, username, email, string(hashedPassword))
		if err != nil {
			http.Error(w, "Erreur lors de l'insertion en base de données", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("frontend/template/home/security/login.html")
		if err != nil {
			http.Error(w, "Erreur lors du chargement de la page de connexion", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	} else if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Erreur lors du traitement du formulaire", http.StatusBadRequest)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		var storedPassword string
		err = db.QueryRow("SELECT password FROM user WHERE email = ?", email).Scan(&storedPassword)
		if err != nil {
			http.Error(w, "Utilisateur non trouvé", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
		if err != nil {
			http.Error(w, "Mot de passe incorrect", http.StatusUnauthorized)
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
	http.Handle("/", http.HandlerFunc(home))
	http.Handle("/register", http.HandlerFunc(register))
	http.Handle("/login", http.HandlerFunc(login))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	fmt.Println("Serveur démarré sur : http://localhost")
	http.ListenAndServe("", nil)
}