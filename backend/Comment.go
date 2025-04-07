package backend

import (
	"database/sql"
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