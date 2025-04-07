package backend

import (
	"database/sql"
	"fmt"
)

type Post struct {
	ID       int
	Title    string
	Content  string
	ImageURL string
	AuthorID int
}

func GetPostByID(db *sql.DB, postID int) (*Post, error) {
    var post Post
    err := db.QueryRow("SELECT id, title, content, image_url, author_id FROM post WHERE id = ?", postID).Scan(
        &post.ID, &post.Title, &post.Content, &post.ImageURL, &post.AuthorID)
    if err != nil {
        return nil, err
    }
    return &post, nil
}

func CreatePost(db *sql.DB, title, content, imageURL string, authorID int) error {
	_, err := db.Exec("INSERT INTO post (title, content, image_url, author_id) VALUES (?, ?, ?, ?)", title, content, imageURL, authorID)
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