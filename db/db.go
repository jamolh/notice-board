package db

import (
	"context"
	"log"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"

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
func Connect() {
	var (
		err error
	)

	dbConnection, exists := os.LookupEnv("DATABASE_URL")
	if !exists {
		log.Fatal("db:Connect DATABASE_URL not found")
	}

	pool, err = pgxpool.Connect(context.Background(), dbConnection)
	if err != nil {
		log.Fatal("Connecting database failed:", err)
	}

	log.Println("------------DATABASE IS CONNECTED------------")
	_, err = pool.Exec(context.Background(), createTableNoticesQuery)
	if err != nil {
		log.Fatal("Executing createTableNotices failed:", err)
	}
}

// Close connection before exit
func Close() {
	log.Println("db:Close closing db connection")
	pool.Close()
}
