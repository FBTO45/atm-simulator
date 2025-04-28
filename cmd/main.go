package main

import (
	"atm-simulator/internal/database"
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID        int
	Name      string
	PIN       string
	Balance   float64
	CreatedAt time.Time
}

type Transaction struct {
	ID          int
	AccountID   int
	Type        string
	Amount      float64
	TargetID    sql.NullInt64
	Description sql.NullString
	CreatedAt   time.Time
}

var db *sql.DB
var currentAccount *Account

func main() {
	// Initialize database
	database.InitDB()
	db = database.DB
	defer db.Close()

	for {
		fmt.Println("\nATM Simulator CLI")
		fmt.Println("1. Register")
		fmt.Println("2. Login")
		fmt.Println("3. Check Balance")
		fmt.Println("4. Deposit")
		fmt.Println("5. Withdraw")
		fmt.Println("6. Transfer")
		fmt.Println("7. Transaction History")
		fmt.Println("8. Logout")
		fmt.Println("9. Exit")
		fmt.Print("Select option: ")

		reader := bufio.NewReader(os.Stdin)
		option, _ := reader.ReadString('\n')
		option = strings.TrimSpace(option)

		switch option {
		case "1":
			register()
		case "2":
			login()
		case "3":
			checkBalance()
		case "4":
			deposit()
		case "5":
			withdraw()
		case "6":
			transfer()
		case "7":
			transactionHistory()
		case "8":
			logout()
		case "9":
			fmt.Println("Thank you for using our ATM service!")
			return
		default:
			fmt.Println("Invalid option, please try again")
		}
	}
}

func register() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Enter your PIN (4 digits): ")
	pin, _ := reader.ReadString('\n')
	pin = strings.TrimSpace(pin)

	if len(pin) != 4 {
		fmt.Println("PIN must be 4 digits")
		return
	}

	hashedPin, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error creating account:", err)
		return
	}

	result, err := db.Exec("INSERT INTO accounts (name, pin) VALUES (?, ?)", name, string(hashedPin))
	if err != nil {
		fmt.Println("Error creating account:", err)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Error creating account:", err)
		return
	}

	fmt.Printf("Account created successfully!\nYour account number is: %d\nPlease remember this number and your PIN\n", id)
}

func login() {
	if currentAccount != nil {
		fmt.Println("You are already logged in as", currentAccount.Name)
		return
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your account number: ")
	accountIDStr, _ := reader.ReadString('\n')
	accountIDStr = strings.TrimSpace(accountIDStr)
	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		fmt.Println("Invalid account number")
		return
	}

	fmt.Print("Enter your PIN: ")
	pin, _ := reader.ReadString('\n')
	pin = strings.TrimSpace(pin)

	var account Account
	err = db.QueryRow("SELECT id, name, pin, balance, created_at FROM accounts WHERE id = ?", accountID).
		Scan(&account.ID, &account.Name, &account.PIN, &account.Balance, &account.CreatedAt)
	if err != nil {
		fmt.Println("Authentication failed: account not found")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.PIN), []byte(pin))
	if err != nil {
		fmt.Println("Authentication failed: invalid PIN")
		return
	}

	currentAccount = &account
	fmt.Printf("Welcome, %s!\n", account.Name)
}

func checkBalance() {
	if currentAccount == nil {
		fmt.Println("Please login first")
		return
	}

	fmt.Printf("Your current balance is: $%.2f\n", currentAccount.Balance)
}

func deposit() {
	if currentAccount == nil {
		fmt.Println("Please login first")
		return
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter amount to deposit: ")
	amountStr, _ := reader.ReadString('\n')
	amountStr = strings.TrimSpace(amountStr)
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		fmt.Println("Invalid amount")
		return
	}

	fmt.Print("Enter description (optional): ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	tx, err := db.Begin()
	if err != nil {
		fmt.Println("Error starting transaction:", err)
		return
	}

	_, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, currentAccount.ID)
	if err != nil {
		tx.Rollback()
		fmt.Println("Error depositing money:", err)
		return
	}

	var desc interface{}
	if description == "" {
		desc = nil
	} else {
		desc = description
	}

	_, err = tx.Exec("INSERT INTO transactions (account_id, type, amount, description) VALUES (?, 'deposit', ?, ?)",
		currentAccount.ID, amount, desc)
	if err != nil {
		tx.Rollback()
		fmt.Println("Error recording transaction:", err)
		return
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println("Error completing transaction:", err)
		return
	}

	// Update current account balance in memory
	currentAccount.Balance += amount
	fmt.Printf("Successfully deposited $%.2f to your account\n", amount)
}

