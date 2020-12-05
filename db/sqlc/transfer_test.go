package db

import (
	"context"
	"testing"

	"github.com/muhajirmuhajir/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createNewTransfer(t *testing.T) Transfer {
	account1 := createNewAccount(t)
	account2 := createNewAccount(t)
	amount := util.RandomMoney()

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        amount,
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotZero(t, transfer.ID)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestTransfer(t *testing.T) {
	createNewTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	transfer := createNewTransfer(t)

	transfer1, err := testQueries.GetTransfer(context.Background(), transfer.ID)

	require.NoError(t, err)
	require.Equal(t, transfer.ID, transfer1.ID)
	require.Equal(t, transfer.Amount, transfer1.Amount)
	require.Equal(t, transfer.FromAccountID, transfer1.FromAccountID)
	require.Equal(t, transfer.ToAccountID, transfer1.ToAccountID)
}

func TestGetListTransfers(t *testing.T) {
	for i := 0; i < 5; i++ {
		createNewTransfer(t)
	}

	transfers, err := testQueries.GetListTransfers(context.Background())
	require.NoError(t, err)
	require.NotZero(t, len(transfers))
}
