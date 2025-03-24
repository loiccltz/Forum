package backend

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

// InitDB initialise la connexion à la base de données MySQL
func InitDB() (*sql.DB, error) {
	// AWS
	dsn := "admin:hardpassword@tcp(forum.cjoaea48gf89.eu-north-1.rds.amazonaws.com:3306)/forum"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("❌ Erreur de connexion à MySQL : %v", err)
	}

	// Vérifie que la connexion fonctionne
	err = db.Ping()
	if err != nil {
		db.Close() // Ferme la connexion si elle est inutilisable
		return nil, fmt.Errorf("❌ Impossible de contacter la BDD : %v", err)
	}

	// Crée la table si elle n'existe pas
	statement, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS user (
			id INT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			session_token VARCHAR(64) DEFAULT '' NOT NULL
		)
	`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("❌ Erreur lors de la préparation de la requête SQL : %v", err)
	}
	defer statement.Close()

	_, err = statement.Exec()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("❌ Erreur lors de l'exécution de la création de table : %v", err)
	}

	fmt.Println("✅ Connexion à MySQL réussie !")
	return db, nil
}
