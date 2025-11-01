package main

import (
 "database/sql"
 "encoding/json"
 "fmt"
 "log"
 "net/http"
 "strconv"
 "strings"
 "time"

 _ "github.com/mattn/go-sqlite3"
)

type Book struct {
 ID     int    json:"id"
 Title  string json:"title"
 Author string json:"author"
 Genre  string json:"genre"
 Price  int    json:"price"
}

func main() {
 db, err := sql.Open("sqlite3", "./books.db")
 if err != nil {
  log.Fatalf("failed to open db: %v", err)
 }
 defer db.Close()

 createTable := `
 CREATE TABLE IF NOT EXISTS books (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  title TEXT NOT NULL,
  author TEXT NOT NULL,
  genre TEXT NOT NULL,
  price INTEGER NOT NULL
 );
 `
 if _, err := db.Exec(createTable); err != nil {
  log.Fatalf("failed to create table: %v", err)
 }

 var count int
 if err := db.QueryRow("SELECT COUNT(1) FROM books").Scan(&count); err != nil {
  log.Fatalf("count error: %v", err)
 }
 if count == 0 {
  _, err := db.Exec(`INSERT INTO books (title, author, genre, price) VALUES
   ('War and Peace', 'Leo Tolstoy', 'fiction', 500),
   ('Crime and Punishment', 'Fyodor Dostoevsky', 'fiction', 450),
   ('Go Programming', 'Alan A. A. Donovan', 'education', 300),
   ('Data Structures', 'N. Wirth', 'education', 350),
   ('The Hobbit', 'J.R.R. Tolkien', 'fantasy', 400)
  `)
  if err != nil {
   log.Fatalf("insert sample data error: %v", err)
  }
 }

 http.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
  q := r.URL.Query()
  genre := q.Get("genre")
  sortParam := q.Get("sort")
  limitParam := q.Get("limit")
  offsetParam := q.Get("offset")

  sqlBase := "SELECT id, title, author, genre, price FROM books"
  var whereParts []string
  var args []interface{}

  if genre != "" {
   whereParts = append(whereParts, "genre = ?")
   args = append(args, genre)
  }

  if len(whereParts) > 0 {
   sqlBase += " WHERE " + strings.Join(whereParts, " AND ")
  }

  if sortParam == "price_asc" {
   sqlBase += " ORDER BY price ASC"
  } else if sortParam == "price_desc" {
   sqlBase += " ORDER BY price DESC"
  }

  if limitParam != "" {
   limit, err := strconv.Atoi(limitParam)
   if err != nil || limit <= 0 {
    http.Error(w, "invalid limit (must be positive integer)", http.StatusBadRequest)
    return
   }
   
   if limit > 1000 {
    limit = 1000
   }
   sqlBase += " LIMIT ?"
   args = append(args, limit)
  }
  if offsetParam != "" {
   offset, err := strconv.Atoi(offsetParam)
   if err != nil || offset < 0 {
    http.Error(w, "invalid offset (must be non-negative integer)", http.StatusBadRequest)
    return
   }
   sqlBase += " OFFSET ?"
   args = append(args, offset)
  }

  startQuery := time.Now()
  rows, err := db.Query(sqlBase, args...)
  queryDuration := time.Since(startQuery)

  log.Printf("SQL: %s -- args=%v -- took=%v", sqlBase, args, queryDuration)

  if err != nil {
   http.Error(w, "internal db error", http.StatusInternalServerError)
   log.Printf("query error: %v", err)
   return
  }
  defer rows.Close()

  var books []Book
  for rows.Next() {
   var b Book
   if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Genre, &b.Price); err != nil {
    http.Error(w, "error scanning rows", http.StatusInternalServerError)
    log.Printf("scan error: %v", err)
    return
   }
   books = append(books, b)
  }
  if err := rows.Err(); err != nil {
   http.Error(w, "rows error", http.StatusInternalServerError)
   log.Printf("rows iteration error: %v", err)
   return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("X-Query-Time", fmt.Sprintf("%v", queryDuration))
  w.WriteHeader(http.StatusOK)

  if err := json.NewEncoder(w).Encode(books); err != nil {
   log.Printf("encode error: %v", err)
  }
 })

 fmt.Println("Server started at :8080")
 log.Fatal(http.ListenAndServe(":8080", nil))
}
