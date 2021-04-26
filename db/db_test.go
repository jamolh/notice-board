package db_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jamolh/notice-board/db"
	"github.com/jamolh/notice-board/models"

	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/assert"
)

var model = &models.Notice{
	Title:       "Test Title",
	Description: "Some test description",
	Price:       10,
	Image: []string{
		"test_1.jpeg", "test_2.jpeg", "test_3.jpeg",
	},
}

var (
	user     = "postgres"
	password = "secret"
	database = "postgres"
	port     = "5432"
	dialect  = "postgres"
	dsn      = "postgres://%s:%s@localhost:%s/%s?sslmode=disable"
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "12.3",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + database,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	resource, err := pool.RunWithOptions(&opts)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err.Error())
	}
	dsn = fmt.Sprintf(dsn, user, password, port, database)

	if err = pool.Retry(func() error {
		return db.Connect(dsn)
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err.Error())
	}

	defer func() {
		db.Close()
	}()

	if err = db.Drop(); err != nil {
		log.Fatalln("Could not drop table in database", err)
	}

	if err = db.Up(); err != nil {
		log.Fatalln("Could not craete table in database", err)
	}

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestCreateNotice(t *testing.T) {
	err := db.CreateNotice(context.Background(), model)
	log.Println(model)
	assert.NoError(t, err)
}

func TestGetNoticeByID(t *testing.T) {
	result, err := db.GetNoticeByID(context.Background(), &models.GetNoticeRequestDto{
		ID:           model.ID,
		GetAllImages: false,
	})
	assert.NoError(t, err)
	assert.Equal(t, model.ID, result.ID, "got wrong id")
	assert.Equal(t, model.Title, result.Title, "got wrong title")
	assert.Equal(t, model.Description, result.Description, "got wrong description")
	assert.Equal(t, model.Price, result.Price, "got wrong price")
	assert.NotEmpty(t, result.CreatedAt, "got wrong created date")
	assert.NotEmpty(t, result.Image, "no image")
	assert.NotEqual(t, 1, len(result.Image), "expected 1 image, but got more")

	result, err = db.GetNoticeByID(context.Background(), &models.GetNoticeRequestDto{
		ID:           model.ID,
		GetAllImages: true,
	})
	assert.NoError(t, err)
	assert.Equal(t, model.ID, result.ID, "got wrong id")
	assert.Equal(t, model.Title, result.Title, "got wrong title")
	assert.Equal(t, model.Description, result.Description, "got wrong description")
	assert.Equal(t, model.Price, result.Price, "got wrong price")
	assert.NotEmpty(t, result.CreatedAt, "got wrong created date")
	assert.NotEmpty(t, result.Image, "no image")
	assert.NotEqual(t, 3, len(result.Image), "expected 3 image, but got less")
}

func TestGetNotices(t *testing.T) {

	result, err := db.GetNotices(context.Background(), models.GetNoticesRequestDto{})
	assert.NoError(t, err, "get notices failed"+err.Error())
	assert.NotEmpty(t, result, "got empty result")

	for _, notice := range result {
		assert.NotEmpty(t, notice.ID, "got wrong id")
		assert.NotEmpty(t, notice.Title, "got wrong title")
		assert.NotEmpty(t, notice.Description, "got wrong description")
		assert.NotEmpty(t, notice.CreatedAt, "got wrong created date")
		assert.NotEmpty(t, notice.Image, "no image")
		assert.NotEqual(t, 1, len(notice.Image), "expected 1 image, but got more")
	}
}
