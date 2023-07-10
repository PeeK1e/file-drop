package db

import (
	"database/sql"
	"log"
)

func RunTransaction(q []byte, args ...any) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(string(q), args...)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func RunTransactionWithResult(q []byte, args ...any) (*sql.Rows, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	r, err := tx.Query(string(q), args...)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return r, nil
}

func QueryWithResult(q []byte, args ...any) (*sql.Rows, error) {
	r, err := db.Query(string(q), args...)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func CloseRows(r *sql.Rows) {
	if err := r.Close(); err != nil {
		log.Printf("WARN: Could not close rows cleanly %s", err)
	}
}
