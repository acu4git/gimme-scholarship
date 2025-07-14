package database

import (
	"testing"

	"github.com/acu4git/gimme-scholarship/internal/domain/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindScholarships(t *testing.T) {
	fdb := NewFakeDatabase()
	t.Cleanup(fdb.TruncateTables)
	require.NoError(t, fdb.TestInitScholarships())

	t.Run("正常系", func(t *testing.T) {
		cases := []struct {
			name    string
			option  repository.FilterOption
			wantLen int
		}{
			{
				name:    "指定なし",
				option:  repository.FilterOption{},
				wantLen: 5,
			}, {
				name:    "指定あり(Target)",
				option:  repository.FilterOption{Target: "学部生"},
				wantLen: 3,
			}, {
				name:    "指定あり(Type)",
				option:  repository.FilterOption{Type: "給付"},
				wantLen: 4,
			},
		}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				scholarships, _, err := fdb.FindScholarships(tt.option)
				require.NoError(t, err)
				assert.Equal(t, tt.wantLen, len(scholarships))
			})
		}
	})
}
