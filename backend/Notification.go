package backend

import (
	"database/sql"
	"time"
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
        err := rows.Scan(&n.ID, &n.Type, &n.SourceID, &n.CreatedAt)
        if err != nil {
            return nil, err
        }
        notifications = append(notifications, n)
    }
    return notifications, nil
}