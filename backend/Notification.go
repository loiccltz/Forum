package backend

import (
	"database/sql"
	"time"
    "fmt"
    "log"
)

type Notification struct {
    ID              int
    UserID          int
    TriggerUserID   int
    TriggerUsername string
    Type            string
    Message         string
    SourceID        int
    PostID          int  // Make sure this exists for your PostID field
    IsRead          bool
    CreatedAt       time.Time
}

func CreateNotification(userID, triggerUserID int, notificationType string, sourceID, postID int, message string) error {
    query := `INSERT INTO notifications (user_id, trigger_user_id, type, source_id, message, post_id, is_read, created_at) 
              VALUES (?, ?, ?, ?, ?, ?, false, NOW())`
    _, err := DB.Exec(query, userID, triggerUserID, notificationType, sourceID, message, postID)
    if err != nil {
        log.Printf("Error creating notification: %v", err)
        return err
    }
    return nil
}

func GetUserNotifications(db *sql.DB, userID int) ([]Notification, error) {
    rows, err := db.Query("SELECT id, type, source_id, created_at FROM notifications WHERE user_id = ?", userID)
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

func GetNotifications(userID int, limit int) ([]Notification, error) {
    return GetNotificationsByUserID(DB, userID, limit) 
}

// CountUnread counts unread notifications for a user
func CountUnread(userID int) int {
    count, err := CountUnreadNotifications(DB, userID) 
    if err != nil {
        return 0
    }
    return count
}