package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateGroup(t *testing.T) {
	user := createRandomUser(t)

	group, err := testQueries.CreateGroup(context.Background(), user.Discordid)
	require.NoError(t, err)
	require.NotEmpty(t, group)
}

func TestGetGroupByHost(t *testing.T) {
	user := createRandomUser(t)

	group, err := testQueries.CreateGroup(context.Background(), user.Discordid)
	require.NoError(t, err)
	require.NotEmpty(t, group)

	group2, err := testQueries.GetGroupByHost(context.Background(), user.Discordid)
	require.NoError(t, err)
	require.Equal(t, group, group2)
}
