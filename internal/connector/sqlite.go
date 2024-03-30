package connector

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
}

func (pg *Sqlite) Connect(url string) *sql.DB {
	db, err := sql.Open("sqlite3", url)

	if err != nil {
		log.Fatalf("Could not connect to sqlite database: %s", err)
	}

	err = db.Ping()

	if err != nil {
		log.Fatalf("Could not connect to sqlite database: %s", err)
	}

	return db
}
