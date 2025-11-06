package dbconfig

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func ConnectDB(DatabaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", DatabaseURL)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	fmt.Println("Connected to the database!")
	return db, nil
}