package backend

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
)

func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			var fileName = "./frontend/template/home/registration/index.html"
			t, err := template.ParseFiles(fileName)
			if err != nil {
				fmt.Println("Erreur pendant le parsing", err)
				return
			}
			t.Execute(w, nil)
		}
		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Erreur lors du traitement du formulaire", http.StatusBadRequest)
				return
			}

			username := r.FormValue("username")
			email := r.FormValue("email")
			password := r.FormValue("password")

			if username == "" || email == "" || password == "" {
				http.Error(w, "Tous les champs sont requis", http.StatusBadRequest)
				return
			}

			err = Register(db, username, email, password)

			if err != nil {
				http.Error(w, "Erreur lors de l'inscription: "+err.Error(), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	}
}

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			tmpl, err := template.ParseFiles("frontend/template/home/security/login.html")
			if err != nil {
				http.Error(w, "Erreur lors du chargement de la page de connexion", http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, nil)
			return
		}

		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Erreur lors du traitement du formulaire", http.StatusBadRequest)
				return
			}

			email := r.FormValue("email")
			password := r.FormValue("password")

			err = AuthenticateUser(db, email, password)
			if err != nil {
				fmt.Println("erreur a l'authentification", err)
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			// on genere un token de session et on le stocke dans un cookie
			token, err := GenerateSessionToken()
			if err != nil {
				http.Error(w, "Erreur lors de la generation du token", http.StatusInternalServerError)
				return
			}

			SetSessionCookie(w, token)
			fmt.Println("Utilisateur connecté :", email)

			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}