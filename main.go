package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"unicode"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func strongPasswordSteps(initPassword string) int {

	lengthCriteriaMet := false
	lowercaseCriteriaMet := false
	uppercaseCriteriaMet := false
	digitCriteriaMet := false
	repeatingCriteriaMet := true

	if len(initPassword) >= 6 && len(initPassword) < 20 {
		lengthCriteriaMet = true
	}

	// Check lowercase, uppercase, and digit criteria
	for _, char := range initPassword {
		switch {
		case unicode.IsLower(char):
			lowercaseCriteriaMet = true
		case unicode.IsUpper(char):
			uppercaseCriteriaMet = true
		case unicode.IsDigit(char):
			digitCriteriaMet = true
		}
	}

	// Check repeating characters criterion
	for i := 2; i < len(initPassword); i++ {
		if initPassword[i] == initPassword[i-1] && initPassword[i-1] == initPassword[i-2] {
			repeatingCriteriaMet = false
			break
		}
	}

	// Calculate the number of steps needed
	numOfSteps := 0
	if !lengthCriteriaMet {
		numOfSteps++
	}
	if !lowercaseCriteriaMet {
		numOfSteps++
	}
	if !uppercaseCriteriaMet {
		numOfSteps++
	}
	if !digitCriteriaMet {
		numOfSteps++
	}
	if !repeatingCriteriaMet {
		numOfSteps++
	}

	return numOfSteps
}

func CheckPasswordHandler(w http.ResponseWriter, r *http.Request) {
	// Extract password from request body
	decoder := json.NewDecoder(r.Body)
	var request struct {
		Password string `json:"password"`
	}

	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call your strongPasswordSteps function or implement your password validation logic
	steps := strongPasswordSteps(request.Password)

	// Respond with the number of steps needed
	response := struct {
		Steps int `json:"steps"`
	}{Steps: steps}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	// Store the data in the PostgreSQL database
	storeDataInDatabase(request.Password, steps)
}

// Function to store data in PostgreSQL database
func storeDataInDatabase(password string, steps int) {
	// Get PostgreSQL connection string from environment variable
	connectionString := os.Getenv("DB_CONNECTION_STRING")

	// Open a connection to the database
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Insert data into the database
	_, err = db.Exec("INSERT INTO password_checks (password, steps) VALUES ($1, $2)", password, steps)
	if err != nil {
		log.Println("Error inserting data into the database:", err)
	}
}

func main() {
	// Create a new router instance from Gorilla Mux
	router := mux.NewRouter()

	// API endpoint for checking password
	router.HandleFunc("/check-password", CheckPasswordHandler).Methods("POST")

	// Specify the port to listen on
	port := ":8080"

	// Start the HTTP server
	fmt.Printf("Server listening on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, router))
	router.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", http.FileServer(http.Dir("docs"))))
}
