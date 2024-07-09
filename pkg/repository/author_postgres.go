package repository

import (
	"database/sql"
	"errors"
	"log"
)

type AuthorPostgres struct{}

const (
	insertAuthor  = "INSERT INTO author (FirstName, LastName, Biography, BirthDate) values ($1, $2, $3, $4) RETURNING id"
	selectAuthors = "SELECT Id, FirstName, LastName, Biography, BirthDate FROM author"
	selectAuthor  = "SELECT Id, FirstName, LastName, Biography, BirthDate FROM author WHERE Id = $1"
	updateAuthor  = "UPDATE author SET FirstName = $2, LastName = $3, Biography = $4, BirthDate = $5 WHERE Id = $1"
	deleteAuthor  = "DELETE FROM author WHERE Id = $1"
)

func NewAuthorPostgres() *AuthorPostgres {
	return &AuthorPostgres{}
}

func (b *AuthorPostgres) Create(tx *sql.Tx, author *Author) (int, error) {
	log.Println("AuthorPostgres. Create")

	var id int
	row := tx.QueryRow(insertAuthor, author.FirstName, author.LastName, author.Biography, author.BirthDateTime())
	err := row.Scan(&id)

	return id, err
}

func (b *AuthorPostgres) GetAll(db *sql.DB) ([]Author, error) {
	log.Println("AuthorPostgres. GetAll")
	rows, err := db.Query(selectAuthors)
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
	row := db.QueryRow(selectAuthor, id)
	author := Author{}
	err := row.Scan(&author.ID, &author.FirstName, &author.LastName, &author.Biography, &author.BirthDate)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return author, err
	}
	return author, nil
}

func (b *AuthorPostgres) Update(tx *sql.Tx, author *Author) (sql.Result, error) {
	log.Println("AuthorPostgres. Update")

	return tx.Exec(updateAuthor, author.ID, author.FirstName, author.LastName, author.Biography, author.BirthDateTime())
}

func (b *AuthorPostgres) Delete(tx *sql.Tx, id int) (sql.Result, error) {
	log.Println("AuthorPostgres. Delete")

	return tx.Exec(deleteAuthor, id)
}
