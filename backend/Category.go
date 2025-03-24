package backend

import (
	"database/sql"
)

type Category struct {
    ID   int
    Name string
}

func GetCategories(db *sql.DB) ([]Category, error) {
    rows, err := db.Query("SELECT id, name FROM category")
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

func CreateCategory(db *sql.DB, name string) error {
    _, err := db.Exec("INSERT INTO category (name) VALUES (?)", name)
    return err
}