package main

import (
	"database/sql"
	"fmt"
	backend "forum/backend"
	"html/template"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
)

func home(w http.ResponseWriter, r *http.Request) {
	var fileName = "./frontend/template/home/forum/accueil.html"
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
	backend.OauthInit()
	errs := godotenv.Load()
	if errs != nil {
	  log.Fatal("Error loading .env file")
	}
	
	fmt.Println("Client ID:", os.Getenv("GOOGLE_CLIENT_ID"))
    fmt.Println("Client Secret:", os.Getenv("GOOGLE_CLIENT_SECRET"))
	backend.InitDB()

	fs := http.FileServer(http.Dir("./frontend/public/"))
    http.Handle("/", backend.LimitRequest(http.HandlerFunc(home)))
    http.Handle("/articles", backend.LimitRequest(http.HandlerFunc(backend.ArticlesHandler())))
    http.Handle("/login", backend.LimitRequest(http.HandlerFunc(backend.LoginHandler(db))))
    http.Handle("/register", backend.LimitRequest(http.HandlerFunc(backend.RegisterHandler(db))))
    http.Handle("/add", backend.LimitRequest(http.HandlerFunc(backend.ArticlesaddHandler(db))))
    http.Handle("/create_post", backend.LimitRequest(http.HandlerFunc(backend.CreatePostHandler(db))))
    http.Handle("/add_comment", backend.LimitRequest(http.HandlerFunc(backend.AddCommentHandler(db))))
    http.Handle("/like_dislike", backend.LimitRequest(http.HandlerFunc(backend.LikePostHandler(db))))
    http.Handle("/auth/google", backend.LimitRequest(http.HandlerFunc(backend.GoogleLoginHandler())))
    http.Handle("/auth/google/callback", backend.LimitRequest(http.HandlerFunc(backend.GoogleCallbackHandler(db))))
    http.Handle("/profile", backend.LimitRequest(http.HandlerFunc(backend.AdminHandler(db))))
    http.Handle("/upload", backend.LimitRequest(http.HandlerFunc(backend.UploadImage)))

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
