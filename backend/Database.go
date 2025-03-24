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
		db.Close()
		return nil, fmt.Errorf("❌ Impossible de contacter la BDD : %v", err)
	}

	// Exécute la création des tables séparément
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS user (
			id INT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			session_token VARCHAR(64) DEFAULT '' NOT NULL
		);
	`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("❌ Erreur lors de la création de la table user : %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS post (
			id INT AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			content TEXT NOT NULL,
			image VARCHAR(255) NOT NULL,
			author_id INT NOT NULL,
			FOREIGN KEY (author_id) REFERENCES user(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("❌ Erreur lors de la création de la table product : %v", err)
	}

	fmt.Println("✅ Connexion à MySQL réussie et tables créées !")
	return db, nil
}
