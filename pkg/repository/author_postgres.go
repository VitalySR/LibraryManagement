package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type AuthorPostgres struct {
	db *sql.DB
}

func NewAuthorPostgres(db *sql.DB) *AuthorPostgres {
	return &AuthorPostgres{db: db}
}

func (b *AuthorPostgres) Create(au Author) (int, error) {
	log.Println("AuthorPostgres. Create")
	tx, err := b.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (FirstName, LastName, Biography, BirthDate) values ($1, $2, $3, $4) RETURNING id", authorTable)

	row := tx.QueryRow(createItemQuery, au.FirstName, au.LastName, au.Biography, au.BirthDate)
	err = row.Scan(&id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (b *AuthorPostgres) GetAll() ([]Author, error) {
	log.Println("AuthorPostgres. GetAll")
	query := fmt.Sprintf("select Id, FirstName, LastName, Biography, BirthDate from %s", authorTable)
	rows, err := b.db.Query(query)
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

func (b *AuthorPostgres) GetById(id int) (Author, error) {
	log.Println("AuthorPostgres. GetById")
	query := fmt.Sprintf("select Id, FirstName, LastName, Biography, BirthDate from %s where Id = $1", authorTable)
	row := b.db.QueryRow(query, id)
	author := Author{}
	err := row.Scan(&author.ID, &author.FirstName, &author.LastName, &author.Biography, &author.BirthDate)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return author, err
	}
	return author, nil
}

func (b *AuthorPostgres) Update(au Author) (int64, error) {
	log.Println("AuthorPostgres. Update")
	var rowCnt int64 = 0
	tx, err := b.db.Begin()
	if err != nil {
		return 0, err
	}
	query := fmt.Sprintf("update %s set FirstName = $2, LastName = $3, Biography = $4, BirthDate = $5 where Id = $1", authorTable)
	result, err := tx.Exec(query, au.ID, au.FirstName, au.LastName, au.Biography, au.BirthDate)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if result != nil {
		rowCnt, _ = result.RowsAffected()
	}
	return rowCnt, tx.Commit()
}

func (b *AuthorPostgres) Delete(id int) (int64, error) {
	log.Println("AuthorPostgres. Delete")
	var rowCnt int64 = 0
	tx, err := b.db.Begin()
	if err != nil {
		return 0, err
	}
	query := fmt.Sprintf("delete from %s where Id = $1", authorTable)
	result, err := tx.Exec(query, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if result != nil {
		rowCnt, _ = result.RowsAffected()
	}
	return rowCnt, tx.Commit()
}
