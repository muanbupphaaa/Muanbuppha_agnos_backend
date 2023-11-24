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

	for i := 2; i < len(initPassword); i++ {
		if initPassword[i] == initPassword[i-1] && initPassword[i-1] == initPassword[i-2] {
			repeatingCriteriaMet = false
			break
		}
	}

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
	decoder := json.NewDecoder(r.Body)
	var request struct {
		Password string `json:"password"`
	}

	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	steps := strongPasswordSteps(request.Password)

	response := struct {
		Steps int `json:"steps"`
	}{Steps: steps}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	storeDataInDatabase(request.Password, steps)
}

func storeDataInDatabase(password string, steps int) {

	connectionString := os.Getenv("DB_CONNECTION_STRING")

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO password_checks (password, steps) VALUES ($1, $2)", password, steps)
	if err != nil {
		log.Println("Error inserting data into the database:", err)
	}
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/check-password", CheckPasswordHandler).Methods("POST")

	port := ":8080"

	fmt.Printf("Server listening on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
