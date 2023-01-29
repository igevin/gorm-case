package basic

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateDb(t *testing.T) {
	db, err := CreateDb()
	defer db.Close()
	require.NoError(t, err)
}
