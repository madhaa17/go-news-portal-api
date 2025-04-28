package entity

type CategoryEntity struct {
	ID    int16
	Title string
	Slug  string
	UserEntity
}
