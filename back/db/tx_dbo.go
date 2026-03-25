package db

import (
	"database/sql"
	"fmt"
)

// TxDBO encapsule une transaction SQL, même API que DBO mais avec commit/rollback
type TxDBO struct {
	tx *sql.Tx
}

func (t *TxDBO) QueryRows(query string, args ...any) (*sql.Rows, error) {
	rows, err := t.tx.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("TxDBO QueryRows error: %w", err)
	}
	return rows, nil
}

func (t *TxDBO) QueryRow(query string, args ...any) *sql.Row {
	return t.tx.QueryRow(query, args...)
}

func (t *TxDBO) Exec(query string, args ...any) (int64, error) {
	result, err := t.tx.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("TxDBO Exec error: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("TxDBO RowsAffected error: %w", err)
	}
	return rowsAffected, nil
}

func (t *TxDBO) ExecReturning(query string, args ...any) *sql.Row {
	return t.tx.QueryRow(query, args...)
}

func (t *TxDBO) Commit() error {
	return t.tx.Commit()
}

func (t *TxDBO) Rollback() error {
	return t.tx.Rollback()
}