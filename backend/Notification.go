package backend

import (
	"database/sql"
	"time"
    "fmt"
)

type Notification struct {
    ID        int
    UserID    int
    Type      string
    SourceID  int
    CreatedAt time.Time
}

func CreateNotification(db *sql.DB, userID int, notifType string, sourceID int) error {
	_, err := db.Exec(`
		INSERT INTO notification (user_id, type, source_id) 
		VALUES (?, ?, ?)`,
		userID, notifType, sourceID,
	)
	return err
}

func GetUserNotifications(db *sql.DB, userID int) ([]Notification, error) {
    rows, err := db.Query("SELECT id, type, source_id, created_at FROM notification WHERE user_id = ?", userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var notifications []Notification
    for rows.Next() {
        var n Notification
        var createdAt []byte // Champ `created_at` récupéré sous forme de []byte
        
        // Scanner les valeurs de la ligne dans les variables correspondantes
        err := rows.Scan(&n.ID, &n.Type, &n.SourceID, &createdAt)
        if err != nil {
            return nil, err
        }

        // Convertir `createdAt` de []byte en time.Time
        n.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
        if err != nil {
            return nil, fmt.Errorf("erreur lors de la conversion de la date : %v", err)
        }

        // Ajouter la notification à la liste
        notifications = append(notifications, n)
    }
    return notifications, nil
}