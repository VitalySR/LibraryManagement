package repository

import (
	"database/sql"
	"fmt"
)

type AuthorPostgres struct {
	db *sql.DB
}

func NewAuthorPostgres(db *sql.DB) *AuthorPostgres {
	return &AuthorPostgres{db: db}
}

func (b *AuthorPostgres) Create(au Author) (int, error) {
	tx, err := b.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (FirstName, LastName, Biography, BirthDate) values ($1, $2, $3, $4) RETURNING id", authorTable)

	row := tx.QueryRow(createItemQuery, au.FirstName, au.LastName, StringWithNull(au.Biography), StringWithNull(au.BirthDate))
	err = row.Scan(&id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (b *AuthorPostgres) GetAll() ([]Author, error) {
	return nil, nil
}

func (b *AuthorPostgres) GetById(id int) (Author, error) {
	return Author{}, nil
}

func (b *AuthorPostgres) Update(Author) error {
	return nil
}

func (b *AuthorPostgres) Delete(id int) error {
	return nil
}
