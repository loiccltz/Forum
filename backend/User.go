package backend

import (
	"database/sql"
	"errors"
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

func GetUserInfoByToken(db *sql.DB, token string) (*User, error) {
	var user User
	err := db.QueryRow("SELECT id, email, username FROM user WHERE session_token = ?", token).Scan(&user.ID, &user.Email, &user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("utilisateur non trouvé")
		}
		return nil, err
	}
	return &user, nil
}