func withdraw() {
	if currentAccount == nil {
		fmt.Println("Please login first")
		return
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter amount to withdraw: ")
	amountStr, _ := reader.ReadString('\n')
	amountStr = strings.TrimSpace(amountStr)
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		fmt.Println("Invalid amount")
		return
	}

	if amount > currentAccount.Balance {
		fmt.Println("Insufficient balance")
		return
	}

	fmt.Print("Enter description (optional): ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	tx, err := db.Begin()
	if err != nil {
		fmt.Println("Error starting transaction:", err)
		return
	}

	_, err = tx.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, currentAccount.ID)
	if err != nil {
		tx.Rollback()
		fmt.Println("Error withdrawing money:", err)
		return
	}

	var desc interface{}
	if description == "" {
		desc = nil
	} else {
		desc = description
	}

	_, err = tx.Exec("INSERT INTO transactions (account_id, type, amount, description) VALUES (?, 'withdraw', ?, ?)",
		currentAccount.ID, amount, desc)
	if err != nil {
		tx.Rollback()
		fmt.Println("Error recording transaction:", err)
		return
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println("Error completing transaction:", err)
		return
	}

	// Update current account balance in memory
	currentAccount.Balance -= amount
	fmt.Printf("Successfully withdrew $%.2f from your account\n", amount)
}

func transfer() {
	if currentAccount == nil {
		fmt.Println("Please login first")
		return
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter recipient account number: ")
	recipientIDStr, _ := reader.ReadString('\n')
	recipientIDStr = strings.TrimSpace(recipientIDStr)
	recipientID, err := strconv.Atoi(recipientIDStr)
	if err != nil {
		fmt.Println("Invalid account number")
		return
	}

	if recipientID == currentAccount.ID {
		fmt.Println("Cannot transfer to yourself")
		return
	}

	fmt.Print("Enter amount to transfer: ")
	amountStr, _ := reader.ReadString('\n')
	amountStr = strings.TrimSpace(amountStr)
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		fmt.Println("Invalid amount")
		return
	}

	if amount > currentAccount.Balance {
		fmt.Println("Insufficient balance")
		return
	}

	fmt.Print("Enter description (optional): ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	// Check if recipient exists
	var recipientExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM accounts WHERE id = ?)", recipientID).Scan(&recipientExists)
	if err != nil {
		fmt.Println("Error checking recipient account:", err)
		return
	}
	if !recipientExists {
		fmt.Println("Recipient account not found")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		fmt.Println("Error starting transaction:", err)
		return
	}

	// Deduct from sender
	_, err = tx.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, currentAccount.ID)
	if err != nil {
		tx.Rollback()
		fmt.Println("Error transferring money:", err)
		return
	}

	// Add to recipient
	_, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, recipientID)
	if err != nil {
		tx.Rollback()
		fmt.Println("Error transferring money:", err)
		return
	}

	// Record sender transaction (transfer_out)
	var desc interface{}
	if description == "" {
		desc = nil
	} else {
		desc = description
	}

	_, err = tx.Exec("INSERT INTO transactions (account_id, type, amount, target_id, description) VALUES (?, 'transfer_out', ?, ?, ?)",
		currentAccount.ID, amount, recipientID, desc)
	if err != nil {
		tx.Rollback()
		fmt.Println("Error recording transaction:", err)
		return
	}

	// Record recipient transaction (transfer_in)
	_, err = tx.Exec("INSERT INTO transactions (account_id, type, amount, target_id, description) VALUES (?, 'transfer_in', ?, ?, ?)",
		recipientID, amount, currentAccount.ID, desc)
	if err != nil {
		tx.Rollback()
		fmt.Println("Error recording transaction:", err)
		return
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println("Error completing transaction:", err)
		return
	}

	// Update current account balance in memory
	currentAccount.Balance -= amount
	fmt.Printf("Successfully transferred $%.2f to account %d\n", amount, recipientID)
}

func transactionHistory() {
	if currentAccount == nil {
		fmt.Println("Please login first")
		return
	}

	rows, err := db.Query(`
		SELECT id, type, amount, target_id, description, created_at 
		FROM transactions 
		WHERE account_id = ? 
		ORDER BY created_at DESC 
		LIMIT 10`, currentAccount.ID)
	if err != nil {
		fmt.Println("Error retrieving transaction history:", err)
		return
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		err := rows.Scan(&t.ID, &t.Type, &t.Amount, &t.TargetID, &t.Description, &t.CreatedAt)
		if err != nil {
			fmt.Println("Error reading transaction:", err)
			return
		}
		transactions = append(transactions, t)
	}

	if len(transactions) == 0 {
		fmt.Println("No transactions found")
		return
	}

	fmt.Println("\nYour recent transactions:")
	fmt.Println("----------------------------------------")
	for _, t := range transactions {
		fmt.Printf("Date: %s\n", t.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Type: %s\n", t.Type)
		fmt.Printf("Amount: $%.2f\n", t.Amount)
		if t.TargetID.Valid {
			fmt.Printf("Target Account: %d\n", t.TargetID.Int64)
		}
		if t.Description.Valid {
			fmt.Printf("Description: %s\n", t.Description.String)
		}
		fmt.Println("----------------------------------------")
	}
}

func logout() {
	if currentAccount == nil {
		fmt.Println("You are not logged in")
		return
	}

	fmt.Printf("Goodbye, %s!\n", currentAccount.Name)
	currentAccount = nil
}
