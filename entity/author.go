package entity

type Author struct {
	ID        int32  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Biography string `json:"biography,omitempty"`
	BirthDate string `json:"birth_date,omitempty"`
}
