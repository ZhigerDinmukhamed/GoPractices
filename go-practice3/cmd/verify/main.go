package verify

import (
	"database/sql"
	"log"
)

func main() {
	dsn := "postgres://postgres:140205091205@localhost:5432/expense_tracker?sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	
	if err != nil{
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
}
