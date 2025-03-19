package backend

import (
	"database/sql"
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
	statement, err := db.Prepare("INSERT INTO user (username, email, password) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer statement.Close()
	
	_, err = statement.Exec(username, email, password)
	return err
}