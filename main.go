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
	mux := http.NewServeMux() // Use a mux for clarity

	// --- Existing Routes (adapt to use mux) ---
	mux.Handle("/", backend.LimitRequest(http.HandlerFunc(backend.HomeHandler))) // Assuming HomeHandler is defined
    // Add other existing handlers to mux...
    mux.Handle("/articles", backend.LimitRequest(http.HandlerFunc(backend.ArticlesHandler())))
	mux.Handle("/login", backend.LimitRequest(http.HandlerFunc(backend.LoginHandler(db))))
	mux.Handle("/register", backend.LimitRequest(http.HandlerFunc(backend.RegisterHandler(db))))
	mux.Handle("/create_post", backend.LimitRequest(http.HandlerFunc(backend.CreatePostHandler(db))))
	// mux.Handle("/add_comment", backend.LimitRequest(http.HandlerFunc(backend.AddCommentHandler(db)))) // Covered by /post/{id}/comment now
	// mux.Handle("/like_dislike", backend.LimitRequest(http.HandlerFunc(backend.LikePostHandler(db)))) // Covered by /post/{id}/like now
	mux.Handle("/auth/google", backend.LimitRequest(http.HandlerFunc(backend.GoogleLoginHandler())))
	mux.Handle("/auth/google/callback", backend.LimitRequest(http.HandlerFunc(backend.GoogleCallbackHandler(db))))
	mux.Handle("/profile", backend.LimitRequest(http.HandlerFunc(backend.ProfileHandler(db))))
	mux.Handle("/report_post", backend.LimitRequest(http.HandlerFunc(backend.ReportPostHandler(db))))
	mux.Handle("/moderation/dashboard", backend.LimitRequest(http.HandlerFunc(backend.ModeratorDashboardHandler(db))))
	mux.Handle("/resolve_report", backend.LimitRequest(http.HandlerFunc(backend.ResolveReportHandler(db))))
	mux.Handle("/notification", backend.LimitRequest(http.HandlerFunc(backend.NotificationHandler(db))))


	// --- Dynamic Post Routes ---
	mux.HandleFunc("/post/", func(w http.ResponseWriter, r *http.Request) {
        path := strings.TrimSuffix(r.URL.Path, "/") // Trim trailing slash

        // Check for specific actions first
        if strings.HasSuffix(path, "/comment") {
            backend.AddCommentHandler(db)(w, r) // POST
        } else if strings.HasSuffix(path, "/like") {
             backend.LikePostHandler(db)(w, r) // POST
        } else if strings.HasSuffix(path, "/edit") {
            if r.Method == http.MethodGet {
                 backend.ShowEditPostFormHandler(db)(w, r) // GET - Show form
            } else if r.Method == http.MethodPost {
                 backend.HandleEditPostHandler(db)(w, r) // POST - Process update
            } else {
                 http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
            }
        } else if strings.HasSuffix(path, "/delete") {
             backend.DeletePostHandler(db)(w, r) // POST
        } else {
            // Default: Show post detail
            backend.PostDetailHandler(db)(w, r) // GET
        }
	})

    // --- Comment Delete Route ---
    mux.HandleFunc("/comment/", func(w http.ResponseWriter, r *http.Request) {
        path := strings.TrimSuffix(r.URL.Path, "/")
        if strings.HasSuffix(path, "/delete") {
            backend.DeleteCommentHandler(db)(w, r) // POST
        } else if strings.HasSuffix(path, "/edit") {
            // Add EditCommentHandler routes here if implementing
             http.Error(w, "Comment editing not yet implemented", http.StatusNotImplemented)
        } else {
            http.NotFound(w, r) // Or handle viewing a single comment if needed
        }
    })


	// --- Static Files ---
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))
	// Ensure correct paths for static assets
    mux.Handle("/frontend/public/", http.StripPrefix("/frontend/public/", fs))

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
	fmt.Println("🔹 Modération             : https://localhost/moderation/dashboard")
	fmt.Println("🔹 Signaler un post       : https://localhost/report_post")
	fmt.Println("🔹 notifications          : https://localhost/notification")

	log.Println("✅ Serveur HTTPS actif : https://localhost")
	http.Handle("/", mux) // Register the mux to handle all requests
    backend.StartSecureServer(nil) // Pass nil to use the DefaultServeMux
}