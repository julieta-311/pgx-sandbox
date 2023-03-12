package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Success!")
}

func run() error {
	connString := os.Getenv("POSTGRES_URL")
	if connString == "" {
		return errors.New("POSTGRES_URL is required")
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return fmt.Errorf("failed to connect to db: %w", err)
	}
	defer func() {
		if err := conn.Close(ctx); err != nil {
			fmt.Printf("Failed to close db connection: %v.\n", err)
		}
	}()

	db := &db{conn: conn}

	schemaMigrations, err := os.ReadFile("./testdata/initial_schema.sql")
	if err != nil {
		return fmt.Errorf("failed to read migrations file: %w", err)
	}

	m := squirrel.Expr(string(schemaMigrations))
	if _, err := squirrel.ExecContextWith(ctx, db, m); err != nil {
		return fmt.Errorf("failed to run initial db migration: %w", err)
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
		return fmt.Errorf("failed to execute insertion: %w", err)
	}

	return nil
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
