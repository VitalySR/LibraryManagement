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
	row := tx.QueryRow(query, bk.Title, authorId, bk.Year, bk.ISBN)
	err = row.Scan(&id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (b *BookPostgres) GetAll() ([]Book, error) {
	log.Println("BookPostgres. GetAll")
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
			return nil, err
		}
		if author.ID != nil {
			bk.Author = &author
		}
		bks = append(bks, bk)
	}
	return bks, nil
}

func (b *BookPostgres) GetById(id int) (Book, error) {
	log.Println("BookPostgres. GetById")
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

func (b *BookPostgres) Update(bk Book) (int64, error) {
	log.Println("BookPostgres. Update")
	var rowCnt int64 = 0
	tx, err := b.db.Begin()
	if err != nil {
		return 0, err
	}

	query := fmt.Sprintf("update %s set Title = $2, Author_Id = $3, Year = $4, ISBN = $5 where Id = $1", bookTable)
	var authorId *int32 = nil
	if bk.Author != nil && bk.Author.ID != nil {
		authorId = bk.Author.ID
	}
	result, err := tx.Exec(query, bk.ID, bk.Title, authorId, bk.Year, bk.ISBN)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if result != nil {
		rowCnt, _ = result.RowsAffected()
	}
	return rowCnt, tx.Commit()
}

func (b *BookPostgres) Delete(id int) (int64, error) {
	log.Println("BookPostgres. Delete")
	var rowCnt int64 = 0
	tx, err := b.db.Begin()
	if err != nil {
		return 0, err
	}
	query := fmt.Sprintf("delete from %s where Id = $1", bookTable)
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
