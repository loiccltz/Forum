package backend

import (
	"database/sql"
)

type Filter struct {
    CategoryID int
    UserID     int
}

func GetPostsByFilter(db *sql.DB, filter Filter) ([]Post, error) {
    query := "SELECT id, title, content FROM posts WHERE 1=1"
    var args []interface{}

    if filter.CategoryID != 0 {
        query += " AND category_id = ?"
        args = append(args, filter.CategoryID)
    }
    if filter.UserID != 0 {
        query += " AND user_id = ?"
        args = append(args, filter.UserID)
    }

    rows, err := db.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var posts []Post
    for rows.Next() {
        var p Post
        if err := rows.Scan(&p.ID, &p.Title, &p.Content); err != nil {
            return nil, err
        }
        posts = append(posts, p)
    }
    return posts, nil
}
