package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ZhigerDinmukhamed/go-practice4/internal/user"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	dsn := "user=user password=password dbname=mydatabase host=localhost port=5430 sslmode=disable"
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatalln(err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	fmt.Println("Connected")

	err = user.Insert(db, user.User{Name: "Zhiga", Email: "zhigae@gmail.com", Balance: 1000})
	if err != nil {
		log.Println(err)
	}
	fmt.Println(user.GetByID(db, 1))

	err = user.Insert(db, user.User{Name: "ivan", Email: "ivan@gmail.com", Balance: 1100})
	if err != nil {
		log.Println(err)
	}
	fmt.Println(user.GetByID(db, 2))

	err = user.TransferBalance(db, 1, 2, 10)
	if err != nil {
		log.Println(err)
	}

	users, _ := user.GetAll(db)
	fmt.Println(users)
}
