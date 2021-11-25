package poststorage

import (
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/lib/pq"
	"github.com/sladonia/news-svc/internal/post"
)

func New(db *goqu.Database, postTableName string) post.Storage {
	return &storage{
		db:            db,
		postTableName: postTableName,
	}
}

type storage struct {
	db            *goqu.Database
	postTableName string
}

func (s *storage) ByID(id string) (post.Post, error) {
	query := s.db.From(s.postTableName).
		Where(goqu.C(columnID).Eq(id))

	var p PostSQL

	ok, err := query.ScanStruct(&p)
	if err != nil {
		return post.Post{}, err
	}

	if !ok {
		return post.Post{}, post.ErrNotFound
	}

	return NewPostFromSQL(p), nil
}

func (s *storage) ByFilter(filter post.Filter) ([]post.Post, error) {
	q := s.db.From(s.postTableName)

	if !filter.From.IsZero() {
		q = q.Where(goqu.C(columnCreatedAt).Gte(filter.From))
	}

	if !filter.To.IsZero() {
		q = q.Where(goqu.C(columnCreatedAt).Lte(filter.To))
	}

	q = q.Order(goqu.C(columnCreatedAt).Desc()).
		Limit(filter.Limit).
		Offset(filter.Offset)

	var postsSQL []PostSQL

	err := q.ScanStructs(&postsSQL)
	if err != nil {
		return nil, err
	}

	posts := make([]post.Post, len(postsSQL))

	for i, postSQL := range postsSQL {
		posts[i] = NewPostFromSQL(postSQL)
	}

	return posts, nil
}

func (s *storage) Insert(p post.Post) error {
	postSQL := NewPostSQL(p)

	_, err := s.db.Insert(s.postTableName).Rows(postSQL).Executor().Exec()

	var errDuplicate *pq.Error

	if errors.As(err, &errDuplicate) {
		if errDuplicate.Code == "23505" {
			return post.ErrorAlreadyExists
		}
	}

	return err
}

func (s *storage) Update(p post.Post) error {
	postSQL := NewPostSQL(p)

	res, err := s.db.Update(s.postTableName).
		Where(goqu.C(columnID).Eq(p.ID)).
		Set(postSQL).
		Executor().
		Exec()

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return post.ErrNotFound
	}

	return err
}

func (s *storage) Remove(id string) error {
	_, err := s.db.Delete(s.postTableName).Where(goqu.C(columnID).Eq(id)).Executor().Exec()

	return err
}
