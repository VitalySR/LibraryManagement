package repository

type Book struct {
	ID     *int32  `json:"id"`
	Title  *string `json:"title"`
	Author *Author `json:"author,omitempty"`
	Year   *int32  `json:"year,omitempty"`
	ISBN   *string `json:"isbn,omitempty"`
}

type Author struct {
	ID        *int32  `json:"id"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Biography *string `json:"biography,omitempty"`
	BirthDate *string `json:"birth_date,omitempty"`
}
