package db

import (
	"context"
	"log"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	pool *pgxpool.Pool
)

func init() {
	// app logging to file
	log.SetOutput(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    2, // megabytes
		MaxBackups: 30,
		MaxAge:     40,   //days
		Compress:   true, // disabled by default
	})

	log.Println("-------- * ------- Starting Logging -------- * -------")
}

// Connect to database
func Connect(dbConnection string) (err error) {
	pool, err = pgxpool.Connect(context.Background(), dbConnection)
	if err != nil {
		log.Println("Connecting database failed:", err)
	}
	log.Println("ooooo")
	return
}

func Up() (err error) {
	_, err = pool.Exec(context.Background(), createTableNoticesQuery)
	if err != nil {
		log.Println("Executing createTableNotices failed:", err)
	}
	return
}

func Drop() (err error) {
	_, err = pool.Exec(context.Background(), dropTableNoticesQuery)
	if err != nil {
		log.Println("Executing dropTableNotices failed:", err)
	}
	return
}

// Close connection before exit
func Close() {
	log.Println("db:Close closing db connection")
	pool.Close()
}

func IsNotFound(err error) bool {
	return pgx.ErrNoRows == err
}
