package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	connString := os.Getenv("POSTGRES_URL")
	if connString == "" {
		fmt.Println("POSTGRES_URL is required")
		return
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		fmt.Printf("failed to connect to db: %v.\n", err)
		return
	}
	defer func() {
		if err := conn.Close(ctx); err != nil {
			log.Printf("failed to close db connection: %v.\n", err)
		}
	}()

	q, err := os.ReadFile("./testdata/initial_schema.sql")
	if err != nil {
		fmt.Printf("failed to read migrations file: %v.\n", err)
		return
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		fmt.Printf("failed to begin transaction: %v.\n", err)
		return
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, string(q))
	if err != nil {
		fmt.Printf("failed to execute initial db migration: %v.\n", err)
		return
	}

	if err := tx.Commit(ctx); err != nil {
		fmt.Printf("failed to commit initial db migration: %v.\n", err)
		return
	}

	fmt.Println("Success")
}
