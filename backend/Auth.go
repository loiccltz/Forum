package backend

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Register(db *sql.DB, username, email, password string) error {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM user WHERE email = ?)", email).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("cet email est déjà utilisé")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("erreur lors du hashage du mot de passe")
	}

	// Générer un token de session
	token, err := GenerateSessionToken()
	if err != nil {
		return errors.New("erreur lors de la génération du token")
	}

	// Insérer l'utilisateur avec le token
	_, err = db.Exec("INSERT INTO user (username, email, password, session_token) VALUES (?, ?, ?, ?)", 
		username, email, hashedPassword, token)
	if err != nil {
		return err
	}

	fmt.Println("✅ Utilisateur enregistré avec succès:", username, "-", email)
	return nil
}

func Login(db *sql.DB, email, password string) (string, error) {
	// Vérifier si l'email existe
	var hashedPassword string
	err := db.QueryRow("SELECT password FROM user WHERE email = ?", email).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("email ou mot de passe incorrect")
		}
		return "", err
	}

	// Comparer le mot de passe avec le hash stocké
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return "", errors.New("email ou mot de passe incorrect")
	}

	// Générer un nouveau token de session
	token, err := GenerateSessionToken()
	if err != nil {
		return "", errors.New("erreur lors de la génération du token")
	}

	// Mettre à jour le token de session dans la base de données
	_, err = db.Exec("UPDATE user SET session_token = ? WHERE email = ?", token, email)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Génération d'un token de session avec UUID
func GenerateSessionToken() (string, error) {
	token := uuid.New().String() // Génère un UUID
	return token, nil
}

func StoreSessionToken(db *sql.DB, email, token string) error {
    _, err := db.Exec("UPDATE user SET session_token = ? WHERE email = ?", token, email)
    return err
}


func SetSessionCookie(w http.ResponseWriter, token string) {
	cookie := http.Cookie{
		Name:     "session_token",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   3600, // 1 heure
	}
	http.SetCookie(w, &cookie)
}


// recupere le token de session stocké dans le cookie
func GetSessionToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}


// utilisé pour verifier si un utilisateur est connecté (pour restreindre l'acces a certaines parges)
func IsAuthenticated(r *http.Request) bool {
	token, err := GetSessionToken(r)
	return err == nil && token != ""
}