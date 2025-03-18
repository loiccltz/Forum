package backend

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
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

func ArticlesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			tmpl, err := template.ParseFiles("frontend/template/home/article/index.html")
			if err != nil {
				http.Error(w, "Erreur lors du chargement de la page de connexion", http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, nil)
			return
		}
	}
}

func CreatePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !IsAuthenticated(r) {
			http.Error(w, "Vous devez être connecté pour créer un post", http.StatusUnauthorized)
			return
		}

		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Erreur lors du traitement du formulaire", http.StatusBadRequest)
				return
			}

			title := r.FormValue("title")
			content := r.FormValue("content")
			sessionToken, _ := GetSessionToken(r)
			imageURL := r.FormValue("image_url")

			var userID int
			err = db.QueryRow("SELECT id FROM user WHERE id = (SELECT user_id FROM sessions WHERE token = ?)", sessionToken).Scan(&userID)
			if err != nil {
				http.Error(w, "Utilisateur non trouvé", http.StatusUnauthorized)
				return
			}

			err = CreatePost(db, title, content, imageURL, userID)
			if err != nil {
				http.Error(w, "Erreur lors de la création du post", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}

func AddCommentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !IsAuthenticated(r) {
			http.Error(w, "Vous devez être connecté pour commenter", http.StatusUnauthorized)
			return
		}

		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Erreur lors du traitement du formulaire", http.StatusBadRequest)
				return
			}

			content := r.FormValue("content")
			postIDStr := r.FormValue("post_id")
			postID, err := strconv.Atoi(postIDStr)
			if err != nil {
				http.Error(w, "ID du post invalide", http.StatusBadRequest)
				return
			}

			sessionToken, _ := GetSessionToken(r)

			var userID int
			err = db.QueryRow("SELECT id FROM user WHERE id = (SELECT user_id FROM sessions WHERE token = ?)", sessionToken).Scan(&userID)
			if err != nil {
				http.Error(w, "Utilisateur non trouvé", http.StatusUnauthorized)
				return
			}

			err = AddComment(db, content, userID, postID)
			if err != nil {
				http.Error(w, "Erreur lors de l'ajout du commentaire", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, fmt.Sprintf("/post?id=%d", postID), http.StatusSeeOther)
		}
	}
}

func LikePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !IsAuthenticated(r) {
			http.Error(w, "Vous devez être connecté pour liker un post", http.StatusUnauthorized)
			return
		}

		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Erreur lors du traitement du formulaire", http.StatusBadRequest)
				return
			}

			postIDStr := r.FormValue("post_id_like")
			likeTypeStr := r.FormValue("like_type")
			postID, err := strconv.Atoi(postIDStr)
			if err != nil {
				http.Error(w, "ID du post invalide", http.StatusBadRequest)
				return
			}
			likeType, err := strconv.Atoi(likeTypeStr)
			if err != nil || (likeType != 0 && likeType != 1) {
				http.Error(w, "Type de like invalide", http.StatusBadRequest)
				return
			}

			sessionToken, _ := GetSessionToken(r)

			var userID int
			err = db.QueryRow("SELECT id FROM user WHERE id = (SELECT user_id FROM sessions WHERE token = ?)", sessionToken).Scan(&userID)
			if err != nil {
				http.Error(w, "Utilisateur non trouvé", http.StatusUnauthorized)
				return
			}

			err = LikePost(db, userID, postID, likeType)
			if err != nil {
				http.Error(w, "Erreur lors de l'ajout du like/dislike", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, fmt.Sprintf("/post?id=%d", postID), http.StatusSeeOther)
		}
	}
}

func ArticlesaddHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			tmpl, err := template.ParseFiles("frontend/template/home/article/add.html")
			if err != nil {
				http.Error(w, "Erreur lors du chargement de la page de connexion", http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, nil)
			return
		}
	}
}

func AdminHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			tmpl, err := template.ParseFiles("frontend/template/home/profile/profil.html")
			if err != nil {
				http.Error(w, "Erreur lors du chargement de la page de connexion", http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, nil)
			return
		}
	}
}

func GoogleLoginHandler() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request)  {
		url := googleOAuthConfig.AuthCodeURL("randomstate")
		fmt.Println(url)
		// randomstate c'est pour prevenir des attaque de type CSRF
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	} 
}

func GoogleCallbackHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// recup le code d'autorisation 
		code := r.URL.Query().Get("code")

		token, err := googleOAuthConfig.Exchange(context.Background(), code)
		// echange du code recup plus haut avec un token
		if err != nil {
			http.Error(w, "Erreur lors de l'échange du token", http.StatusInternalServerError)
			return
		}

		// on recupere les info de l'user avec le token
		client := googleOAuthConfig.Client(context.Background(), token)
		resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
		if err != nil {
			http.Error(w, "Erreur lors de la récupération des infos utilisateur", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close() // fermer la requete ( bonne pratique )

		userData, _ := io.ReadAll(resp.Body)
		fmt.Println("Données utilisateur :", string(userData))

		
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
