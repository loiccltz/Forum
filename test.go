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
	fmt.Println("🔹 Test : Ajout d'une notification...")

	// Ajouter une notification test
	err := backend.CreateNotification(db, 1, "new_comment", 123)
	if err != nil {
		log.Fatalf("❌ Erreur lors de l'ajout de la notification : %v", err)
	} else {
		fmt.Println("✅ Notification ajoutée avec succès.")
	}

	// Vérifier si la notification a bien été insérée
	fmt.Println("🔹 Test : Récupération des notifications...")
	rows, err := db.Query("SELECT id, user_id, type, source_id, created_at FROM notification WHERE user_id = ?", 1)
	if err != nil {
		log.Fatalf("❌ Erreur lors de la récupération des notifications : %v", err)
	}
	defer rows.Close()

	fmt.Println("📜 Liste des notifications :")
	for rows.Next() {
		var id, userID, sourceID int
		var notifType string
		var createdAt string

		err := rows.Scan(&id, &userID, &notifType, &sourceID, &createdAt)
		if err != nil {
			log.Fatalf("❌ Erreur lors du scan des résultats : %v", err)
		}

		fmt.Printf("🔔 Notification %d | Utilisateur: %d | Type: %s | Source: %d | Date: %s\n",
			id, userID, notifType, sourceID, createdAt)
	}

	fmt.Println("✅ Test terminé avec succès.")
}

func main() {
	// Initialisation de la base de données
	db, err := backend.InitDB()
	if err != nil {
		log.Fatal("❌ Erreur de connexion à la base de données :", err)
	}
	defer db.Close()

	// Exécuter le test
	testNotification(db)
}
