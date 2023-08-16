package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/broemp/broempSignal/util"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:   util.RandomUser(),
		Discordid:  util.RandomId(),
		Telegramid: sql.NullInt64{Int64: util.RandomId(), Valid: true},
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Discordid, user.Discordid)
	require.Equal(t, arg.Telegramid, user.Telegramid)

	require.NotZero(t, user.Userid)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Userid)

	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Discordid, user2.Discordid)
	require.WithinDuration(t, user1.CreatedAt.Time, user2.CreatedAt.Time, time.Second)
}

func TestGetUserByDiscordId(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUserByDiscordId(context.Background(), user1.Discordid)

	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1, user2)
}

func TestGetUserByTelegramId(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUserByTelegramId(context.Background(), user1.Telegramid)

	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1, user2)
}

func TestDeleteUser(t *testing.T) {
	user1 := createRandomUser(t)
	err := testQueries.DeleteUser(context.Background(), user1.Userid)
	require.NoError(t, err)

	user2, err := testQueries.GetUser(context.Background(), user1.Userid)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user2)
}

func TestUpdateTelegramId(t *testing.T) {
	user1 := createRandomUser(t)
	arg := UpdateTelegramIdParams{
		Userid:     user1.Userid,
		Telegramid: sql.NullInt64{Int64: util.RandomId(), Valid: true},
	}

	err := testQueries.UpdateTelegramId(context.Background(), arg)
	require.NoError(t, err)
	user2, err := testQueries.GetUser(context.Background(), user1.Userid)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, arg.Telegramid, user2.Telegramid)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Discordid, user2.Discordid)
	require.WithinDuration(t, user1.CreatedAt.Time, user2.CreatedAt.Time, time.Second)
}
