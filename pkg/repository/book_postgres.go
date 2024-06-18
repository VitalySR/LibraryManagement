package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type BookPostgres struct {
	db *sql.DB
}

func NewBookPostgres(db *sql.DB) *BookPostgres {
	return &BookPostgres{db: db}
}

func (b *BookPostgres) Create(bk Book) (int, error) {
	log.Println("BookPostgres. Create")
	tx, err := b.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	var authorId *int32 = nil
	if bk.Author != nil && bk.Author.ID != nil {
		authorId = bk.Author.ID
	}
	query := fmt.Sprintf("INSERT INTO %s (Title, Author_Id, Year, ISBN) values ($1, $2, $3, $4) RETURNING Id", bookTable)
	row := tx.QueryRow(query, bk.Title, Int32WithNull(authorId), Int32WithNull(bk.Year), StringWithNull(bk.ISBN))
	err = row.Scan(&id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (b *BookPostgres) GetAll() ([]Book, error) {
	query := fmt.Sprintf("select b.id, b.title, b.year, b.isbn, b.author_id, a.FirstName, a.LastName, a.Biography, a.BirthDate from %s b left join %s a on a.ID = b.Author_Id", bookTable, authorTable)
	rows, err := b.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	bks := make([]Book, 0, 5)
	for rows.Next() {
		bk := Book{}
		author := Author{}
		if err := rows.Scan(&bk.ID, &bk.Title, &bk.Year, &bk.ISBN, &author.ID, &author.FirstName, &author.LastName, &author.Biography, &author.BirthDate); err != nil {
			return bks, err
		}
		if author.ID != nil {
			bk.Author = &author
		}
		bks = append(bks, bk)
	}
	return bks, nil
}

func (b *BookPostgres) GetById(id int) (Book, error) {
	query := fmt.Sprintf("select b.id, b.title, b.year, b.isbn, b.author_id, a.FirstName, a.LastName, a.Biography, a.BirthDate from %s b left join %s a on a.ID = b.Author_Id where b.ID = $1", bookTable, authorTable)
	row := b.db.QueryRow(query, id)
	bk := Book{}
	author := Author{}
	err := row.Scan(&bk.ID, &bk.Title, &bk.Year, &bk.ISBN, &author.ID, &author.FirstName, &author.LastName, &author.Biography, &author.BirthDate)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return bk, err
	}
	if author.ID != nil {
		bk.Author = &author
	}
	return bk, nil
}

func (b *BookPostgres) Update(Book) error {
	return nil
}

func (b *BookPostgres) Delete(id int) (sql.Result, error) {
	query := fmt.Sprintf("delete from %s where Id = $1", bookTable)
	return b.db.Exec(query, id)
}
