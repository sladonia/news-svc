package poststorage

import (
	"time"

	"github.com/sladonia/news-svc/internal/post"
)

const (
	columnID        = "id"
	columnTitle     = "title"
	columnContent   = "content"
	columnCreatedAt = "created_at"
	columnUpdatedAt = "updated_at"
)

type PostSQL struct {
	ID        string    `db:"id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewPostSQL(post post.Post) PostSQL {
	return PostSQL{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}
}

func NewPostFromSQL(postSQL PostSQL) post.Post {
	return post.Post{
		ID:        postSQL.ID,
		Title:     postSQL.Title,
		Content:   postSQL.Content,
		CreatedAt: postSQL.CreatedAt.UTC(),
		UpdatedAt: postSQL.UpdatedAt.UTC(),
	}
}
