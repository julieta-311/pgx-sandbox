package main

import (
	"context"
	"math"
	"os"
	"testing"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestReadThingByID(t *testing.T) {
	connString := os.Getenv("POSTGRES_URL")
	require.NotEmpty(t, connString, "POSTGRES_URL is required")

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connString)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close(ctx)) })

	db := &db{conn: conn}

	schemaMigrations, err := os.ReadFile("./testdata/initial_schema.sql")
	require.NoError(t, err)

	m := squirrel.Expr(string(schemaMigrations))
	_, err = squirrel.ExecContextWith(ctx, db, m)
	require.NoError(t, err)

	th := thing{
		ID:        id(uuid.New().String()),
		Name:      "Foo",
		Labels:    []string{"cat", "dog"},
		N:         7,
		X:         1.283,
		CreatedAt: time.Now().Truncate(time.Millisecond),
		Stuff:     map[string]any{"a": math.Pi, "b": "yes"},
	}
	require.NoError(t, db.insertThing(ctx, th))

	gotThing, err := db.readThingByID(ctx, th.ID)
	require.NoError(t, err)
	require.Equal(t, th, gotThing)
}
