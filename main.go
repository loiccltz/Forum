package main

import (
	"database/sql"
	"fmt"
	backend "forum/backend"
	"html/template"
	"log"
	"net/http"
	"os"
)

func home(w http.ResponseWriter, r *http.Request) {
	var fileName = "./frontend/template/home/index.html"
	t, err := template.ParseFiles(fileName)
	if err != nil {
		fmt.Println("Erreur pendant le parsing", err)
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
		return
	}

	t.Execute(w, nil)
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./forum.db")
	fmt.Println("Client ID:", os.Getenv("GOOGLE_CLIENT_ID"))
    fmt.Println("Client Secret:", os.Getenv("GOOGLE_CLIENT_SECRET"))
	if err != nil {
		log.Fatal("Erreur lors de l'ouverture de la base de données:", err)
	}
	backend.InitDB(db)

	fs := http.FileServer(http.Dir("./frontend/public/"))
	http.Handle("/", backend.LimitRequest(http.HandlerFunc(home)))
	http.HandleFunc("/articles", backend.ArticlesHandler())
	http.HandleFunc("/login", backend.LoginHandler(db))
	http.HandleFunc("/register", backend.RegisterHandler(db))
	http.HandleFunc("/add", backend.ArticlesaddHandler(db))
	http.HandleFunc("/create_post", backend.CreatePostHandler(db))
	http.HandleFunc("/add_comment", backend.AddCommentHandler(db))
	http.HandleFunc("/like_dislike", backend.LikePostHandler(db))
	http.HandleFunc("/auth/google", backend.GoogleLoginHandler())
	http.HandleFunc("/auth/google/callback", backend.GoogleCallbackHandler(db))
	http.HandleFunc("/profile", backend.AdminHandler(db))
	http.Handle("/public/", http.StripPrefix("/public/", fs))
	http.Handle("frontend/public/js", http.StripPrefix("frontend/public/js", fs))

	fmt.Println("\n📌 Pages disponibles :")
	fmt.Println("🔹 Page d'accueil         : https://localhost/")
	fmt.Println("🔹 Page d'inscription     : https://localhost/register")
	fmt.Println("🔹 Page de connexion      : https://localhost/login")
	fmt.Println("🔹 Ajouter un article     : https://localhost/add")
	fmt.Println("🔹 Voir les articles      : https://localhost/articles")
	fmt.Println("🔹 Création de post       : https://localhost/create_post")
	fmt.Println("🔹 Ajouter un commentaire : https://localhost/add_comment")
	fmt.Println("🔹 Like/Dislike un post   : https://localhost/like_dislike")
	fmt.Println("🔹 Profil utilisateur     : https://localhost/profile")

	log.Println("✅ Serveur HTTPS actif : https://localhost")
	backend.StartSecureServer(http.DefaultServeMux)
}
