package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createNewAccount(t)
	account2 := createNewAccount(t)
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	n := 5

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	// check results
	exited := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.NotZero(t, transfer.ID)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)

		// check if the transfer is exists
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check if the entry is exists
		entry1 := result.FromEntry
		require.NotEmpty(t, entry1)
		require.NotZero(t, entry1.ID)
		require.Equal(t, account1.ID, entry1.AccountID)
		require.Equal(t, -amount, entry1.Amount)
		require.NotZero(t, entry1.CreatedAt)

		_, err = store.GetEntry(context.Background(), entry1.ID)
		require.NoError(t, err)

		entry2 := result.ToEntry
		require.NotEmpty(t, entry2)
		require.NotZero(t, entry2.ID)
		require.Equal(t, account2.ID, entry2.AccountID)
		require.Equal(t, amount, entry2.Amount)
		require.NotZero(t, entry2.CreatedAt)

		_, err = store.GetEntry(context.Background(), entry2.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check accounts balances

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)

		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, exited, k)
		exited[k] = true
	}

	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance-int64(n)*amount, updateAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updateAccount2.Balance)

}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createNewAccount(t)
	account2 := createNewAccount(t)
	amount := int64(10)

	errs := make(chan error)

	n := 10

	for i := 0; i < n; i++ {
		fromAccount := account1.ID
		toAccount := account2.ID
		if i%2 == 1 {
			fromAccount = account2.ID
			toAccount = account1.ID
		}
		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccount,
				ToAccountID:   toAccount,
				Amount:        amount,
			})

			errs <- err

		}()
	}

	// check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance, updateAccount1.Balance)
	require.Equal(t, account2.Balance, updateAccount2.Balance)

}
