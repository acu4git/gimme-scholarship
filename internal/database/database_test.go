package database

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_FindScholarships(t *testing.T) {
	fdb := NewFakeDatabase()
	t.Cleanup(fdb.TruncateTables)
	require.NoError(t, fdb.TestInitScholarships())
}
