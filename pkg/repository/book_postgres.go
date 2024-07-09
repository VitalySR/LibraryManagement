package repository

import (
	"database/sql"
	"errors"
	"log"
)

type BookPostgres struct{}

const (
	insertBook  = "INSERT INTO book (Title, Author_Id, Year, ISBN) values ($1, $2, $3, $4) RETURNING Id"
	selectBooks = "SELECT b.id, b.title, b.year, b.isbn, b.author_id, a.FirstName, a.LastName, a.Biography, a.BirthDate FROM book b LEFT JOIN author a ON a.ID = b.Author_Id"
	selectBook  = "SELECT b.id, b.title, b.year, b.isbn, b.author_id, a.FirstName, a.LastName, a.Biography, a.BirthDate FROM book b LEFT JOIN author a ON a.ID = b.Author_Id WHERE b.ID = $1"
	updateBook  = "UPDATE book SET Title = $2, Author_Id = $3, Year = $4, ISBN = $5 WHERE Id = $1"
	deleteBook  = "DELETE FROM book WHERE Id = $1"
)

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

	row := tx.QueryRow(insertBook, bk.Title, authorId, bk.Year, bk.ISBN)
	err := row.Scan(&id)

	return id, err
}

func (b *BookPostgres) GetAll(db *sql.DB) ([]Book, error) {
	log.Println("BookPostgres. GetAll")
	rows, err := db.Query(selectBooks)
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
	row := db.QueryRow(selectBook, id)
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

	var authorId *int32 = nil
	if bk.Author != nil && bk.Author.ID != nil {
		authorId = bk.Author.ID
	}
	return tx.Exec(updateBook, bk.ID, bk.Title, authorId, bk.Year, bk.ISBN)
}

func (b *BookPostgres) Delete(tx *sql.Tx, id int) (sql.Result, error) {
	log.Println("BookPostgres. Delete")

	return tx.Exec(deleteBook, id)
}
