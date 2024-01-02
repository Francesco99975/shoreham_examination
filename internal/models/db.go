package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func Setup(dsn string) {
	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}

	go func ()  {
		for {
			cleanDatabase(db)
			time.Sleep(24 * time.Hour)
		}
	}()

	var count int

	rows, err := db.Query("SELECT COUNT(*) FROM members;")

	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			fmt.Println(err)
		}
	}

	if count == 0 {
		email := os.Getenv("ADMIN_EMAIL")
		password := os.Getenv("ADMIN_PASSWORD")

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
		if err != nil {
			fmt.Println(err)
		}

		statement := "INSERT INTO members(email, password) VALUES($1, $2);"

		_, err = db.Exec(statement, email, hashedPassword)

		if err != nil {
			fmt.Println(err)
		}
	}

}



func cleanDatabase(db *sql.DB) {
	cutoffTime := time.Now().Add(-7 * 24 * time.Hour)

	// Execute the delete query
	_, err := db.Exec("DELETE FROM patients WHERE created < $1;", cutoffTime)
	if err != nil {
		log.Println("Error deleting old data:", err)
	}

	// Execute the delete query
	_, err = db.Exec("DELETE FROM adminresults WHERE created < $1;", cutoffTime)
	if err != nil {
		log.Println("Error deleting old data:", err)
	}

	// Execute the delete query
	_, err = db.Exec("DELETE FROM examinations WHERE created < $1;", cutoffTime)
	if err != nil {
		log.Println("Error deleting old data:", err)
	}
}
