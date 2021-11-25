package post

import "time"

type Post struct {
	ID        string
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Storage interface {
	ByID(id string) (Post, error)
	ByFilter(filter Filter) ([]Post, error)
	Insert(post Post) error
	Replace(post Post) error
	Remove(post Post) error
}

type Filter struct {
	From   time.Time
	To     time.Time
	Limit  int
	Offset int
}
