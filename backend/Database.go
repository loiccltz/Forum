package backend

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

// InitDB initialise la connexion à la base de données MySQL
func InitDB() *sql.DB {
	// AWS
	dsn := "admin:hardpassword@tcp(forum.cjoaea48gf89.eu-north-1.rds.amazonaws.com:3306)/forum"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("Erreur de connexion à MySQL :", err)
		return nil
	}

	// Vérifie que la connexion fonctionne
	err = db.Ping()
	if err != nil {
		fmt.Println("Impossible de contacter la bdd :", err)
		return nil
	}

	// Crée la table si elle n'existe pas
	statement, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS user (
			id INT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL,
			session_token VARCHAR(64) NOT NULL
		)
	`)
	if err != nil {
		fmt.Println("Erreur lors de la préparation de la requête :", err)
		return nil
	}
	statement.Exec()

	fmt.Println("Connexion à MySQL réussie !")
	return db
}
