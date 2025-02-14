package DatabaseConn

import (
	"OPS/module/Orders"
	ord "OPS/module/Orders"
	"context"
	"fmt"
	"log"

	//"error"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"
)

const dbURL = "postgres://user:password@localhost:5432/orders?sslmode=disable"

func ConnectDB() (dbpool *pgxpool.Pool, err1 error) {

	dbpool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %s\n", err)
		return nil, fmt.Errorf("Unable to connect to DB")
	}
	//defer dbpool.Close()

	fmt.Println("Connected to PostgreSQL successfully!")

	fmt.Println("Successfully connected to DB")

	//getTableSchema(dbpool, "order")
	createTable(dbpool)
	return dbpool, nil
}

func AddRow(dbpool *pgxpool.Pool, od ord.Order) (int, error) {

	insertQuery := `INSERT INTO "order" (user_id,product_id,quantity,total_price) VALUES ($1, $2, $3,$4) RETURNING user_id`
	var id int
	err := dbpool.QueryRow(context.Background(), insertQuery, od.User_Id, od.Product_Id, od.Quantity, od.Total_Price).Scan(&id)
	if err != nil {
		//log.Fatalf("Failed to insert row: %v\n", err)
		return 0, err
	}

	fmt.Printf("Inserted new order with ID: %d\n", id)
	return id, nil
}

func GetAll(dbpool *pgxpool.Pool) ([]Orders.Order, error) {
	allrow := `SELECT * from "order"`

	rows, err := dbpool.Query(context.Background(), allrow)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders: %w", err)
	}
	defer rows.Close() // Close the rows iterator

	var orders []Orders.Order

	// Iterate over rows and scan into Order struct
	for rows.Next() {
		var o Orders.Order
		err := rows.Scan(&o.User_Id, &o.Product_Id, &o.Quantity, &o.Total_Price)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		orders = append(orders, o)
	}

	// Check for iteration errors
	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating rows: %w", rows.Err())
	}

	return orders, nil
}

func createTable(dbpool *pgxpool.Pool) {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS "order" (
		user_id SERIAL PRIMARY KEY,
		product_id INT NOT NULL,
		quantity INT NOT NULL,
		total_price FLOAT NOT NULL
	);
	`
	var err error
	_, err = dbpool.Exec(context.Background(), createTableQuery)
	if err != nil {
		log.Fatalf("Failed to create table: %v\n", err)
	}

	fmt.Println("Table 'order' created successfully!")
}

func getTableSchema(dbpool *pgxpool.Pool, tableName string) {
	query := `
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' AND table_type = 'BASE TABLE';
	`

	rows, err := dbpool.Query(context.Background(), query)
	if err != nil {
		log.Fatalf("Error fetching tables: %v\n", err)
	}
	defer rows.Close()
	if !rows.Next() {
		fmt.Println("No tables found in the 'orders' database.")
		return
	}
	fmt.Println("Tables in the 'orders' database:")
	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}
		fmt.Println(tableName)
	}
}
