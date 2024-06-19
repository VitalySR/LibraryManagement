package repository

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
