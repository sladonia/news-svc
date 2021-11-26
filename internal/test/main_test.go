package test

import (
	"database/sql"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/gorilla/mux"
	"github.com/ory/dockertest/v3"
	"github.com/sladonia/news-svc/internal/handler"
	"github.com/sladonia/news-svc/internal/logger"
	"github.com/sladonia/news-svc/internal/post"
	"github.com/sladonia/news-svc/internal/poststorage"
	"github.com/sladonia/news-svc/internal/testtool"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

const (
	postTableName = "post"
	testDBNAME    = "news_test"
	migrationsDir = "../../migration"
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

	log               *zap.Logger
	srv               *httptest.Server
	db                *goqu.Database
	storage           post.Storage
	service           post.Service
	handler           *handler.Handler
	dockerPool        *dockertest.Pool
	postgresContainer *dockertest.Resource
}

func (s *Suite) SetupSuite() {
	log, err := logger.NewZapLogger("debug")
	if err != nil {
		panic(err)
	}
	s.log = log

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Panic("failed to init pool", zap.Error(err))
	}
	s.dockerPool = pool

	postgresContainer, dbDSN, err := testtool.NewPostgresContainer(pool, testDBNAME)
	if err != nil {
		log.Panic("failed to create postgres container", zap.Error(err))
	}
	s.postgresContainer = postgresContainer

	var postgresClient *sql.DB

	pool.MaxWait = 120 * time.Second
	err = pool.Retry(func() error {
		postgresClient, err = sql.Open("postgres", dbDSN)
		if err != nil {
			log.Panic("failed to connect to db", zap.Error(err))
		}

		return postgresClient.Ping()
	})
	if err != nil {
		log.Panic("database up timeout", zap.Error(err))
	}

	s.db = goqu.New("postgres", postgresClient)
	s.storage = poststorage.New(s.db, postTableName)
	s.service = post.NewService(s.storage)
	s.handler = handler.NewHandler(log, 100, s.service, "news-sv")

	r := mux.NewRouter()
	s.handler.Register(r)

	s.srv = httptest.NewServer(r)

	loc, err := time.LoadLocation("UTC")
	if err != nil {
		panic(err)
	}
	goqu.SetTimeLocation(loc)

	err = testtool.ApplyDBMigrations(s.db, migrationsDir)
	if err != nil {
		panic(err)
	}

}

func (s *Suite) TearDownSuite() {
	s.srv.Close()

	err := s.dockerPool.Purge(s.postgresContainer)
	if err != nil {
		panic(err)
	}
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

func (s *Suite) insertFixtures() error {
	post1SQL := poststorage.NewPostSQL(post1)
	_, err := s.db.Insert(postTableName).Rows(post1SQL).Executor().Exec()

	return err
}

func (s *Suite) deleteAllData() error {
	_, err := s.db.Delete(postTableName).Executor().Exec()

	return err
}
