package main

import (
	"OPS/module/API"
	"OPS/module/DatabaseConn"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	//"log"
)

var dbpool *pgxpool.Pool // Global database pool

func init() {
	var err error
	dbpool, err = DatabaseConn.ConnectDB()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Call table creation function once the DB connection is established
	DatabaseConn.CreateTable(dbpool)
	DatabaseConn.CreateOCD(dbpool)
	DatabaseConn.CreatePaymentTable(dbpool)
	fmt.Println("Hi")
}

func main() {
	fmt.Println("Hello World")

	API.StartServer()

}
