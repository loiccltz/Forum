 # Forum Project

 ## Membres du projet
 - Loic
 - Theo
 - Axel
 - Lois

 ## Description du projet
 Le projet **Forum** est un forum en **Go (Golang)**. Les utilisateurs peuvent s'inscrire, se connecter, créer des posts, commenter et interagir avec les autres membres.

 ## Lancer le projet

 ### Utiliser Docker
 - **Lancer le projet avec Docker** :
   Utilise la commande suivante pour lancer le projet avec Docker :
   ```bash
   docker-compose up --build
   ```

 - **Problèmes avec Docker ?** 
   Si tu rencontres des erreurs avec Docker, tu peux lancer le projet directement en Go avec la commande suivante :
   ```bash
   go run main.go
   ```
 Il faut avoir une base de données "forum" créée. Si elle n'est pas encore créée, utilise la commande suivante :
 ```bash
  mysql -u root -p
  ```
- Renseigne ton mot de passe.
- Puis, crée la base de données en exécutant la commande suivante :
  ```bash
  CREATE DATABASE forum;
  ```

 ### Configuration des certificats SSL
 Le projet nécessite un certificat SSL pour fonctionner en local.

 - **Générer les certificats** :
   Utilise `mkcert` pour générer les certificats SSL nécessaires. La commande est la suivante :
   ```bash
   mkcert localhost 127.0.0.1 ::1
   ```

 - **Noms des fichiers** :
   Assure-toi que les certificats générés se nomment bien :
   - `localhost+2.pem`
   - `localhost+2-key.pem`

## Axe d'amélioration
- Finaliser la connexion avec Google
- Finir l'integration de la modération (fonctionnelle partiellement acutellement)
- 
