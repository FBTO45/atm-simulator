package services

import (
	"atm-simulator/internal/database"
	"atm-simulator/internal/models"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AccountService struct{}

func NewAccountService() *AccountService {
	return &AccountService{}
}

func (s *AccountService) CreateAccount(name, pin string) (*models.Account, error) {
	hashedPin, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	result, err := database.DB.Exec(
		"INSERT INTO accounts (name, pin, balance) VALUES (?, ?, 0)",
		name, string(hashedPin),
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return s.GetAccountByID(int(id))
}

func (s *AccountService) GetAccountByID(id int) (*models.Account, error) {
	var account models.Account
	err := database.DB.QueryRow(
		"SELECT id, name, pin, balance, created_at FROM accounts WHERE id = ?",
		id,
	).Scan(&account.ID, &account.Name, &account.PIN, &account.Balance, &account.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("account not found")
		}
		return nil, err
	}
	return &account, nil
}

func (s *AccountService) Authenticate(accountID int, pin string) (*models.Account, error) {
	account, err := s.GetAccountByID(accountID)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.PIN), []byte(pin))
	if err != nil {
		return nil, errors.New("invalid PIN")
	}

	return account, nil
}

func (s *AccountService) UpdateBalance(accountID int, amount float64) error {
	_, err := database.DB.Exec(
		"UPDATE accounts SET balance = balance + ? WHERE id = ?",
		amount, accountID,
	)
	return err
}

func (s *AccountService) GetAccountBalance(accountID int) (float64, error) {
	var balance float64
	err := database.DB.QueryRow(
		"SELECT balance FROM accounts WHERE id = ?",
		accountID,
	).Scan(&balance)
	return balance, err
}
