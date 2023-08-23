package main

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type id string

func (i id) Valid() bool {
	_, err := uuid.Parse(string(i))
	return err == nil
}

// Value implements the driver.Valuer interface.
func (i id) Value() string {
	return string(i) + "::uuid"
}

type thing struct {
	ID        id             `db:"thing_id"`
	Name      string         `db:"name"`
	Labels    []string       `db:"labels"`
	N         int            `db:"n"`
	X         float64        `db:"x"`
	CreatedAt time.Time      `db:"created_at"`
	Stuff     map[string]any `db:"stuff"`
}

func (d *db) insertThing(ctx context.Context, t thing) error {
	cols := []string{"name", "labels", "n", "x", "stuff"}
	vals := []any{t.Name, t.Labels, t.N, t.X, t.Stuff}

	if !t.CreatedAt.IsZero() {
		cols = append(cols, "created_at")
		vals = append(vals, t.CreatedAt)
	}

	if t.ID != "" {
		cols = append(cols, "thing_id")
		vals = append(vals, t.ID)
	}

	insertion := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("thing").Columns(cols...).Values(vals...)

	if _, err := sq.ExecContextWith(ctx, d, insertion); err != nil {
		return fmt.Errorf("executing insertion: %w", err)
	}

	return nil
}

// readThingByID reads a thing identified by the given id. Since I want to use
// pgx's methods to scan the row, call squirrel's 'ToSql' to build the statement
// and execute it so it returns a collectible row.
func (d *db) readThingByID(ctx context.Context, id id) (th thing, err error) {
	q, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("thing_id", "name", "labels", "n", "x", "created_at", "stuff").
		From("thing").
		Where(sq.Eq{"thing_id": id}).
		ToSql()

	if err != nil {
		return th, fmt.Errorf("building statement: %w", err)
	}

	rows, err := d.conn.Query(ctx, q, args...)
	if err != nil {
		return th, fmt.Errorf("querying for thing %q: %w", id, err)
	}

	th, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[thing])
	if err != nil {
		return th, fmt.Errorf("scanning thing %q: %w", id, err)
	}

	return th, nil
}
