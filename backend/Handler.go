package backend

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	db, err := InitDB()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// recup les posts
	posts, err := GetPosts(db)
	if err != nil {
		http.Error(w, "Error fetching posts", http.StatusInternalServerError)
		return
	}

	// il faut parser les templates
	tmpl, err := template.ParseFiles(
		"./frontend/template/home/forum/accueil.html", 
		"./frontend/template/home/article/posts.html",
	)
	if err != nil {
		log.Printf("Template parsing error: %v", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	// on donne les données des posts a la template
	err = tmpl.ExecuteTemplate(w, "accueil", struct {
		Posts []Post
	}{
		Posts: posts,
	})
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Rendering error", http.StatusInternalServerError)
	}
}

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

		fmt.Println("Fonction CreatePost appelée")
		if !IsAuthenticated(r) {
			http.Error(w, "Vous devez être connecté pour créer un post", http.StatusUnauthorized)
			return
		}

		if r.Method == "GET" {
			fmt.Println("GET")
			tmpl, err := template.ParseFiles("frontend/template/home/article/add.html")
			if err != nil {
				http.Error(w, "Erreur lors du chargement du formulaire", http.StatusInternalServerError)
				fmt.Println("Erreur template :", err)
				return
			}
			tmpl.Execute(w, nil)
		}

		if r.Method == http.MethodPost {
			fmt.Println(" POST")
			err := r.ParseMultipartForm(20 << 20) // Limite a 20MB
			if err != nil {
				http.Error(w, "Erreur lors du traitement du formulaire", http.StatusBadRequest)
				return
			}

			title := r.FormValue("title")
			content := r.FormValue("content")
			sessionToken, _ := GetSessionToken(r)

			file, handler, err := r.FormFile("image")
			if err != nil {
				http.Error(w, "Erreur lors de l'upload de l'image", http.StatusBadRequest)
				return
			}
			defer file.Close()

			// il faut changer le chemin
			uploadDir := "uploads/"
			if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
				os.Mkdir(uploadDir, os.ModePerm)
			}

			imagePath := filepath.Join(uploadDir, handler.Filename)

			// save l'image
			dst, err := os.Create(imagePath)
			if err != nil {
				http.Error(w, "Erreur lors de la sauvegarde de l'image", http.StatusInternalServerError)
				return
			}
			defer dst.Close()
			io.Copy(dst, file)

			// recup id de l'utilisateur
			var userID int
			err = db.QueryRow("SELECT id FROM user WHERE session_token = ?", sessionToken).Scan(&userID)
			if err != nil {
				http.Error(w, "Utilisateur non trouvé", http.StatusUnauthorized)
				return
			}

			err = CreatePost(db, title, content, imagePath, userID)
			if err != nil {
				http.Error(w, "Erreur lors de la création du post", http.StatusInternalServerError)
				return
			}

		
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}
}


func PostDetailHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Extraire l'ID du post de l'URL
        postIDStr := r.URL.Path[len("/post/"):]
        postID, err := strconv.Atoi(postIDStr)
        if err != nil {
            http.Error(w, "ID de post invalide", http.StatusBadRequest)
            return
        }

        // recup les details du post
        post, err := GetPostByID(db, postID)
        if err != nil {
            if err == sql.ErrNoRows {
                http.Error(w, "Post non trouvé", http.StatusNotFound)
            } else {
                http.Error(w, "Erreur lors de la récupération du post", http.StatusInternalServerError)
            }
            return
        }

        // recup les commentaires lié au post
        comments, err := GetCommentsByPostID(db, postID)
        if err != nil {
            http.Error(w, "Erreur lors de la récupération des commentaires", http.StatusInternalServerError)
            return
        }
        
        // Récupérer le nombre de likes et dislikes
        likes, dislikes, err := CountLikes(db, postID)
        if err != nil {
            http.Error(w, "Erreur lors de la récupération des likes", http.StatusInternalServerError)
            return
        }

        // Charger le template
        tmpl, err := template.ParseFiles("./frontend/template/home/article/post.html")
        if err != nil {
            http.Error(w, "Erreur lors du chargement du template", http.StatusInternalServerError)
            return
        }

        // passer les données a la template
        err = tmpl.ExecuteTemplate(w, "post", struct {
            Post     *Post
            Comments []Comment
            Likes    int
            Dislikes int
        }{
            Post:     post,
            Comments: comments,
            Likes:    likes,
            Dislikes: dislikes,
        })
        if err != nil {
            http.Error(w, "Erreur lors du rendu du template", http.StatusInternalServerError)
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
            // Extraire l'ID du post de l'URL: /post/{ID}/comment
            pathParts := strings.Split(r.URL.Path, "/")
            if len(pathParts) < 3 {
                http.Error(w, "URL invalide", http.StatusBadRequest)
                return
            }
            postIDStr := pathParts[2]
            postID, err := strconv.Atoi(postIDStr)
            if err != nil {
                http.Error(w, "ID du post invalide", http.StatusBadRequest)
                return
            }

            sessionToken, _ := GetSessionToken(r)

            var userID int
            err = db.QueryRow("SELECT id FROM user WHERE session_token = ?", sessionToken).Scan(&userID)
            if err != nil {
                http.Error(w, "Utilisateur non trouvé", http.StatusUnauthorized)
                return
            }

            err = AddComment(db, content, userID, postID)
            if err != nil {
                http.Error(w, "Erreur lors de l'ajout du commentaire", http.StatusInternalServerError)
                return
            }

            // Rediriger vers la page du post
            http.Redirect(w, r, fmt.Sprintf("/post/%d", postID), http.StatusSeeOther)
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
            err = db.QueryRow("SELECT id FROM user WHERE session_token = ?", sessionToken).Scan(&userID)
            if err != nil {
                http.Error(w, "Utilisateur non trouvé", http.StatusUnauthorized)
                return
            }

            err = LikePost(db, userID, postID, likeType)
            if err != nil {
                http.Error(w, "Erreur lors de l'ajout du like/dislike", http.StatusInternalServerError)
                return
            }

            http.Redirect(w, r, fmt.Sprintf("/post/%d", postID), http.StatusSeeOther)
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
	return func(w http.ResponseWriter, r *http.Request) {
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

func GithubLoginHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        url := githubOAuthConfig.AuthCodeURL("randomstate")
        fmt.Println("Redirection vers :", url)
        http.Redirect(w, r, url, http.StatusTemporaryRedirect)
    }
}

func GithubCallbackHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        code := r.URL.Query().Get("code")

        token, err := githubOAuthConfig.Exchange(context.Background(), code)
        if err != nil {
            http.Error(w, "Erreur lors de l'échange du token", http.StatusInternalServerError)
            return
        }

        // Utilisation du client pour récupérer les infos de l'utilisateur via l'API GitHub
        client := githubOAuthConfig.Client(context.Background(), token)
        resp, err := client.Get("https://api.github.com/user")
        if err != nil {
            http.Error(w, "Erreur lors de la récupération des infos utilisateur", http.StatusInternalServerError)
            return
        }
        defer resp.Body.Close()

        userData, err := io.ReadAll(resp.Body)
        if err != nil {
            http.Error(w, "Erreur lors de la lecture des données utilisateur", http.StatusInternalServerError)
            return
        }
        fmt.Println("Données utilisateur GitHub :", string(userData))

        http.Redirect(w, r, "/", http.StatusSeeOther)
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

			if email == "" || password == "" {
				http.Error(w, "Tous les champs sont requis", http.StatusBadRequest)
				return
			}

			// Appel à la fonction Login
			token, err := Login(db, email, password)
			if err != nil {
				http.Error(w, "Erreur lors de la connexion: "+err.Error(), http.StatusUnauthorized)
				return
			}

			// Mettre le token dans le cookie
			SetSessionCookie(w, token)

			// Rediriger vers la page d'accueil ou autre
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}

func ProfileHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Vérifier d'abord si l'utilisateur est authentifié
		if !IsAuthenticated(r) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		
		// Récupérer le token depuis le cookie
		cookie, err := r.Cookie("session_token")
		if err != nil {
			fmt.Println("Erreur cookie:", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		
		// Vérifier que le token n'est pas vide
		if cookie.Value == "" {
			fmt.Println("Token vide dans le cookie")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Récupérer les informations de l'utilisateur
		user, err := GetUserInfoByToken(db, cookie.Value)
		if err != nil {
			fmt.Printf("Erreur GetUserInfoByToken: %v\n", err)
			http.Error(w, "Erreur lors de la récupération des informations de l'utilisateur", http.StatusInternalServerError)
			return
		}

		// Passer les informations de l'utilisateur à la vue
		tmpl, err := template.ParseFiles("frontend/template/home/profile/profil.html")
		if err != nil {
			fmt.Printf("Erreur parsing template: %v\n", err)
			http.Error(w, "Erreur lors du chargement de la page de profil", http.StatusInternalServerError)
			return
		}

		// Rendre la page avec les données utilisateur
		err = tmpl.Execute(w, user)
		if err != nil {
			fmt.Printf("Erreur exécution template: %v\n", err)
		}
	}
}

func GetCurrentUser(db *sql.DB, r *http.Request) *User {
    cookie, err := r.Cookie("session_token")
    if err != nil {
        return nil
    }

    user, err := GetUserInfoByToken(db, cookie.Value) 
    if err != nil {
        return nil
    }

    return user
}

func ActivityHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        user := GetCurrentUser(db, r)
        if user == nil {
            http.Error(w, "Utilisateur non authentifié", http.StatusUnauthorized)
            return
        }

        fmt.Fprintf(w, "Activités de l'utilisateur : %s", user.Username)
    }
}

func UpdateUserRoleHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        user := GetCurrentUser(db, r)
        if user == nil || !IsAdmin(user) {
            http.Error(w, "Accès refusé : seuls les admins peuvent modifier les rôles", http.StatusForbidden)
            return
        }

        // Récupérer les paramètres, par exemple via un formulaire
        targetUserIDStr := r.FormValue("user_id")
        newRole := r.FormValue("role")
        targetUserID, err := strconv.Atoi(targetUserIDStr)
        if err != nil {
            http.Error(w, "ID utilisateur invalide", http.StatusBadRequest)
            return
        }

        // Vérifier que le nouveau rôle est valide
        if newRole != "user" && newRole != "moderator" && newRole != "admin" {
            http.Error(w, "Rôle invalide", http.StatusBadRequest)
            return
        }

        // Mettre à jour le rôle dans la base de données
        _, err = db.Exec("UPDATE user SET role = ? WHERE id = ?", newRole, targetUserID)
        if err != nil {
            http.Error(w, "Erreur lors de la mise à jour du rôle", http.StatusInternalServerError)
            return
        }

        http.Redirect(w, r, "/admin", http.StatusSeeOther)
    }
}

func ReportPostHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        user := GetCurrentUser(db, r)
        
        if !IsModerator(user) {
            http.Error(w, "Accès refusé : seuls les modérateurs peuvent signaler du contenu", http.StatusForbidden)
            return
        }

        if r.Method == "POST" {
            err := r.ParseForm()
            if err != nil {
                http.Error(w, "Erreur lors du traitement du formulaire", http.StatusBadRequest)
                return
            }

            postIDStr := r.FormValue("post_id")
            reason := r.FormValue("reason")
            
            postID, err := strconv.Atoi(postIDStr)
            if err != nil {
                http.Error(w, "ID de post invalide", http.StatusBadRequest)
                return
            }

            err = CreatePostReport(db, user.ID, postID, reason)
            if err != nil {
                http.Error(w, "Erreur lors du signalement", http.StatusInternalServerError)
                return
            }

            http.Redirect(w, r, fmt.Sprintf("/post?id=%d", postID), http.StatusSeeOther)
        }
    }
}

func ModeratorDashboardHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        user := GetCurrentUser(db, r)
        
        if !IsModerator(user) {
            http.Error(w, "Accès refusé : seuls les modérateurs peuvent accéder à ce tableau de bord", http.StatusForbidden)
            return
        }

        if r.Method == "GET" {
            reports, err := GetPendingReports(db)
            if err != nil {
                http.Error(w, "Erreur lors de la récupération des signalements", http.StatusInternalServerError)
                return
            }

            tmpl, err := template.ParseFiles("frontend/template/moderation/dashboard.html")
            if err != nil {
                http.Error(w, "Erreur lors du chargement du tableau de bord", http.StatusInternalServerError)
                return
            }

            tmpl.Execute(w, reports)
        }
    }
}

func ResolveReportHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        user := GetCurrentUser(db, r)
        
        if !IsAdmin(user) {
            http.Error(w, "Accès refusé : seuls les administrateurs peuvent résoudre les signalements", http.StatusForbidden)
            return
        }

        if r.Method == "POST" {
            err := r.ParseForm()
            if err != nil {
                http.Error(w, "Erreur lors du traitement du formulaire", http.StatusBadRequest)
                return
            }

            reportIDStr := r.FormValue("report_id")
            actionStr := r.FormValue("action")
            
            reportID, err := strconv.Atoi(reportIDStr)
            if err != nil {
                http.Error(w, "ID de signalement invalide", http.StatusBadRequest)
                return
            }

            approve := (actionStr == "approve")

            err = ResolveReport(db, reportID, user.ID, approve)
            if err != nil {
                http.Error(w, "Erreur lors du traitement du signalement", http.StatusInternalServerError)
                return
            }

            http.Redirect(w, r, "/moderation/dashboard", http.StatusSeeOther)
        }
    }
}

func NotificationHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Récupérer le token de session
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Utilisateur non authentifié", http.StatusUnauthorized)
			return
		}

		// Récupérer les infos de l'utilisateur
		user, err := GetUserInfoByToken(db, cookie.Value)
		if err != nil {
			http.Error(w, "Erreur lors de la récupération de l'utilisateur", http.StatusInternalServerError)
			return
		}

		// Récupérer les notifications en utilisant la fonction existante GetUserNotifications
		notifications, err := GetUserNotifications(db, user.ID)
		if err != nil {
			http.Error(w, "Erreur lors de la récupération des notifications", http.StatusInternalServerError)
			return
		}

		// Passer les notifications à la vue
		tmpl, err := template.ParseFiles("frontend/template/home/notification/notifications.html")
		if err != nil {
			http.Error(w, "Erreur lors du chargement de la page des notifications", http.StatusInternalServerError)
			return
		}

		data := struct {
			User          User
			Notifications []Notification
		}{
			User:          *user,
			Notifications: notifications,
		}

		// Exécuter le template
		tmpl.Execute(w, data)
	}
}



