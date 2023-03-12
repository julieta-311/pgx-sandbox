package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/Masterminds/squirrel"
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
		fmt.Printf("Failed to connect to db: %v.\n", err)
		return
	}
	defer func() {
		if err := conn.Close(ctx); err != nil {
			fmt.Printf("Failed to close db connection: %v.\n", err)
		}
	}()

	db := &db{conn: conn}

	schemaMigrations, err := os.ReadFile("./testdata/initial_schema.sql")
	if err != nil {
		fmt.Printf("Failed to read migrations file: %v.\n", err)
		return
	}

	m := squirrel.Expr(string(schemaMigrations))
	if _, err := squirrel.ExecContextWith(ctx, db, m); err != nil {
		fmt.Printf("Failed to run initial db migration: %v.\n", err)
		return
	}

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	insert := psql.Insert("thing").
		Columns("name", "labels", "n", "x", "created_at", "stuff").
		Values(
			"Foo",
			[]string{"cat", "dog"},
			7,
			1.283,
			time.Now(),
			`{"a": 7, "b": "yes"}`)

	if _, err := squirrel.ExecContextWith(ctx, db, insert); err != nil {
		fmt.Printf("Failed to execute insertion: %v.\n", err)
		return
	}

	fmt.Println("Success")
}

type db struct {
	conn *pgx.Conn
}

func (d *db) ExecContext(ctx context.Context, query string, args ...any) (_ sql.Result, err error) {
	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q with args %q: %w", query, args, err)
	}

	return nil, tx.Commit(ctx)
}
