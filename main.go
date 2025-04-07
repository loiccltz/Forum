package main

import (
	"database/sql"
	"fmt"
	backend "forum/backend"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var db *sql.DB

func main() {
	backend.OauthInit()
	errs := godotenv.Load()
	if errs != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("Client ID:", os.Getenv("GOOGLE_CLIENT_ID"))
	fmt.Println("Client Secret:", os.Getenv("GOOGLE_CLIENT_SECRET"))

	db, err := backend.InitDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	fs := http.FileServer(http.Dir("./frontend/public/"))

	http.Handle("/", backend.LimitRequest(http.HandlerFunc(backend.HomeHandler)))
	http.Handle("/articles", backend.LimitRequest(http.HandlerFunc(backend.ArticlesHandler())))
	http.Handle("/login", backend.LimitRequest(http.HandlerFunc(backend.LoginHandler(db))))
	http.Handle("/register", backend.LimitRequest(http.HandlerFunc(backend.RegisterHandler(db))))
	http.Handle("/create_post", backend.LimitRequest(http.HandlerFunc(backend.CreatePostHandler(db))))
	http.Handle("/add_comment", backend.LimitRequest(http.HandlerFunc(backend.AddCommentHandler(db))))
	http.Handle("/like_dislike", backend.LimitRequest(http.HandlerFunc(backend.LikePostHandler(db))))
	http.Handle("/auth/google", backend.LimitRequest(http.HandlerFunc(backend.GoogleLoginHandler())))
	http.Handle("/auth/google/callback", backend.LimitRequest(http.HandlerFunc(backend.GoogleCallbackHandler(db))))
	http.Handle("/profile", backend.LimitRequest(http.HandlerFunc(backend.ProfileHandler(db))))
	// Mise à jour du routage
	http.HandleFunc("/post/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Route pour les commentaires
		if strings.HasSuffix(path, "/comment") {
			backend.AddCommentHandler(db)(w, r)
			return
		}

		// Route pour les likes
		if strings.HasSuffix(path, "/like") {
			backend.LikePostHandler(db)(w, r)
			return
		}

		// Affichage du post
		backend.PostDetailHandler(db)(w, r)
	})

	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))
	http.Handle("/public/", http.StripPrefix("/public/", fs))
	http.Handle("frontend/public/js", http.StripPrefix("frontend/public/js", fs))

	fmt.Println("\n📌 Pages disponibles :")
	fmt.Println("🔹 Page d'accueil         : https://localhost/")
	fmt.Println("🔹 Page d'inscription     : https://localhost/register")
	fmt.Println("🔹 Page de connexion      : https://localhost/login")
	fmt.Println("🔹 Ajouter un article     : https://localhost/add")
	fmt.Println("🔹 Voir les articles      : https://localhost/articles")
	fmt.Println("🔹 Création de post       : https://localhost/create_post")
	fmt.Println("🔹 Ajouter une image      : https://localhost/upload")
	fmt.Println("🔹 Ajouter un commentaire : https://localhost/add_comment")
	fmt.Println("🔹 Like/Dislike un post   : https://localhost/like_dislike")
	fmt.Println("🔹 Profil utilisateur     : https://localhost/profile")

	log.Println("✅ Serveur HTTPS actif : https://localhost")
	backend.StartSecureServer(http.DefaultServeMux)
}
