package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateAFK(t *testing.T) {
	user := createRandomUser(t)
	afk, err := testQueries.CreateAFK(context.Background(), sql.NullInt64{Int64: user.Discordid, Valid: true})

	require.NoError(t, err)
	require.NotEmpty(t, afk)
}

func TestGetAFK(t *testing.T) {
	user := createRandomUser(t)
	i := sql.NullInt64{Int64: 0, Valid: true}
	for i.Int64 <= 5 {
		afk, err := testQueries.CreateAFK(context.Background(), sql.NullInt64{Int64: user.Discordid, Valid: true})
		require.NoError(t, err)
		require.NotEmpty(t, afk)
		i.Int64++
	}
	num, err := testQueries.GetAFKCount(context.Background(), sql.NullInt64{Int64: user.Discordid, Valid: true})
	require.NoError(t, err)
	require.Equal(t, i.Int64, num)
}
