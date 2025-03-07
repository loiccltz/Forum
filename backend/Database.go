package backend

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)


func InitDB(db *sql.DB) *sql.DB {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		fmt.Println("Erreur de connexion à la base de données :", err)
		return nil
	}

	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS user (id INTEGER PRIMARY KEY, username TEXT, email TEXT, password TEXT)")
	if err != nil {
		fmt.Println("Erreur lors de la préparation de la requête :", err)
		return nil
	}
	statement.Exec()
	
	return db
}