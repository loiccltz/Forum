package backend

import (
	"database/sql"
	"errors"
	"fmt"
)

type Comment struct {
    ID       int
    Content  string
    Username string
    PostID   int
    AuthorID int
}

func GetCommentsByPostID(db *sql.DB, postID int) ([]Comment, error) {
    rows, err := db.Query(`
        SELECT c.id, c.content, u.username, c.post_id, c.author_id 
        FROM comment c 
        JOIN user u ON c.author_id = u.id 
        WHERE c.post_id = ?
        ORDER BY c.created_at DESC
    `, postID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var comments []Comment
    for rows.Next() {
        var c Comment
        err := rows.Scan(&c.ID, &c.Content, &c.Username, &c.PostID, &c.AuthorID)
        if err != nil {
            return nil, err
        }
        comments = append(comments, c)
    }

    return comments, nil
}


// recup l'id de l'utilisateur sur un comment
func GetCommentAuthor(db *sql.DB, commentID int) (int, error) {
	var authorID int
	err := db.QueryRow("SELECT author_id FROM comment WHERE id = ?", commentID).Scan(&authorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("comment not found")
		}
		return 0, err
	}
	return authorID, nil
}

// mettre a jour un commentaire
func UpdateComment(db *sql.DB, commentID int, currentUserID int, content string) error {
	commentAuthorID, err := GetCommentAuthor(db, commentID)
	if err != nil {
		return err
	}
	if commentAuthorID != currentUserID {
		return errors.New("user not authorized to edit this comment")
	}

	_, err = db.Exec("UPDATE comment SET content = ? WHERE id = ?", content, commentID)
	if err != nil {
		fmt.Printf("Error updating comment (ID: %d): %v\n", commentID, err)
		return errors.New("failed to update comment")
	}

	fmt.Printf("✅ Comment (ID: %d) updated successfully by user (ID: %d)\n", commentID, currentUserID)
	return nil
}

// suppr un commentaire de la bdd
func DeleteComment(db *sql.DB, commentID int, currentUserID int) error {
	commentAuthorID, err := GetCommentAuthor(db, commentID)
	if err != nil {
		return err
	}
	if commentAuthorID != currentUserID {
		// user, _ := GetUserInfoByID(db, currentUserID) // Need a function like this
		// if !IsAdmin(user) && !IsModerator(user) {
		//     return errors.New("user not authorized to delete this comment")
		// }
		return errors.New("user not authorized to delete this comment")
	}

	_, err = db.Exec("DELETE FROM comment WHERE id = ?", commentID)
	if err != nil {
		fmt.Printf("Error deleting comment (ID: %d): %v\n", commentID, err)
		return errors.New("failed to delete comment")
	}

	fmt.Printf("✅ Comment (ID: %d) deleted successfully by user (ID: %d)\n", commentID, currentUserID)
	return nil
}