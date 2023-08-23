package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func main() {
	fmt.Println("Nothing to see, only tests to run move it along!")
}

type db struct {
	conn *pgx.Conn
}

// ExecContext implements the squirrel ExecerContext interface, that is it's used whenever
// sq.ExecContext or sq.ExecContextWith are invoked.
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
