package user

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func Insert(db *sqlx.DB, u User) error {
	_, err := db.NamedExec(`INSERT INTO users (name, email, balance) VALUES (:name, :email, :balance)`, u)
	return err
}

func GetAll(db *sqlx.DB) ([]User, error) {
	var users []User
	err := db.Select(&users, `SELECT * FROM users`)
	return users, err
}

func GetByID(db *sqlx.DB, id int) (User, error) {
	var u User
	err := db.Get(&u, `SELECT * FROM users WHERE id=$1`, id)
	return u, err
}

func TransferBalance(db *sqlx.DB, fromID, toID int, amount float64) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	var from, to User

	if err = tx.Get(&from, `SELECT * FROM users WHERE id=$1`, fromID); err != nil {
		tx.Rollback()
		return fmt.Errorf("sender not found")
	}
	if err = tx.Get(&to, `SELECT * FROM users WHERE id=$1`, toID); err != nil {
		tx.Rollback()
		return fmt.Errorf("receiver not found")
	}
	if from.Balance < amount {
		tx.Rollback()
		return fmt.Errorf("insufficient funds")
	}

	if _, err = tx.Exec(`UPDATE users SET balance = balance - $1 WHERE id=$2`, amount, fromID); err != nil {
		tx.Rollback()
		return err
	}
	if _, err = tx.Exec(`UPDATE users SET balance = balance + $1 WHERE id=$2`, amount, toID); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}