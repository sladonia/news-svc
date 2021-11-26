package testtool

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/doug-martin/goqu/v9"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

const (
	password = "password"
	user     = "user"
)

func NewPostgresContainer(pool *dockertest.Pool, dbName string) (*dockertest.Resource, string, error) {
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14.1-alpine",
		Env: []string{
			fmt.Sprintf("POSTGRES_PASSWORD=%s", password),
			fmt.Sprintf("POSTGRES_USER=%s", user),
			fmt.Sprintf("POSTGRES_DB=%s", dbName),
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		return nil, "", err
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	dbDSN := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", user, password, hostAndPort, dbName)

	return resource, dbDSN, nil
}

func ApplyDBMigrations(db *goqu.Database, migrationsDIR string) error {
	fileInfos, err := ioutil.ReadDir(migrationsDIR)
	if err != nil {
		return err
	}

	sort.SliceStable(fileInfos, func(i, j int) bool {
		return fileInfos[i].Name() < fileInfos[j].Name()
	})

	for _, fileInfo := range fileInfos {
		filePath := fmt.Sprintf("%s/%s", migrationsDIR, fileInfo.Name())

		err = runMigrationFile(db, filePath)
		if err != nil {
			return err
		}
	}

	return nil
}

func runMigrationFile(db *goqu.Database, filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	sqlStr, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(sqlStr))

	return err
}

