package backend

import (
	"database/sql"
	"fmt"
)

func InsertUser(db *sql.DB, username, email, password string) error {
	statement, err := db.Prepare("INSERT INTO user (username, email, password) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Println("Erreur de préparation de la requête :", err)
		return err
	}
	_, err = statement.Exec(username, email, password)
	if err != nil {
		fmt.Println("Erreur d'exécution de la requête :", err)
	}
	return err
}