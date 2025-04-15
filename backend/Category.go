package backend

import (
	"database/sql"
)

type Category struct {
    ID   int
    Name string
}

func AddDefaultCategories(db *sql.DB) error {
    defaultCategories := []string{"Hardware", "CPU", "GPU"}
    
    for _, categoryName := range defaultCategories {
        // Vérifier si la catégorie existe déjà
        var exists bool
        err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM category WHERE name = ?)", categoryName).Scan(&exists)
        if err != nil {
            return err
        }
        
        // Créer la catégorie si elle n'existe pas
        if !exists {
            if err := CreateCategory(db, categoryName); err != nil {
                return err
            }
        }
    }
    
    return nil
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