// Package main contains the phone number normalizer application.
package main

import (
	"fmt"
	"regexp"

	phonedb "github.com/ibrahaknazarov/phone/db"

	_ "github.com/lib/pq"
)

// Database connection configuration
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "phone"
)

// main orchestrates the phone number normalization process:
// 1. Resets the database
// 2. Runs migrations to create tables
// 3. Seeds test data
// 4. Normalizes all phone numbers and consolidates duplicates
func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	must(phonedb.Reset("postgres", psqlInfo, dbname))

	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	must(phonedb.Migrate("postgres", psqlInfo))

	db, err := phonedb.Open("postgres", psqlInfo)
	must(err)
	defer db.Close()

	err = db.Seed()
	must(err)

	phones, err := db.AllPhones()
	must(err)
	for _, p := range phones {
		fmt.Printf("Working on... %+v\n", p)
		number := normalize(p.Number)
		if number != p.Number {
			fmt.Println("Updating or removing...", number)
			existing, err := db.FindPhone(number)
			must(err)
			if existing != nil {
				must(db.DeletePhone(p.ID))
			} else {
				p.Number = number
				must(db.UpdatePhone(&p))
			}
		} else {
			fmt.Println("No changes required")
		}
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
