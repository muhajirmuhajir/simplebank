package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createNewEntry(t *testing.T) Entry {
	account := createNewAccount(t)

	amount := 10000

	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    int64(amount),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	createNewEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry := createNewEntry(t)

	entry1, err := testQueries.GetEntry(context.Background(), entry.ID)

	require.NoError(t, err)
	require.NotZero(t, entry1.CreatedAt)
	require.Equal(t, entry.ID, entry1.ID)
	require.Equal(t, entry.AccountID, entry1.AccountID)
	require.Equal(t, entry.Amount, entry1.Amount)
}

func TestListEntry(t *testing.T) {
	for i := 0; i < 5; i++ {
		createNewEntry(t)
	}

	entries, err := testQueries.GetEntries(context.Background())

	require.NoError(t, err)
	require.NotZero(t, len(entries))
}
