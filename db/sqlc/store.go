package simplebank

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB

}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error{
	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil{
		return err
	}

	q := New(tx)
	err =fn(q)
	if err != nil{
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId int64 `json:"to_account_id"`
	Amount int64 `json:"amount"`
}


type TransferTxResult struct {

	Transfer Transfer `json:"transfer"`
	FromAccount Account `json:"from_account"`
	ToAccount Account `json:"to_account"`
	FromEntry Entry `json:"from_entry"`
	ToEntry Entry `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error){
	var result TransferTxResult

	err :=store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID: arg.ToAccountId,
			Amount: arg.Amount,
		})

		if err != nil {
			return err
		}
		
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount: -arg.Amount,	
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountId,
			Amount: arg.Amount,	
		})

		if err != nil {
			return err
		}

		// update account balance
		account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountId)
		if err != nil {
			return err
		}

		result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID: arg.FromAccountId,
			Balance: account1.Balance-arg.Amount,
		})
		
		if err != nil {
			return err
		}

		account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountId)
		if err != nil {
			return err
		}


		result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID: arg.ToAccountId,
			Balance: account2.Balance+arg.Amount,
		})
		
		if err != nil {
			return err
		}

		return nil
	})
	return result, err
}