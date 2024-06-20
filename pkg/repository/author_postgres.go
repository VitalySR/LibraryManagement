package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type AuthorPostgres struct{}

func NewAuthorPostgres() *AuthorPostgres {
	return &AuthorPostgres{}
}

func (b *AuthorPostgres) Create(tx *sql.Tx, author *Author) (int, error) {
	log.Println("AuthorPostgres. Create")

	var id int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (FirstName, LastName, Biography, BirthDate) values ($1, $2, $3, $4) RETURNING id", authorTable)

	row := tx.QueryRow(createItemQuery, author.FirstName, author.LastName, author.Biography, author.BirthDate)
	err := row.Scan(&id)

	return id, err
}

func (b *AuthorPostgres) GetAll(db *sql.DB) ([]Author, error) {
	log.Println("AuthorPostgres. GetAll")
	query := fmt.Sprintf("select Id, FirstName, LastName, Biography, BirthDate from %s", authorTable)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	aus := make([]Author, 0, 5)
	for rows.Next() {
		author := Author{}
		if err := rows.Scan(&author.ID, &author.FirstName, &author.LastName, &author.Biography, &author.BirthDate); err != nil {
			return nil, err
		}

		aus = append(aus, author)
	}
	return aus, nil
}

func (b *AuthorPostgres) GetById(db *sql.DB, id int) (Author, error) {
	log.Println("AuthorPostgres. GetById")
	query := fmt.Sprintf("select Id, FirstName, LastName, Biography, BirthDate from %s where Id = $1", authorTable)
	row := db.QueryRow(query, id)
	author := Author{}
	err := row.Scan(&author.ID, &author.FirstName, &author.LastName, &author.Biography, &author.BirthDate)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return author, err
	}
	return author, nil
}

func (b *AuthorPostgres) Update(tx *sql.Tx, author *Author) (sql.Result, error) {
	log.Println("AuthorPostgres. Update")

	query := fmt.Sprintf("update %s set FirstName = $2, LastName = $3, Biography = $4, BirthDate = $5 where Id = $1", authorTable)
	return tx.Exec(query, author.ID, author.FirstName, author.LastName, author.Biography, author.BirthDate)
}

func (b *AuthorPostgres) Delete(tx *sql.Tx, id int) (sql.Result, error) {
	log.Println("AuthorPostgres. Delete")

	query := fmt.Sprintf("delete from %s where Id = $1", authorTable)
	return tx.Exec(query, id)
}
