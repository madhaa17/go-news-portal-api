package entity

type CategoryEntity struct {
	ID    int16
	Title string
	Slug  string
	UserEntity
}

// Error implements error.
func (c *CategoryEntity) Error() string {
	panic("unimplemented")
}
