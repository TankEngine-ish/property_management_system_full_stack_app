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

type User struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Statement int    `json:"statement"`
}

// main function
func main() {
	// Try to load from .env file, but don't fail if it doesn't exist
	_ = godotenv.Load()

	// Get DATABASE_URL from environment or use default
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Println("Warning: DATABASE_URL not set, using default connection string")
		dbURL = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	}

	log.Println("Connecting to database with connection string:", dbURL)

	//connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		log.Fatal(err)
	}
	defer db.Close()

	// Verify database connection
	err = db.Ping()
	if err != nil {
		log.Printf("Error pinging database: %v", err)
		log.Fatal(err)
	}
	log.Println("Successfully connected to database")

	// create table if not exists
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT, statement INT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		log.Fatal(err)
	}
	log.Println("Database table checked/created")

	// create router
	router := mux.NewRouter()

	// Add health endpoint for Kubernetes
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// Check if we can connect to the database as part of health check
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

	// wrap the router with CORS and JSON content type middlewares
	enhancedRouter := enableCORS(jsonContentTypeMiddleware(router))

	// start server
	log.Println("Starting server on :8000")
	log.Fatal(http.ListenAndServe(":8000", enhancedRouter))
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow any origin
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Check if the request is for CORS preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass down the request to the next middleware (or final handler)
		next.ServeHTTP(w, r)
	})
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set JSON Content-Type
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// get all users
func getUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM users")
		if err != nil {
			log.Printf("Error querying users: %v", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		users := []User{} // array of users
		for rows.Next() {
			var u User
			if err := rows.Scan(&u.Id, &u.Name, &u.Statement); err != nil {
				log.Printf("Error scanning user row: %v", err)
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
			users = append(users, u)
		}
		if err := rows.Err(); err != nil {
			log.Printf("Error iterating user rows: %v", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(users)
	}
}

// get user by id
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

// create user
func createUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u User
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			log.Printf("Error decoding request body: %v", err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		err = db.QueryRow("INSERT INTO users (name, statement) VALUES ($1, $2) RETURNING id", u.Name, u.Statement).Scan(&u.Id)
		if err != nil {
			log.Printf("Error creating user: %v", err)
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(u)
	}
}

// update user
func updateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u User
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			log.Printf("Error decoding request body: %v", err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		vars := mux.Vars(r)
		id := vars["id"]

		// Check if the user exists
		var exists bool
		err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)", id).Scan(&exists)
		if err != nil || !exists {
			log.Printf("User %s not found: %v", id, err)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		// Execute the update query
		_, err = db.Exec("UPDATE users SET name = $1, statement = $2 WHERE id = $3", u.Name, u.Statement, id)
		if err != nil {
			log.Printf("Error updating user %s: %v", id, err)
			http.Error(w, "Failed to update user", http.StatusInternalServerError)
			return
		}

		// Retrieve the updated user data from the database
		var updatedUser User
		err = db.QueryRow("SELECT id, name, statement FROM users WHERE id = $1", id).Scan(&updatedUser.Id, &updatedUser.Name, &updatedUser.Statement)
		if err != nil {
			log.Printf("Error retrieving updated user %s: %v", id, err)
			http.Error(w, "Failed to retrieve updated user", http.StatusInternalServerError)
			return
		}

		// Send the updated user data in the response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedUser)
	}
}

// delete user
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
				http.Error(w, "Failed to delete user", http.StatusInternalServerError)
				return
			}

			json.NewEncoder(w).Encode("User deleted")
		}
	}
}

// package main

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"log"
// 	"net/http"
// 	"os"

// 	"github.com/joho/godotenv"

// 	"github.com/gorilla/mux"
// 	_ "github.com/lib/pq"
// )

// type User struct {
// 	Id        int    `json:"id"`
// 	Name      string `json:"name"`
// 	Statement int    `json:"statement"`
// }

// // main function
// func main() {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}
// 	//connect to database
// 	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()

// 	// create table if not exists
// 	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT, statement INT)")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// create router
// 	router := mux.NewRouter()
// 	router.HandleFunc("/api/repair/users", getUsers(db)).Methods("GET")
// 	router.HandleFunc("/api/repair/users", createUser(db)).Methods("POST")
// 	router.HandleFunc("/api/repair/users/{id}", getUser(db)).Methods("GET")
// 	router.HandleFunc("/api/repair/users/{id}", updateUser(db)).Methods("PUT")
// 	router.HandleFunc("/api/repair/users/{id}", deleteUser(db)).Methods("DELETE")

