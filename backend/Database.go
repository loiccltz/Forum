package backend

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// InitDB initialise la connexion à la base de données MySQL
func InitDB() (*sql.DB, error) {

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

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

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS user (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		session_token VARCHAR(64) DEFAULT '' NOT NULL,
		role ENUM('user', 'moderator', 'admin') NOT NULL DEFAULT 'user',
		google_id VARCHAR(255) UNIQUE,  -- Ajout de google_id pour l'authentification via Google
		auth_type ENUM('password', 'google') DEFAULT 'password'  -- Ajout du type d'authentification
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
	err = AddDefaultCategories(db)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("❌ Erreur lors de l'ajout des catégories par défaut : %v", err)
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

// CreateDefaultAdmin vérifie si un admin existe et le crée sinon.
func CreateDefaultAdmin(db *sql.DB) error {
	defaultAdminEmail := "admin@admin.com"
	defaultAdminUsername := "admin"
	defaultPassword := "admin" // pour test 

	// 1. Vérifier si l'admin existe déjà par email
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM user WHERE email = ?)", defaultAdminEmail).Scan(&exists)
	if err != nil {
		return fmt.Errorf("erreur lors de la vérification de l'existence de l'admin: %w", err)
	}

	if exists {
		log.Println("ℹ️ L'utilisateur admin par défaut existe déjà.")
		return nil 
	}


	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM user WHERE username = ?)", defaultAdminUsername).Scan(&exists)
	if err != nil {
		return fmt.Errorf("erreur lors de la vérification de l'existence du username admin: %w", err)
	}
    if exists {
        log.Printf("⚠️ Le username '%s' existe déjà, impossible de créer l'admin par défaut avec ce username.", defaultAdminUsername)
        return fmt.Errorf("le username '%s' existe déjà", defaultAdminUsername)
    }


	// 2. Hasher le mot de passe par défaut
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("erreur lors du hashage du mot de passe admin: %w", err)
	}

	// 3. Insérer le nouvel admin
	_, err = db.Exec("INSERT INTO user (username, email, password, role) VALUES (?, ?, ?, ?)",
		defaultAdminUsername,
		defaultAdminEmail,
		string(hashedPassword),
		RoleAdmin, // Utilise la constante RoleAdmin de roles.go
	)
	if err != nil {
		return fmt.Errorf("erreur lors de l'insertion de l'admin par défaut: %w", err)
	}

	log.Printf("✅ Admin par défaut créé avec succès : email=%s, username=%s", defaultAdminEmail, defaultAdminUsername)
	log.Printf("🔑 Mot de passe admin par défaut : %s (À CHANGER IMMÉDIATEMENT !)", defaultPassword)

	return nil
}
