package backend

import (
	"database/sql"
	"fmt"
	"log"
)

type Post struct {
	ID       int
	Title    string
	Content  string
	ImageURL string
	AuthorID int
}

func CreatePost(db *sql.DB, title, content, imageURL string, authorID int, categories []int) error {
    result, err := db.Exec("INSERT INTO post (title, content, image_url, author_id) VALUES (?, ?, ?, ?)", 
        title, content, imageURL, authorID)
    if err != nil {
        return err
    }
    
    postID, _ := result.LastInsertId()

    // Créer une notification pour l'utilisateur auteur du post
    err = CreateNotification(db, authorID, "Post créé", int(postID))
    if err != nil {
        return fmt.Errorf("Erreur lors de la création de la notification: %v", err)
    }

    for _, catID := range categories {
        _, err := db.Exec("INSERT INTO post_category VALUES (?, ?)", postID, catID)
        if err != nil {
            log.Printf("Erreur d'insertion catégorie : %v", err)
        }
    }
    return nil
}

func AddComment(db *sql.DB, content string, authorID, postID int) error {
    _, err := db.Exec("INSERT INTO comment (content, author_id, post_id) VALUES (?, ?, ?)", content, authorID, postID)
    if err != nil {
        fmt.Println("Erreur lors de l'ajout du commentaire :", err)
        return err
    }
    fmt.Println("✅ Commentaire ajouté avec succès.")

    // Récupérer l'auteur du post pour lui envoyer une notification
    var postAuthorID int
    err = db.QueryRow("SELECT author_id FROM post WHERE id = ?", postID).Scan(&postAuthorID)
    if err != nil {
        return fmt.Errorf("Erreur lors de la récupération de l'auteur du post: %v", err)
    }

    // Créer une notification pour l'auteur du post
    err = CreateNotification(db, postAuthorID, "Nouveau commentaire sur votre post", postID)
    if err != nil {
        return fmt.Errorf("Erreur lors de la création de la notification pour le commentaire: %v", err)
    }

    return nil
}

func LikePost(db *sql.DB, userID, postID, likeType int) error {
	_, err := db.Exec(`
		INSERT INTO like_dislike (user_id, post_id, type)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE type = VALUES(type);
	`, userID, postID, likeType)
	if err != nil {
		fmt.Println("Erreur lors du like/dislike :", err)
	}
	return err

    // Récupérer l'auteur du post
    var postAuthorID int
    err = db.QueryRow("SELECT author_id FROM post WHERE id = ?", postID).Scan(&postAuthorID)
    if err != nil {
        return fmt.Errorf("Erreur lors de la récupération de l'auteur du post: %v", err)
    }

    // Créer une notification pour l'auteur du post
    err = CreateNotification(db, postAuthorID, "Votre post a été aimé", postID)
    if err != nil {
        return fmt.Errorf("Erreur lors de la création de la notification pour le like: %v", err)
    }

    return nil
}

func GetPosts(db *sql.DB) ([]Post, error) {
	rows, err := db.Query("SELECT id, title, content, image_url, author_id FROM post")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.ImageURL, &p.AuthorID)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}