// 	// wrap the router with CORS and JSON content type middlewares
// 	enhancedRouter := enableCORS(jsonContentTypeMiddleware(router))

// 	// start server
// 	log.Fatal(http.ListenAndServe(":8000", enhancedRouter))
// }

// func enableCORS(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Set CORS headers
// 		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow any origin
// 		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

// 		// Check if the request is for CORS preflight
// 		if r.Method == "OPTIONS" {
// 			w.WriteHeader(http.StatusOK)
// 			return
// 		}

// 		// Pass down the request to the next middleware (or final handler)
// 		next.ServeHTTP(w, r)
// 	})

// }

// func jsonContentTypeMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Set JSON Content-Type
// 		w.Header().Set("Content-Type", "application/json")
// 		next.ServeHTTP(w, r)
// 	})
// }

// // get all users
// func getUsers(db *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		rows, err := db.Query("SELECT * FROM users")
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		defer rows.Close()

// 		users := []User{} // array of users
// 		for rows.Next() {
// 			var u User
// 			if err := rows.Scan(&u.Id, &u.Name, &u.Statement); err != nil {
// 				log.Fatal(err)
// 			}
// 			users = append(users, u)
// 		}
// 		if err := rows.Err(); err != nil {
// 			log.Fatal(err)
// 		}

// 		json.NewEncoder(w).Encode(users)
// 	}
// }

// // get user by id
// func getUser(db *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		vars := mux.Vars(r)
// 		id := vars["id"]

// 		var u User
// 		err := db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&u.Id, &u.Name, &u.Statement)
// 		if err != nil {
// 			w.WriteHeader(http.StatusNotFound)
// 			return
// 		}

// 		json.NewEncoder(w).Encode(u)
// 	}
// }

// // create user
// func createUser(db *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var u User
// 		json.NewDecoder(r.Body).Decode(&u)

// 		err := db.QueryRow("INSERT INTO users (name, statement) VALUES ($1, $2) RETURNING id", u.Name, u.Statement).Scan(&u.Id)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		json.NewEncoder(w).Encode(u)
// 	}
// }

// // update user
// func updateUser(db *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var u User
// 		err := json.NewDecoder(r.Body).Decode(&u)
// 		if err != nil {
// 			http.Error(w, "Invalid request payload", http.StatusBadRequest)
// 			return
// 		}

// 		vars := mux.Vars(r)
// 		id := vars["id"]

// 		// Check if the user exists
// 		var exists bool
// 		err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)", id).Scan(&exists)
// 		if err != nil || !exists {
// 			http.Error(w, "User not found", http.StatusNotFound)
// 			return
// 		}

// 		// Execute the update query
// 		_, err = db.Exec("UPDATE users SET name = $1, statement = $2 WHERE id = $3", u.Name, u.Statement, id)
// 		if err != nil {
// 			http.Error(w, "Failed to update user", http.StatusInternalServerError)
// 			return
// 		}

// 		// Retrieve the updated user data from the database
// 		var updatedUser User
// 		err = db.QueryRow("SELECT id, name, statement FROM users WHERE id = $1", id).Scan(&updatedUser.Id, &updatedUser.Name, &updatedUser.Statement)
// 		if err != nil {
// 			http.Error(w, "Failed to retrieve updated user", http.StatusInternalServerError)
// 			return
// 		}

// 		// Send the updated user data in the response
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(updatedUser)
// 	}
// }

// // delete user
// func deleteUser(db *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		vars := mux.Vars(r)
// 		id := vars["id"]

// 		var u User
// 		err := db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&u.Id, &u.Name, &u.Statement)
// 		if err != nil {
// 			w.WriteHeader(http.StatusNotFound)
// 			return
// 		} else {
// 			_, err := db.Exec("DELETE FROM users WHERE id = $1", id)
// 			if err != nil {
// 				//todo : fix error handling
// 				w.WriteHeader(http.StatusNotFound)
// 				return
// 			}

// 			json.NewEncoder(w).Encode("User deleted")
// 		}
// 	}
// }
