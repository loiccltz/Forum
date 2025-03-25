/*package main

import (
	"database/sql"
	"fmt"
	"log"
	"forum/backend"

	_ "github.com/go-sql-driver/mysql"
)

// Fonction pour tester la création et la récupération des catégories
func testCategory(db *sql.DB) {
	fmt.Println("🔹 Test : Ajout d'une catégorie...")

	// Ajouter une catégorie test
	categoryName := "Technologie"
	err := backend.CreateCategory(db, categoryName)
	if err != nil {
		log.Fatalf("❌ Erreur lors de l'ajout de la catégorie : %v", err)
	} else {
		fmt.Println("✅ Catégorie ajoutée avec succès.")
	}

	// Vérifier si la catégorie a bien été insérée
	fmt.Println("🔹 Test : Récupération des catégories...")
	categories, err := backend.GetCategories(db)
	if err != nil {
		log.Fatalf("❌ Erreur lors de la récupération des catégories : %v", err)
	}

	fmt.Println("📜 Liste des catégories :")
	for _, c := range categories {
		fmt.Printf("📂 ID: %d | Nom: %s\n", c.ID, c.Name)
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
	testCategory(db)
}
