package models

import (
	"database/sql"
	"log"
)

func OpenConnections() *sql.DB {
	connStr := "host=host.docker.internal port=5432 user=postgres password=fkubifkom10 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalln(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}
	return db
}
