package repository

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"log"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const (
	bookTable   = "book"
	authorTable = "author"
)

var db *sql.DB

func InitDB(dataSourceName string) (*sql.DB, error) {
	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	log.Println("Successfully connected!")

	return db, nil
}

func MigrateDB() error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance("file://schema", "postgres", driver)
	if err != nil {
		return err
	}
	err = m.Up()
	mResult := "successfully"
	if err != nil {
		mResult = err.Error()
	}
	log.Printf("Migrate result: %s\n", mResult)

	return nil
}

func CloseDB() error {
	return db.Close()
}
