package repository

import (
	"errors"
)

type Book struct {
	ID     *int32  `json:"id,omitempty"`
	Title  *string `json:"title,omitempty"`
	Author *Author `json:"author,omitempty"`
	Year   *int32  `json:"year,omitempty"`
	ISBN   *string `json:"isbn,omitempty"`
}

type Author struct {
	ID        *int32  `json:"id,omitempty"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Biography *string `json:"biography,omitempty"`
	BirthDate *string `json:"birth_date,omitempty"`
}

func (b Book) Validate(withAuthor bool) error {
	if b.Title == nil || len(*b.Title) == 0 {
		return errors.New("book's title is mandatory field")
	}

	if withAuthor {
		if b.Author == nil {
			return errors.New("author is mandatory")
		}
		if err := b.Author.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (a Author) Validate() error {
	if a.FirstName == nil || len(*a.FirstName) == 0 {
		return errors.New("author's first_name is mandatory fields")
	}

	if a.LastName == nil || len(*a.LastName) == 0 {
		return errors.New("author's last_name is mandatory fields")
	}
	return nil
}
