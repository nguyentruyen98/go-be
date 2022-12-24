package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions ti execute ad queries and transactions

type Store interface {
	Querier
	TransfersTx(ctx context.Context, arg TransferTxParam) (TransferTxResult, error)
}

// SQLStore provides all functions ti execute SQL queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

// * NewStore creates a new Store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

// * execTx executes a functions within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParam struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}
type TransferTxResult struct {
	Transfers   Transfers `json:"transfers"`
	FromAccount Accounts  `json:"from_account"`
	ToAccount   Accounts  `json:"to_account"`
	FromEntry   Entries   `json:"from_entry"`
	ToEntry     Entries   `json:"to_entry"`
}

var txKey = struct{}{}

// * TransfersTx performs a money transfer from one account to another account
// * It creates a transfer record, add account entries and update accounts balance within a single database transaction

func (store *SQLStore) TransfersTx(ctx context.Context, arg TransferTxParam) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		txName := ctx.Value(txKey)

		fmt.Println(txName, "create transfer")
		var err error
		result.Transfers, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}
		// TODO: update account's balance
		// get account --> update it balance

		result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.FromAccountID,
			Amount: -arg.Amount,
		})

		if err != nil {
			return err
		}

		result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.ToAccountID,
			Amount: arg.Amount})

		if err != nil {
			return err
		}

		return nil
	})
	return result, err
}
