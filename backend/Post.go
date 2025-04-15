package backend

import (
	"database/sql"
	"errors"
	"fmt"
)

type Post struct {
	ID         int
	Title      string
	Content    string
	ImageURL   string
	AuthorID   int
	AuthorName string
	Categories []Category
}

func GetPostAuthor(db *sql.DB, postID int) (int, error) {
	var authorID int
	err := db.QueryRow("SELECT author_id FROM post WHERE id = ?", postID).Scan(&authorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("post not found")
		}
		return 0, err
	}
	return authorID, nil
}

func UpdatePost(db *sql.DB, postID int, currentUserID int, title string, content string) error {
	// 1. Verify ownership
	postAuthorID, err := GetPostAuthor(db, postID)
	if err != nil {
		return err
	}
	if postAuthorID != currentUserID {
		return errors.New("user not authorized to edit this post")
	}

	// mettre a jour un post
	_, err = db.Exec("UPDATE post SET title = ?, content = ? WHERE id = ?", title, content, postID)
	if err != nil {
		fmt.Printf("Error updating post (ID: %d): %v\n", postID, err)
		return errors.New("failed to update post")
	}

	fmt.Printf("✅ Post (ID: %d) updated successfully by user (ID: %d)\n", postID, currentUserID)
	return nil
}

// supprime un post de la bdd
func DeletePost(db *sql.DB, postID int, currentUserID int) error {
	postAuthorID, err := GetPostAuthor(db, postID)
	if err != nil {
		return err // Post not found or DB error
	}
	if postAuthorID != currentUserID {
		// permettre aux admin de delete
		// user, _ := GetUserInfoByID(db, currentUserID) // Need a function like this
		// if !IsAdmin(user) && !IsModerator(user) {
		//     return errors.New("user not authorized to delete this post")
		// }
		return errors.New("user not authorized to delete this post")
	}


	_, err = db.Exec("DELETE FROM post_category WHERE post_id = ?", postID)
	if err != nil {
		fmt.Printf("Error deleting post_category entries for post (ID: %d): %v\n", postID, err)
		return errors.New("failed to delete post-category links")
	}

	// supprimer le post
	_, err = db.Exec("DELETE FROM post WHERE id = ?", postID)
	if err != nil {
		fmt.Printf("Error deleting post (ID: %d): %v\n", postID, err)
		return errors.New("failed to delete post")
	}

	//TODO : supprimer l'image avec le post

	fmt.Printf("✅ Post (ID: %d) deleted successfully by user (ID: %d)\n", postID, currentUserID)
	return nil
}

func GetPostsByCategory(db *sql.DB, categoryID int) ([]Post, error) {
	query := `
        SELECT p.id, p.title, p.content, p.image_url, p.author_id
        FROM post p
        JOIN post_category pc ON p.id = pc.post_id
        WHERE pc.category_id = ?
    `
	rows, err := db.Query(query, categoryID)
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

		categories, err := GetPostCategories(db, p.ID)
		if err != nil {
			return nil, err
		}
		p.Categories = categories

		posts = append(posts, p)
	}

	return posts, nil
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

		// Récupérer les catégories pour ce post
		categories, err := GetPostCategories(db, p.ID)
		if err != nil {
			return nil, err
		}
		p.Categories = categories

		posts = append(posts, p)
	}

	return posts, nil
}

func GetPostByID(db *sql.DB, postID int) (*Post, error) {
	var post Post
	err := db.QueryRow(`
        SELECT p.id, p.title, p.content, p.image_url, p.author_id, u.username
        FROM post p
        JOIN user u ON p.author_id = u.id
        WHERE p.id = ?`, postID).Scan(
		&post.ID, &post.Title, &post.Content, &post.ImageURL, &post.AuthorID, &post.AuthorName, // Scan AuthorName
	)
	if err != nil {
		return nil, err
	}

	// Récupérer les catégories pour ce post
	categories, err := GetPostCategories(db, post.ID)
	if err != nil {
		return nil, err
	}
	post.Categories = categories

	return &post, nil
}

func CreatePost(db *sql.DB, title, content, imageURL string, authorID int, categoryIDs []int) (int, error) {

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	// Insérer le post
	result, err := tx.Exec("INSERT INTO post (title, content, image_url, author_id) VALUES (?, ?, ?, ?)",
		title, content, imageURL, authorID)
	if err != nil {
		tx.Rollback()
		fmt.Println("Erreur lors de la création du post :", err)
		return 0, err
	}

	// Récupérer l'ID du nouveau post
	postID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Associer les catégories au post
	if len(categoryIDs) > 0 {
		stmt, err := tx.Prepare("INSERT INTO post_category (post_id, category_id) VALUES (?, ?)")
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		defer stmt.Close()

		for _, categoryID := range categoryIDs {
			_, err := stmt.Exec(postID, categoryID)
			if err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	fmt.Println("✅ Post créé avec succès :", title)
	return int(postID), nil
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
		return fmt.Errorf("erreur lors de la récupération de l'auteur du post: %v", err)
	}

	// Créer une notification pour l'auteur du post
	err = CreateNotification(db, postAuthorID, "Nouveau commentaire sur votre post", postID)
	if err != nil {
		return fmt.Errorf("erreur lors de la création de la notification pour le commentaire: %v", err)
	}

	return nil
}

func LikePost(db *sql.DB, userID, postID, likeType int) error {
	_, err := db.Exec(`
        INSERT INTO like_dislike (user_id, post_id, type)
        VALUES (?, ?, ?)
        ON DUPLICATE KEY UPDATE type = ?;
    `, userID, postID, likeType, likeType)
	if err != nil {
		fmt.Println("Erreur lors du like/dislike :", err)
	}
	return err
}

func CountLikes(db *sql.DB, postID int) (int, int, error) {
	var likes, dislikes int

	// Compter les likes (type = 1 = like)
	err := db.QueryRow("SELECT COUNT(*) FROM like_dislike WHERE post_id = ? AND type = 1", postID).Scan(&likes)
	if err != nil {
		return 0, 0, err
	}

	// Compter les dislikes (type = 0 = dislike)
	err = db.QueryRow("SELECT COUNT(*) FROM like_dislike WHERE post_id = ? AND type = 0", postID).Scan(&dislikes)
	if err != nil {
		return 0, 0, err
	}

	return likes, dislikes, nil
}

// associe un post a une ou plusieurs categories
func AssignCategoriesToPost(db *sql.DB, postID int, categoryIDs []int) error {

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Préparer la requête
	stmt, err := tx.Prepare("INSERT INTO post_category (post_id, category_id) VALUES (?, ?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	// Insérer chaque catégorie
	for _, categoryID := range categoryIDs {
		_, err := stmt.Exec(postID, categoryID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// reécupère les catégories associés a un post
func GetPostCategories(db *sql.DB, postID int) ([]Category, error) {
	rows, err := db.Query(`
        SELECT c.id, c.name 
        FROM category c
        JOIN post_category pc ON c.id = pc.category_id
        WHERE pc.post_id = ?
    `, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}
