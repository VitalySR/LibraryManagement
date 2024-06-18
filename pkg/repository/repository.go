package repository

import (
	"database/sql"
)

type BookWorker interface {
	Create(Book) (int, error)
	GetAll() ([]Book, error)
	GetById(id int) (Book, error)
	Update(Book) error
	Delete(id int) (sql.Result, error)
}

type AuthorWorker interface {
	Create(author Author) (int, error)
	GetAll() ([]Author, error)
	GetById(id int) (Author, error)
	Update(author Author) error
	Delete(id int) error
}

type Repository struct {
	BookWorker
	AuthorWorker
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		BookWorker:   NewBookPostgres(db),
		AuthorWorker: NewAuthorPostgres(db),
	}
}
