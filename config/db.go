package config

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
)

var DB *pgx.Conn

func ConnectDatabase() {
	var err error
	DB, err = pgx.Connect(context.Background(), "postgres:///mygodb?user=postgres&sslmode=disable")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	err = DB.Ping(context.Background())
	if err != nil {
		log.Fatalf("Database connection is not alive: %v", err)
	}

	fmt.Println("Connected to Database")
}