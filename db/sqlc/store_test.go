package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTranferTx(t *testing.T) {
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// * Run n concurrent transactions
	errs := make(chan error)
	results := make(chan TransferTxResult)
	n := 5
	amount := int64(10)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransfersTx(context.Background(), TransferTxParam{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result

		}()
	}

	// * Check result
	for i := 0; i < n; i++ {
		err := <-errs

		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// * Check transfer
		transfer := result.Transfers
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)

		require.NoError(t, err)

		// * Check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

	}
}
