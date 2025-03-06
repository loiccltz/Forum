package main

import (
	"database/sql"
	"fmt"
	backend "forum/backend"
	"html/template"
	"log"
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

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal("Erreur lors de l'ouverture de la base de données:", err)
	}

	backend.InitDB(db)


	fs := http.FileServer(http.Dir("./frontend/public/"))
	http.HandleFunc("/", home)
	http.HandleFunc("/articles", backend.ArticlesHandler())
	http.HandleFunc("/login", backend.LoginHandler(db))
	http.HandleFunc("/register", backend.RegisterHandler(db))
	http.Handle("/public/", http.StripPrefix("/public/", fs))
	http.Handle("frontend/public/js", http.StripPrefix("frontend/public/js", fs))
	fmt.Println("Serveur démarré sur : http://localhost")
	http.ListenAndServe("", nil)
}