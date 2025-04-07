package backend

import (
	"database/sql"
	"errors"
	"fmt"
)


type User struct {
	ID           int
	Username     string
	Email        string
	Password     string // hashé
	SessionToken string
	Role         string
}

//insere un nouvel utilisateur dans la base de donnees
func InsertUser(db *sql.DB, username, email, password string) error {
	statement, err := db.Prepare("INSERT INTO user (username, email, password) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer statement.Close()
	
	_, err = statement.Exec(username, email, password)
	return err
}

//met a jour le token de session d'un utilisateur
func (u *User) UpdateSessionToken(db *sql.DB, token string) error {
	_, err := db.Exec("UPDATE user SET session_token = ? WHERE id = ?", token, u.ID)
	if err != nil {
		return err
	}
	u.SessionToken = token
	return nil
}


func GetUserInfoByToken(db *sql.DB, token string) (*User, error) {
	var user User
	
	// Debug: imprimer le token pour vérification
	fmt.Printf("Recherche de l'utilisateur avec le token: %s\n", token)
	
	// Vérifions que le token n'est pas vide
	if token == "" {
		return nil, errors.New("token de session vide")
	}
	
	err := db.QueryRow("SELECT id, email, username, role FROM user WHERE session_token = ?", token).Scan(&user.ID, &user.Email, &user.Username, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("Aucun utilisateur trouvé avec ce token: %s\n", token)
			return nil, errors.New("utilisateur non trouvé")
		}
		fmt.Printf("Erreur SQL: %v\n", err)
		return nil, err
	}
	
	fmt.Printf("Utilisateur trouvé: ID=%d, Username=%s\n", user.ID, user.Username)
	return &user, nil
}