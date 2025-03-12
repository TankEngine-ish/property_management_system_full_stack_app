package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// as per the sonarqube report, I have defined the below constants to avoid the hardcoding of the values
const (
	headerContentType               = "Content-Type"
	headerAccessControlAllowOrigin  = "Access-Control-Allow-Origin"
	headerAccessControlAllowMethods = "Access-Control-Allow-Methods"
	headerAccessControlAllowHeaders = "Access-Control-Allow-Headers"

	errMsgDatabase                    = "Database error"
	errMsgInvalidRequestPayload       = "Invalid request payload"
	errMsgUserNotFound                = "User not found"
	errMsgFailedToCreateUser          = "Failed to create user"
	errMsgFailedToUpdateUser          = "Failed to update user"
	errMsgFailedToDeleteUser          = "Failed to delete user"
	errMsgFailedToRetrieveUpdatedUser = "Failed to retrieve updated user"

	contentTypeJSON = "application/json"
)

type User struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Statement int    `json:"statement"`
}

func main() {
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Println("Error: DATABASE_URL environment variable not set")
		log.Fatal("Please set the DATABASE_URL environment variable with a valid connection string")
	}

	log.Println("Connecting to database")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Printf("Error pinging database: %v", err)
		log.Fatal(err)
	}
	log.Println("Successfully connected to database")

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT, statement INT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		log.Fatal(err)
	}
	log.Println("Database table checked/created")

	router := mux.NewRouter()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		err := db.Ping()
		if err != nil {
			log.Printf("Health check failed - DB connection error: %v", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(fmt.Sprintf("DB Connection Error: %v", err)))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	router.HandleFunc("/api/repair/users", getUsers(db)).Methods("GET")
	router.HandleFunc("/api/repair/users", createUser(db)).Methods("POST")
	router.HandleFunc("/api/repair/users/{id}", getUser(db)).Methods("GET")
	router.HandleFunc("/api/repair/users/{id}", updateUser(db)).Methods("PUT")
	router.HandleFunc("/api/repair/users/{id}", deleteUser(db)).Methods("DELETE")

	enhancedRouter := enableCORS(jsonContentTypeMiddleware(router))

	log.Println("Starting server on :8000")
	log.Fatal(http.ListenAndServe(":8000", enhancedRouter))
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(headerAccessControlAllowOrigin, "*")
		w.Header().Set(headerAccessControlAllowMethods, "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set(headerAccessControlAllowHeaders, headerContentType)

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(headerContentType, contentTypeJSON)
		next.ServeHTTP(w, r)
	})
}

func getUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM users")
		if err != nil {
			log.Printf("Error querying users: %v", err)
			http.Error(w, errMsgDatabase, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		users := []User{}
		for rows.Next() {
			var u User
			if err := rows.Scan(&u.Id, &u.Name, &u.Statement); err != nil {
				log.Printf("Error scanning user row: %v", err)
				http.Error(w, errMsgDatabase, http.StatusInternalServerError)
				return
			}
			users = append(users, u)
		}
		if err := rows.Err(); err != nil {
			log.Printf("Error iterating user rows: %v", err)
			http.Error(w, errMsgDatabase, http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(users)
	}
}

func getUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var u User
		err := db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&u.Id, &u.Name, &u.Statement)
		if err != nil {
			log.Printf("Error getting user %s: %v", id, err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(u)
	}
}

func createUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u User
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			log.Printf("Error decoding request body: %v", err)
			http.Error(w, errMsgInvalidRequestPayload, http.StatusBadRequest)
			return
		}

		err = db.QueryRow("INSERT INTO users (name, statement) VALUES ($1, $2) RETURNING id", u.Name, u.Statement).Scan(&u.Id)
		if err != nil {
			log.Printf("Error creating user: %v", err)
			http.Error(w, errMsgFailedToCreateUser, http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(u)
	}
}

func updateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u User
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			log.Printf("Error decoding request body: %v", err)
			http.Error(w, errMsgInvalidRequestPayload, http.StatusBadRequest)
			return
		}

		vars := mux.Vars(r)
		id := vars["id"]

		var exists bool
		err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)", id).Scan(&exists)
		if err != nil || !exists {
			log.Printf("User %s not found: %v", id, err)
			http.Error(w, errMsgUserNotFound, http.StatusNotFound)
			return
		}

		_, err = db.Exec("UPDATE users SET name = $1, statement = $2 WHERE id = $3", u.Name, u.Statement, id)
		if err != nil {
			log.Printf("Error updating user %s: %v", id, err)
			http.Error(w, errMsgFailedToUpdateUser, http.StatusInternalServerError)
			return
		}

		var updatedUser User
		err = db.QueryRow("SELECT id, name, statement FROM users WHERE id = $1", id).Scan(&updatedUser.Id, &updatedUser.Name, &updatedUser.Statement)
		if err != nil {
			log.Printf("Error retrieving updated user %s: %v", id, err)
			http.Error(w, errMsgFailedToRetrieveUpdatedUser, http.StatusInternalServerError)
			return
		}

		w.Header().Set(headerContentType, contentTypeJSON)
		json.NewEncoder(w).Encode(updatedUser)
	}
}

func deleteUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var u User
		err := db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&u.Id, &u.Name, &u.Statement)
		if err != nil {
			log.Printf("User %s not found: %v", id, err)
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			_, err := db.Exec("DELETE FROM users WHERE id = $1", id)
			if err != nil {
				log.Printf("Error deleting user %s: %v", id, err)
				http.Error(w, errMsgFailedToDeleteUser, http.StatusInternalServerError)
				return
			}

			json.NewEncoder(w).Encode("User deleted")
		}
	}
}
