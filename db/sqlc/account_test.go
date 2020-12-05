package db

import (
	"context"
	"testing"

	"github.com/muhajirmuhajir/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createNewAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createNewAccount(t)
}

func TestGetAccount(t *testing.T) {
	account := createNewAccount(t)

	account1, err := testQueries.GetAccount(context.Background(), account.ID)

	require.NoError(t, err)

	require.Equal(t, account.ID, account1.ID)
	require.Equal(t, account.Owner, account1.Owner)
	require.Equal(t, account.Balance, account1.Balance)

}

func TestDeleteAccount(t *testing.T) {
	account := createNewAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account.ID)

	require.NoError(t, err)

}

func TestUpdateAccount(t *testing.T) {
	account := createNewAccount(t)

	arg := UpdateAccountParams{
		ID:      account.ID,
		Balance: account.Balance + 10000,
	}

	account1, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.Equal(t, account.ID, account1.ID)
	require.Equal(t, int64(10000), account1.Balance-account.Balance)
}

func TestGetListAccounts(t *testing.T) {
	for i := 0; i < 5; i++ {
		createNewAccount(t)
	}

	accounts, err := testQueries.ListAccounts(context.Background())
	require.NoError(t, err)
	require.NotZero(t, len(accounts))
}
