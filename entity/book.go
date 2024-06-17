package entity

type Book struct {
	ID     int32  `json:"id"`
	Title  string `json:"title"`
	Year   int32  `json:"year,omitempty"`
	ISBN   string `json:"isbn,omitempty"`
	Author `json:"author,omitempty"`
}
