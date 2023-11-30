package models

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
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
}
