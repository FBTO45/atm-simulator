package utils

import (
	"database/sql"
	"fmt"
)

func HandleTransactionError(tx *sql.Tx, err error) error {
	if rbErr := tx.Rollback(); rbErr != nil {
		return fmt.Errorf("transaction error: %v, rollback error: %v", err, rbErr)
	}
	return err
}
