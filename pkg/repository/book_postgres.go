package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type BookPostgres struct{}

func NewBookPostgres() *BookPostgres {
	return &BookPostgres{}
}

func (b *BookPostgres) Create(tx *sql.Tx, bk *Book) (int, error) {
	log.Println("BookPostgres. Create")

	var id int
	var authorId *int32 = nil
	if bk.Author != nil && bk.Author.ID != nil {
		authorId = bk.Author.ID
	}

	query := fmt.Sprintf("INSERT INTO %s (Title, Author_Id, Year, ISBN) values ($1, $2, $3, $4) RETURNING Id", bookTable)
	row := tx.QueryRow(query, bk.Title, authorId, bk.Year, bk.ISBN)
	err := row.Scan(&id)

	return id, err
}

func (b *BookPostgres) GetAll(db *sql.DB) ([]Book, error) {
	log.Println("BookPostgres. GetAll")
	query := fmt.Sprintf("select b.id, b.title, b.year, b.isbn, b.author_id, a.FirstName, a.LastName, a.Biography, a.BirthDate from %s b left join %s a on a.ID = b.Author_Id", bookTable, authorTable)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	bks := make([]Book, 0, 5)
	for rows.Next() {
		bk := Book{}
		author := Author{}
		if err := rows.Scan(&bk.ID, &bk.Title, &bk.Year, &bk.ISBN, &author.ID, &author.FirstName, &author.LastName, &author.Biography, &author.BirthDate); err != nil {
			return nil, err
		}
		if author.ID != nil {
			bk.Author = &author
		}
		bks = append(bks, bk)
	}
	return bks, nil
}

func (b *BookPostgres) GetById(db *sql.DB, id int) (Book, error) {
	log.Println("BookPostgres. GetById")
	query := fmt.Sprintf("select b.id, b.title, b.year, b.isbn, b.author_id, a.FirstName, a.LastName, a.Biography, a.BirthDate from %s b left join %s a on a.ID = b.Author_Id where b.ID = $1", bookTable, authorTable)
	row := db.QueryRow(query, id)
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

func (b *BookPostgres) Update(tx *sql.Tx, bk *Book) (sql.Result, error) {
	log.Println("BookPostgres. Update")

	query := fmt.Sprintf("update %s set Title = $2, Author_Id = $3, Year = $4, ISBN = $5 where Id = $1", bookTable)
	var authorId *int32 = nil
	if bk.Author != nil && bk.Author.ID != nil {
		authorId = bk.Author.ID
	}
	return tx.Exec(query, bk.ID, bk.Title, authorId, bk.Year, bk.ISBN)
}

func (b *BookPostgres) Delete(tx *sql.Tx, id int) (sql.Result, error) {
	log.Println("BookPostgres. Delete")

	query := fmt.Sprintf("delete from %s where Id = $1", bookTable)
	return tx.Exec(query, id)
}
