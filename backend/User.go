package backend

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)


type User struct {
	ID           int
	Username     string
	Email        string
	Password     string // hashé
	SessionToken string
	Role         string
}

//insere un nouvel utilisateur dans la base de donnees
func InsertUser(db *sql.DB, username, email, password string) error {
	statement, err := db.Prepare("INSERT INTO user (username, email, password) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer statement.Close()
	
	_, err = statement.Exec(username, email, password)
	return err
}

//met a jour le token de session d'un utilisateur
func (u *User) UpdateSessionToken(db *sql.DB, token string) error {
	_, err := db.Exec("UPDATE user SET session_token = ? WHERE id = ?", token, u.ID)
	if err != nil {
		return err
	}
	u.SessionToken = token
	return nil
}


func GetUserInfoByToken(db *sql.DB, token string) (*User, error) {
	var user User
	
	// Debug: imprimer le token pour vérification
	fmt.Printf("Recherche de l'utilisateur avec le token: %s\n", token)
	
	// Vérifions que le token n'est pas vide
	if token == "" {
		return nil, errors.New("token de session vide")
	}
	
	err := db.QueryRow("SELECT id, email, username, role FROM user WHERE session_token = ?", token).Scan(&user.ID, &user.Email, &user.Username, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("Aucun utilisateur trouvé avec ce token: %s\n", token)
			return nil, errors.New("utilisateur non trouvé")
		}
		fmt.Printf("Erreur SQL: %v\n", err)
		return nil, err
	}
	
	fmt.Printf("Utilisateur trouvé: ID=%d, Username=%s\n", user.ID, user.Username)
	return &user, nil
}

func GetAllUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT id, username, email, role FROM user ORDER BY id ASC")
	if err != nil {
		log.Printf("Erreur BDD lors de la récupération de tous les utilisateurs: %v", err)
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Role); err != nil {
			log.Printf("Erreur lors du scan d'un utilisateur: %v", err)
			continue 
		}
		users = append(users, u)
	}
	if err = rows.Err(); err != nil {
		log.Printf("Erreur après l'itération sur les utilisateurs: %v", err)
		return nil, err
	}

	return users, nil
}


func UpdateUserRole(db *sql.DB, userID int, newRole string) error {
	// Valider le rôle pour s'assurer qu'il correspond aux valeurs ENUM
	isValidRole := false
	for _, validRole := range []string{RoleUser, RoleModerator, RoleAdmin} {
		if newRole == validRole {
			isValidRole = true
			break
		}
	}
	if !isValidRole {
        log.Printf("Tentative de mise à jour vers un rôle invalide '%s' pour l'utilisateur ID %d", newRole, userID)
		return fmt.Errorf("rôle invalide : %s", newRole)
	}

	result, err := db.Exec("UPDATE user SET role = ? WHERE id = ?", newRole, userID)
	if err != nil {
		log.Printf("Erreur BDD lors de la mise à jour du rôle pour l'utilisateur ID %d: %v", userID, err)
		return fmt.Errorf("erreur lors de la mise à jour du rôle")
	}

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        log.Printf("Erreur lors de la vérification des lignes affectées pour la mise à jour du rôle (User ID %d): %v", userID, err)
        // Continuer même si on ne peut pas vérifier, la mise à jour a peut-être réussi
    } else if rowsAffected == 0 {
        log.Printf("Aucun utilisateur trouvé avec l'ID %d pour la mise à jour du rôle.", userID)
        return fmt.Errorf("utilisateur non trouvé (ID: %d)", userID)
    }


	log.Printf("Rôle mis à jour avec succès pour l'utilisateur ID %d vers '%s'", userID, newRole)
	return nil
}