# h3-travel

Mini projet de gestion de voyages avec Go, Gin, GORM, JWT, Swagger et PGSQL.

Plus d'infos dans le fichier backend/consignes.txt

---

## 1️⃣ Prérequis

- [Go](https://golang.org/dl/) >= 1.21  
- [Docker](https://www.docker.com/get-started/) & Docker Compose  
- Git  

---

## 2️⃣ Cloner le projet

```bash
git clone https://github.com/Vincent-Murienne/h3-travel.git
cd h3-travel/backend
```

---

## 3️⃣ Initialiser le projet Go
```bash
go mod init h3-travel ## Création du module Go.
go mod tidy ## Récupère toutes les dépendances importées dans le code.
```

---

## 4️⃣ Installer les packages nécessaires
```bash
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/postgres  # ou sqlite selon la config
go get github.com/golang-jwt/jwt/v5
go get github.com/joho/godotenv
go get github.com/stretchr/testify
go get modernc.org/sqlite        # pour tests unitaires sans CGO
```

---

## 5️⃣ Configurer les variables d’environnement
```bash
## Créez un fichier .env.development à la racine du dossier backend
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=h3travel
JWT_SECRET=une_cle_secrete_longue_et_aleatoire_genere_manuellement
```

---

## 6️⃣ Utilisation avec Docker

Normalement pas besoin d'y toucher

---

## 7️⃣ Lancer le projet :
```bash
docker-compose up --build -d

## Création d'un compte admin
## Remplacer : default:'user' par default:'admin'

type User struct {
	gorm.Model
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
	Role     string `gorm:"type:varchar(10);default:'user'"` // "user" ou "admin"
}

## Modification / Ajout de requêtes accessibles sur Swagger
swag init

## Rentrer dans le containeur db h3_travel :
docker exec -it h3_travel_db psql -U admin -d h3_travel

## Rentrer dans le containeur postgres :
docker exec -it h3_travel_db psql -U admin -d postgres

\dt              ## liste toutes les tables
\d users         ## structure de la table users (par exemple)
SELECT * FROM users;  ## voir les utilisateurs

## Serveur Gin accessible sur http://localhost:8080
## Swagger accessible sur http://localhost:8080/swagger/index.html
## SonarQube accessible sur (intégration avec Github Actions) http://localhost:9000
## Grafana accessible sur http://localhost:3000
## Prometheus accessible sur http://localhost:9090
## Node Exporter accessible sur http://localhost:9100
## Base PostgreSQL sur le port 5432.

## Commandes tests unitaires depuis le dossier backend 
go test ./... -coverprofile=coverage.out

## Pour faire fonctionner SonarQube : retirer '.example' dans sonar-project.properties 
```
