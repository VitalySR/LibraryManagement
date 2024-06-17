package repository

import (
	"database/sql"
	"fmt"
	"library/entity"
	"log"
)

type BookPostgres struct {
	db *sql.DB
}

func NewBookPostgres(db *sql.DB) *BookPostgres {
	return &BookPostgres{db: db}
}

func (b *BookPostgres) Create(bk entity.Book) (int, error) {
	log.Println("BookPostgres. Create")
	tx, err := b.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	query := fmt.Sprintf("INSERT INTO %s (Title, Author_Id, Year, ISBN) values ($1, $2, $3, $4) RETURNING Id", bookTable)
	row := tx.QueryRow(query, bk.Title, Int32WithNull(bk.Author.ID), Int32WithNull(bk.Year), StringWithNull(bk.ISBN))
	err = row.Scan(&id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (b *BookPostgres) GetAll() ([]entity.Book, error) {
	query := fmt.Sprintf("select b.id, b.title, b.year, b.isbn, b.author_id, a.firstname, a.lastname, a.biography, a.birthdate"+
		" from %s b left join %s a on b.author_id = a.id", bookTable, authorTable)
	rows, err := b.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	bks := make([]entity.Book, 0, 5)
	for rows.Next() {
		bk := entity.Book{}
		log.Println(bk)
		if err := rows.Scan(&bk.ID, &bk.Title, &bk.Year, &bk.ISBN, &bk.Author.ID, &bk.Author.FirstName, &bk.Author.LastName, &bk.Author.Biography, &bk.Author.BirthDate); err != nil {
			return bks, err
		}
		bks = append(bks, bk)
	}
	return bks, nil
}

func (b *BookPostgres) GetById(id int) (entity.Book, error) {
	return entity.Book{}, nil
}

func (b *BookPostgres) Update(entity.Book) error {
	return nil
}

func (b *BookPostgres) Delete(id int) error {
	return nil
}
