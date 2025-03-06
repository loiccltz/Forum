package backend

import (
	"database/sql"
	"fmt"
)

func CreatePost(db *sql.DB, title, content string, authorID int) error {
	_, err := db.Exec("INSERT INTO post (title, content, author_id) VALUES (?, ?, ?)", title, content, authorID)
	if err != nil {
		fmt.Println("Erreur lors de la création du post :", err)
		return err
	}
	fmt.Println("✅ Post créé avec succès :", title)
	return nil
}

func AddComment(db *sql.DB, content string, authorID, postID int) error {
	_, err := db.Exec("INSERT INTO comment (content, author_id, post_id) VALUES (?, ?, ?)", content, authorID, postID)
	if err != nil {
		fmt.Println("Erreur lors de l'ajout du commentaire :", err)
		return err
	}
	fmt.Println("✅ Commentaire ajouté avec succès.")
	return nil
}

func LikePost(db *sql.DB, userID, postID, likeType int) error {
	_, err := db.Exec(`
		INSERT INTO like_dislike (user_id, post_id, type)
		VALUES (?, ?, ?)
		ON CONFLICT(user_id, post_id) 
		DO UPDATE SET type = excluded.type;
	`, userID, postID, likeType)
	if err != nil {
		fmt.Println("Erreur lors du like/dislike :", err)
	}
	return err
}