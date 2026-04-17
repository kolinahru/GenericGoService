package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func initDB() *sql.DB {
	connStr := "host=localhost port=5432 user=postgres password=T@p030379 dbname=goapp sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(time.Hour)

	log.Println("Connected to DB")
	return db
}
