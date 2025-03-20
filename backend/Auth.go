package backend

import (
	"crypto/rand"
	"encoding/hex"
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
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

	_, err = db.Exec("INSERT INTO user (username, email, password) VALUES (?, ?, ?)", username, email, hashedPassword)
	if err != nil {
		return err
	}

	fmt.Println("✅ Utilisateur enregistré avec succès:", username, "-", email)
	return nil
}

func AuthenticateUser(db *sql.DB, email, password string) (string, error) {
    if db == nil {
        return "", errors.New("connexion à la base de données invalide")
    }
    
    // Récupère l'utilisateur par email
    user, err := GetUserByEmail(db, email)
    if err != nil {
        return "", errors.New("utilisateur non trouvé")
    }
    
    // Vérifie le mot de passe
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        return "", errors.New("mot de passe incorrect")
    }
    
    // Génère un token de session
    token, err := GenerateSessionToken()
    if err != nil {
        return "", err
    }
    
    // Met à jour le token de session dans la base de données
    err = user.UpdateSessionToken(db, token)
    if err != nil {
        return "", err
    }
    
    // Retourne le token pour que le handler puisse créer le cookie
    return token, nil
}

func GenerateSessionToken() (string, error) {
	bytes := make([]byte, 32) // 64 caractères hexadecimaux
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
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