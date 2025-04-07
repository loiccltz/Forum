package backend

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

// InitDB initialise la connexion à la base de données MySQL
func InitDB() (*sql.DB, error) {

	dsn := "root:Test@tcp(127.0.0.1:3306)/forum"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf(" Erreur de connexion à MySQL : %v", err)
	}

	// Vérifie que la connexion fonctionne
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf(" Impossible de contacter la BDD : %v", err)
	}

	// Exécute la création des tables séparément
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS user (
			id INT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			session_token VARCHAR(64) DEFAULT '' NOT NULL,
			role ENUM('user', 'moderator', 'admin') NOT NULL DEFAULT 'user'
		);
	`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf(" Erreur lors de la création de la table user : %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS post (
			id INT AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			content TEXT NOT NULL,
			image_url VARCHAR(255) NOT NULL,
			author_id INT NOT NULL,
			FOREIGN KEY (author_id) REFERENCES user(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		db.Close()

		return nil, fmt.Errorf("❌ Erreur lors de la création de la table post : %v", err)
	}

	_, err = db.Exec(`
	    CREATE TABLE IF NOT EXISTS category (
        	id INT AUTO_INCREMENT PRIMARY KEY,
        	name VARCHAR(255) NOT NULL UNIQUE
    	);
	`)
		if err != nil {
		db.Close()
		return nil, fmt.Errorf("❌ Erreur lors de la création de la table category : %v", err)
	}

	_, err = db.Exec(`
    	CREATE TABLE IF NOT EXISTS post_category (
        	post_id INT,
        	category_id INT,
        	PRIMARY KEY (post_id, category_id),
        	FOREIGN KEY (post_id) REFERENCES post(id),
        	FOREIGN KEY (category_id) REFERENCES category(id)
    	);
	`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("❌ Erreur lors de la création de la table post_category : %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS notification (
        	id INT AUTO_INCREMENT PRIMARY KEY,
        	user_id INT NOT NULL,
        	type VARCHAR(50) NOT NULL,
        	source_id INT NOT NULL,
        	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    	);
	`)
		if err != nil {
		db.Close()
		return nil, fmt.Errorf("❌ Erreur lors de la création de la table notification : %v", err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS comment (
		id INT AUTO_INCREMENT PRIMARY KEY,
		content TEXT NOT NULL,
		author_id INT NOT NULL,
		post_id INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (author_id) REFERENCES user(id) ON DELETE CASCADE,
		FOREIGN KEY (post_id) REFERENCES post(id) ON DELETE CASCADE
	);
`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf(" Erreur lors de la création de la table comment : %v", err)
	}
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS post_reports (
		id INT AUTO_INCREMENT PRIMARY KEY,
		reporter_id INT NOT NULL,
		post_id INT,
		comment_id INT,
		reason VARCHAR(255) NOT NULL,
		status ENUM('pending', 'approved', 'rejected') DEFAULT 'pending',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		resolved_at TIMESTAMP NULL,
		resolved_by_id INT,
		FOREIGN KEY (reporter_id) REFERENCES user(id),
		FOREIGN KEY (post_id) REFERENCES post(id) ON DELETE CASCADE,
		FOREIGN KEY (comment_id) REFERENCES comment(id) ON DELETE CASCADE,
		FOREIGN KEY (resolved_by_id) REFERENCES user(id)
	);
	`)
	if err != nil {
	db.Close()
	return nil, fmt.Errorf("❌ Erreur lors de la création de la table post_reports : %v", err)
	}

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS like_dislike (
        user_id INT NOT NULL,
        post_id INT NOT NULL,
        type INT NOT NULL, -- 0 pour dislike, 1 pour like
        PRIMARY KEY (user_id, post_id),
        FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
        FOREIGN KEY (post_id) REFERENCES post(id) ON DELETE CASCADE
    );
`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf(" Erreur lors de la création de la table like_dislike : %v", err)
	}

	fmt.Println("✅ Connexion à MySQL réussie et tables créées !")
	return db, nil
}

