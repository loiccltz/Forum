/*package main

import (
	"database/sql"
	"fmt"
	"log"
	"forum/backend"

	_ "github.com/go-sql-driver/mysql"
)

// Fonction pour tester la création et la récupération des notifications
func testNotification(db *sql.DB) {
	// Utilisateur avec ID 4 pour tester
	userID := 4
	notifType := "Nouveau commentaire"
	sourceID := 123 // L'ID du post ou autre source d'une notification

	// Créer une notification pour cet utilisateur
	err := backend.CreateNotification(db, userID, notifType, sourceID)
	if err != nil {
		log.Fatalf("❌ Erreur lors de la création de la notification : %v", err)
	} else {
		fmt.Println("✅ Notification ajoutée avec succès.")
	}

	// Vérifier si la notification a bien été ajoutée
	notifications, err := backend.GetUserNotifications(db, userID)
	if err != nil {
		log.Fatalf("❌ Erreur lors de la récupération des notifications : %v", err)
	}

	// Afficher les notifications
	fmt.Println("📜 Liste des notifications :")
	for _, notif := range notifications {
		fmt.Printf("🔔 ID: %d | Type: %s | Source ID: %d | Date: %v\n", notif.ID, notif.Type, notif.SourceID, notif.CreatedAt)
	}
}

func main() {
	// Initialisation de la base de données
	db, err := backend.InitDB()
	if err != nil {
		log.Fatal("❌ Erreur de connexion à la base de données :", err)
	}
	defer db.Close()

	// Exécuter le test des notifications
	testNotification(db)
}
