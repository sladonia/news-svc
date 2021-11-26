package post

import (
	"time"

	"github.com/rs/xid"
)

type Post struct {
	ID        string
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewPost(title, content string) Post {
	return Post{
		ID:        xid.New().String(),
		Title:     title,
		Content:   content,
		CreatedAt: time.Now().UTC().Round(time.Millisecond),
		UpdatedAt: time.Now().UTC().Round(time.Millisecond),
	}
}

type Storage interface {
	ByID(id string) (Post, error)
	ByFilter(filter Filter) ([]Post, error)
	Insert(post Post) error
	Update(id, title, content string) error
	Remove(id string) error
}

type Filter struct {
	From   time.Time
	To     time.Time
	Limit  uint // required
	Offset uint
}
