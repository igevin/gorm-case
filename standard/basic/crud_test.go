package basic

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateDb(t *testing.T) {
	_ = createDbForTest(t)
}

func TestCreateTable(t *testing.T) {
	db := createDbForTest(t)
	defer db.Close()
	createTableForTest(t, db)
}

func createDbForTest(t *testing.T) *sql.DB {
	db, err := CreateDb()
	require.NoError(t, err)
	return db
}

func createTableForTest(t *testing.T, db *sql.DB) {
	affected, err := CreateTable(db, context.Background())
	require.NoError(t, err)
	assert.Equal(t, int64(0), affected)
}
