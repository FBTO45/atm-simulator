package services

import (
	"atm-simulator/internal/database"
	"atm-simulator/internal/models"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type TransactionService struct{}

func NewTransactionService() *TransactionService {
	return &TransactionService{}
}

func (s *TransactionService) CreateTransaction(tx *sql.Tx, transaction *models.Transaction) error {
	_, err := tx.Exec(
		"INSERT INTO transactions (account_id, type, amount, target_id, description) VALUES (?, ?, ?, ?, ?)",
		transaction.AccountID, transaction.Type, transaction.Amount, transaction.TargetID, transaction.Description,
	)
	return err
}

func (s *TransactionService) Deposit(accountID int, amount float64, description string) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Update account balance
	_, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, accountID)
	if err != nil {
		return err
	}

	// Create transaction record
	transaction := &models.Transaction{
		AccountID:   accountID,
		Type:        models.Deposit,
		Amount:      amount,
		Description: description,
		CreatedAt:   time.Now(),
	}
	if err := s.CreateTransaction(tx, transaction); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *TransactionService) Withdraw(accountID int, amount float64, description string) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Check balance
	var balance float64
	err = tx.QueryRow("SELECT balance FROM accounts WHERE id = ? FOR UPDATE", accountID).Scan(&balance)
	if err != nil {
		return err
	}

	if balance < amount {
		return errors.New("insufficient balance")
	}

	// Update account balance
	_, err = tx.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, accountID)
	if err != nil {
		return err
	}

	// Create transaction record
	transaction := &models.Transaction{
		AccountID:   accountID,
		Type:        models.Withdraw,
		Amount:      amount,
		Description: description,
		CreatedAt:   time.Now(),
	}
	if err := s.CreateTransaction(tx, transaction); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *TransactionService) Transfer(senderID, receiverID int, amount float64, description string) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	if senderID == receiverID {
		return errors.New("cannot transfer to yourself")
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Check if receiver exists
	var receiverExists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM accounts WHERE id = ?)", receiverID).Scan(&receiverExists)
	if err != nil {
		return err
	}
	if !receiverExists {
		return errors.New("receiver account not found")
	}

	// Check sender balance
	var senderBalance float64
	err = tx.QueryRow("SELECT balance FROM accounts WHERE id = ? FOR UPDATE", senderID).Scan(&senderBalance)
	if err != nil {
		return err
	}

	if senderBalance < amount {
		return errors.New("insufficient balance")
	}

	// Update sender balance
	_, err = tx.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, senderID)
	if err != nil {
		return err
	}

	// Update receiver balance
	_, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, receiverID)
	if err != nil {
		return err
	}

	// Create transfer_out transaction for sender
	senderTransaction := &models.Transaction{
		AccountID:   senderID,
		Type:        models.TransferOut,
		Amount:      amount,
		TargetID:    &receiverID,
		Description: fmt.Sprintf("Transfer to account %d: %s", receiverID, description),
		CreatedAt:   time.Now(),
	}
	if err := s.CreateTransaction(tx, senderTransaction); err != nil {
		return err
	}

	// Create transfer_in transaction for receiver
	receiverTransaction := &models.Transaction{
		AccountID:   receiverID,
		Type:        models.TransferIn,
		Amount:      amount,
		TargetID:    &senderID,
		Description: fmt.Sprintf("Transfer from account %d: %s", senderID, description),
		CreatedAt:   time.Now(),
	}
	if err := s.CreateTransaction(tx, receiverTransaction); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *TransactionService) GetTransactionHistory(accountID int) ([]models.Transaction, error) {
	rows, err := database.DB.Query(
		"SELECT id, account_id, type, amount, target_id, description, created_at FROM transactions WHERE account_id = ? ORDER BY created_at DESC LIMIT 10",
		accountID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		err := rows.Scan(&t.ID, &t.AccountID, &t.Type, &t.Amount, &t.TargetID, &t.Description, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}
