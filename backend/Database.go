package backend

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var DB *sql.DB


// InitDB initialise la connexion à la base de données MySQL
func InitDB() (*sql.DB, error) {

	dsn := "root:NOUVEAUMDP@tcp(127.0.0.1:3306)/forum?parseTime=true&loc=Local"

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

	err = AddDefaultCategories(db)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("❌ Erreur lors de l'ajout des catégories par défaut : %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS notifications (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT NOT NULL,
			trigger_user_id INT NOT NULL,
			type VARCHAR(50) NOT NULL,
			message TEXT NOT NULL,
			source_id INT NOT NULL,
			post_id INT NULL,
			is_read BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
			FOREIGN KEY (trigger_user_id) REFERENCES user(id) ON DELETE CASCADE
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

func GetNotificationsByUserID(db *sql.DB, userID int, limit int) ([]Notification, error) {
	var notifications []Notification
	query := `SELECT n.id, n.user_id, n.trigger_user_id, u.username, 
			  n.type, n.message, n.source_id, n.is_read, n.created_at, IFNULL(n.post_id, 0) 
			  FROM notifications n
			  LEFT JOIN user u ON n.trigger_user_id = u.id
			  WHERE n.user_id = ? 
			  ORDER BY n.created_at DESC LIMIT ?`
			  
	rows, err := db.Query(query, userID, limit)
	if err != nil {
		log.Printf("Error querying notifications: %v", err)
		return nil, err
	}
	defer rows.Close()
	
	for rows.Next() {
		var n Notification
		var nullablePostID sql.NullInt64
		
		// Changez ceci pour scanner directement dans un champ de type time.Time
		if err := rows.Scan(&n.ID, &n.UserID, &n.TriggerUserID, &n.TriggerUsername, 
						   &n.Type, &n.Message, &n.SourceID, &n.IsRead, &n.CreatedAt, &nullablePostID); err != nil {
			log.Printf("Error scanning notification row: %v", err)
			continue
		}
		
		// Vérifiez si nullablePostID est non nul avant de l'assigner
		if nullablePostID.Valid {
			n.PostID = int(nullablePostID.Int64)
		} else {
			n.PostID = 0
		}
		
		notifications = append(notifications, n)
	}
	
	return notifications, nil
}

func CountUnreadNotifications(db *sql.DB, userID int) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM notifications WHERE user_id = ? AND is_read = false"
	err := db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		log.Printf("Error counting unread notifications: %v", err)
		return 0, err
	}
	return count, nil
}

// MarkNotificationAsRead marque une notification comme lue
func MarkNotificationAsRead(db *sql.DB, notificationID int, userID int) error {
	query := "UPDATE notifications SET is_read = true WHERE id = ? AND user_id = ?"
	_, err := db.Exec(query, notificationID, userID)
	if err != nil {
		log.Printf("Error marking notification as read: %v", err)
		return err
	}
	return nil
}