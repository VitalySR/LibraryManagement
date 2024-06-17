package repository

import (
	"database/sql"
	"library/entity"
)

type BookWorker interface {
	Create(entity.Book) (int, error)
	GetAll() ([]entity.Book, error)
	GetById(id int) (entity.Book, error)
	Update(entity.Book) error
	Delete(id int) error
}

type AuthorWorker interface {
	Create(author entity.Author) (int, error)
	GetAll() ([]entity.Author, error)
	GetById(id int) (entity.Author, error)
	Update(author entity.Author) error
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
