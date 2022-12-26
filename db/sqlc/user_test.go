package db

import (
	"context"
	"testing"
	"time"

	"github.com/nguyentruyen98/go-be/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) Users {
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: "myPassword",
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Username, user.Username)

	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())

	return user
}
func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.Username, user2.Username)

	require.WithinDuration(t, user1.PasswordChangedAt, user1.PasswordChangedAt, 1*time.Second)
	require.WithinDuration(t, user1.CreatedAt, user1.CreatedAt, 1*time.Second)

}