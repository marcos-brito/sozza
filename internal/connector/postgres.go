package connector

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Postgres struct {
}

func (pg *Postgres) Connect(url string) *sql.DB {
	db, err := sql.Open("postgres", url)

	if err != nil {
		log.Fatalf("Could not connect to postgres database: %s", err)
	}

	err = db.Ping()

	if err != nil {
		log.Fatalf("Could not connect to postgres database: %s", err)
	}

	return db
}
