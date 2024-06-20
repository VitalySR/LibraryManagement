package repository

import (
	"database/sql"
)

type BookWorker interface {
	Create(*sql.Tx, *Book) (int, error)
	GetAll(*sql.DB) ([]Book, error)
	GetById(*sql.DB, int) (Book, error)
	Update(*sql.Tx, *Book) (sql.Result, error)
	Delete(*sql.Tx, int) (sql.Result, error)
}

type AuthorWorker interface {
	Create(*sql.Tx, *Author) (int, error)
	GetAll(*sql.DB) ([]Author, error)
	GetById(*sql.DB, int) (Author, error)
	Update(*sql.Tx, *Author) (sql.Result, error)
	Delete(*sql.Tx, int) (sql.Result, error)
}

type Repository struct {
	BookWorker   BookWorker
	AuthorWorker AuthorWorker
	db           *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		BookWorker:   NewBookPostgres(),
		AuthorWorker: NewAuthorPostgres(),
		db:           db,
	}
}

func (repo *Repository) CreateBook(book *Book) (int, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return 0, err
	}

	id, err := repo.BookWorker.Create(tx, book)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (repo *Repository) GetAllBook() ([]Book, error) {
	return repo.BookWorker.GetAll(repo.db)
}

func (repo *Repository) GetBookById(id int) (Book, error) {
	return repo.BookWorker.GetById(repo.db, id)
}

func (repo *Repository) UpdateBook(book *Book) (int64, error) {
	var rowCnt int64 = 0

	tx, err := repo.db.Begin()
	if err != nil {
		return 0, err
	}

	result, err := repo.BookWorker.Update(tx, book)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}
	if result != nil {
		rowCnt, _ = result.RowsAffected()
	}
	return rowCnt, tx.Commit()
}

func (repo *Repository) DeleteBook(id int) (int64, error) {
	var rowCnt int64 = 0
	tx, err := repo.db.Begin()
	if err != nil {
		return 0, err
	}
	result, err := repo.BookWorker.Delete(tx, id)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}
	if result != nil {
		rowCnt, _ = result.RowsAffected()
	}
	return rowCnt, tx.Commit()
}

func (repo *Repository) CreateAuthor(author *Author) (int, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return 0, err
	}

	id, err := repo.AuthorWorker.Create(tx, author)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (repo *Repository) GetAllAuthor() ([]Author, error) {
	return repo.AuthorWorker.GetAll(repo.db)
}

func (repo *Repository) GetAuthorById(id int) (Author, error) {
	return repo.AuthorWorker.GetById(repo.db, id)
}

func (repo *Repository) UpdateAuthor(author *Author) (int64, error) {
	var rowCnt int64 = 0

	tx, err := repo.db.Begin()
	if err != nil {
		return 0, err
	}

	result, err := repo.AuthorWorker.Update(tx, author)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}
	if result != nil {
		rowCnt, _ = result.RowsAffected()
	}
	return rowCnt, tx.Commit()
}

func (repo *Repository) DeleteAuthor(id int) (int64, error) {
	var rowCnt int64 = 0
	tx, err := repo.db.Begin()
	if err != nil {
		return 0, err
	}
	result, err := repo.AuthorWorker.Delete(tx, id)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}
	if result != nil {
		rowCnt, _ = result.RowsAffected()
	}
	return rowCnt, tx.Commit()
}
