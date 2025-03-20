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
	// ajouter le reste de nos propriété
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

// GetUserByEmail récupère un utilisateur par son email
func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	user := &User{}
	err := db.QueryRow("SELECT id, username, email, password, session_token FROM user WHERE email = ?", email).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.SessionToken)
	
	if err == sql.ErrNoRows {
		return nil, errors.New("aucun utilisateur trouvé avec cet email")
	} else if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération de l'utilisateur : %v", err)
	}

	return user, nil
}