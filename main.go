// Package main contains the phone number normalizer application.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/ibrahaknazarov/phone/config"
	phonedb "github.com/ibrahaknazarov/phone/db"

	_ "github.com/lib/pq"
)

// main orchestrates the phone number normalization process and starts a web server.
func main() {
	cfg := config.Load()

	// Initial DB setup
	must(phonedb.Reset("postgres", cfg.PSQLInfo(), cfg.DBName))
	must(phonedb.Migrate("postgres", cfg.PSQLInfoWithDB(cfg.DBName)))

	db, err := phonedb.Open("postgres", cfg.PSQLInfoWithDB(cfg.DBName))
	must(err)
	defer db.Close()

	// Seed initial data
	err = db.Seed()
	must(err)

	// Run initial normalization cleanup
	syncDatabase(db)

	// Handlers
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/api/normalize", handleNormalize)
	http.HandleFunc("/api/phones", handleGetPhones(db))
	http.HandleFunc("/api/sync", handleSync(db))

	addr := ":8080"
	fmt.Printf("Server starting on %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func syncDatabase(db *phonedb.DB) {
	phones, err := db.AllPhones()
	if err != nil {
		log.Printf("Error getting phones: %v", err)
		return
	}
	for _, p := range phones {
		number := normalize(p.Number)
		if number != p.Number {
			existing, err := db.FindPhone(number)
			if err != nil {
				log.Printf("Error finding phone %s: %v", number, err)
				continue
			}
			if existing != nil {
				if err := db.DeletePhone(p.ID); err != nil {
					log.Printf("Error deleting duplicate phone %d: %v", p.ID, err)
				}
			} else {
				p.Number = number
				if err := db.UpdatePhone(&p); err != nil {
					log.Printf("Error updating phone %d: %v", p.ID, err)
				}
			}
		}
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func handleNormalize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		Phone string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	normalized := normalize(data.Phone)
	json.NewEncoder(w).Encode(map[string]string{
		"original":   data.Phone,
		"normalized": normalized,
	})
}

func handleGetPhones(db *phonedb.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phones, err := db.AllPhones()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(phones)
	}
}

func handleSync(db *phonedb.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		syncDatabase(db)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Database synced and normalized")
	}
}

// must panics if err is not nil. Used for error handling in main function.
func must(err error) {
	if err != nil {
		panic(err)
	}
}

// normalize removes all non-digit characters from a phone number string.
// Example: "(123) 456-7890" becomes "1234567890"
func normalize(phone string) string {
	re := regexp.MustCompile("\\D")
	return re.ReplaceAllString(phone, "")
}

// func normalize(phone string) string {
// 	var buf bytes.Buffer
// 	for _, ch := range phone {
// 		if ch >= '0' && ch <= '9' {
// 			buf.WriteRune(ch)
// 		}
// 	}
// 	return buf.String()
// }
