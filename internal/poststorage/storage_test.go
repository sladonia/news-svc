package poststorage

import (
	"database/sql"
	"testing"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/lib/pq"
	"github.com/sladonia/news-svc/internal/post"
	"github.com/stretchr/testify/suite"
)

const (
	testDBDSN     = "postgresql://user:password@localhost:5432/news?sslmode=disable"
	postTableName = "post"
)

// db fixtures
var (
	post1 = post.Post{
		ID:        "1",
		Title:     "exclusive",
		Content:   "new era beginning!",
		CreatedAt: time.Now().UTC().Round(time.Millisecond),
		UpdatedAt: time.Now().UTC().Round(time.Millisecond),
	}
)

type Suite struct {
	suite.Suite

	db      *goqu.Database
	storage post.Storage
}

func (s *Suite) SetupSuite() {
	postgresClient, err := sql.Open("postgres", testDBDSN)
	if err != nil {
		panic(err)
	}

	s.db = goqu.New("postgres", postgresClient)
	s.storage = New(s.db, postTableName)

	loc, err := time.LoadLocation("UTC")
	if err != nil {
		panic(err)
	}
	goqu.SetTimeLocation(loc)

	// TODO: run migration
}

func TestRunSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) SetupTest() {
	err := s.insertFixtures()
	if err != nil {
		panic(err)
	}
}

func (s *Suite) TearDownTest() {
	err := s.deleteAllData()
	if err != nil {
		panic(err)
	}
}

func (s *Suite) TestInsertPost() {
	s.Run("success", func() {
		p2 := post.NewPost("title2", "content2")
		p2.ID = "2"

		err := s.storage.Insert(p2)
		s.NoError(err)

		fromStorage, err := s.storage.ByID("2")

		s.NoError(err)
		s.Equal(p2, fromStorage)
	})

	s.Run("conflict", func() {
		err := s.storage.Insert(post1)
		s.Error(err)
		s.ErrorIs(err, post.ErrorAlreadyExists)
	})
}

func (s *Suite) TestByID() {
	s.Run("success", func() {
		p, err := s.storage.ByID("1")

		s.NoError(err)
		s.Equal(post1, p)
	})

	s.Run("no_documents", func() {
		_, err := s.storage.ByID("unexisting")

		s.Error(err)
		s.ErrorIs(err, post.ErrNotFound)
	})
}

func (s *Suite) TestByFilter() {
	s.Run("latest", func() {
		f := post.Filter{
			Limit:  10,
			Offset: 0,
		}

		posts, err := s.storage.ByFilter(f)
		s.NoError(err)
		s.Len(posts, 1)
	})

	s.Run("after_date", func() {
		f := post.Filter{
			From: time.Now().Add(-time.Hour),
			Limit:  10,
			Offset: 0,
		}

		posts, err := s.storage.ByFilter(f)
		s.NoError(err)
		s.Len(posts, 1)
	})

	s.Run("no_results_found", func() {
		f := post.Filter{
			From: time.Now().Add(-2*time.Hour),
			To: time.Now().Add(-1*time.Hour),
			Limit:  10,
			Offset: 0,
		}

		posts, err := s.storage.ByFilter(f)
		s.NoError(err)
		s.Len(posts, 0)
	})
}

func (s *Suite) TestReplace() {
	s.Run("replace", func() {
		updatedContent := "content2"

		updated := post1
		updated.Content = updatedContent

		err := s.storage.Update(updated)
		s.NoError(err)

		retrieved, err := s.storage.ByID(post1.ID)
		s.NoError(err)
		s.Equal(updated.Content, retrieved.Content)
	})

	s.Run("not_found", func() {
		newPost := post.NewPost("extra new", "something happened")

		err := s.storage.Update(newPost)
		s.Error(err)
		s.ErrorIs(err, post.ErrNotFound)
	})
}

func (s *Suite) TestRemove() {
	s.Run("no_documents", func() {
		err := s.storage.Remove("42")
		s.NoError(err)
	})

	s.Run("success", func() {
		err := s.storage.Remove("1")
		s.NoError(err)

		_, err = s.storage.ByID("1")
		s.ErrorIs(err, post.ErrNotFound)
	})
}

func (s *Suite) insertFixtures() error {
	post1SQL := NewPostSQL(post1)
	_, err := s.db.Insert(postTableName).Rows(post1SQL).Executor().Exec()

	return err
}

func (s *Suite) deleteAllData() error {
	_, err := s.db.Delete(postTableName).Executor().Exec()

	return err
}
