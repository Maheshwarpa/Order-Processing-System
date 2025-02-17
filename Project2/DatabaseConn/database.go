package DatabaseConn

import (
	//"OPS/module/Orders"
	"OPS/module/Event"
	ord "OPS/module/Orders"
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	//"error"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"
)

const dbURL = "postgres://user:password@localhost:5432/orders?sslmode=disable"

var DbPool *pgxpool.Pool

func ConnectDB() (*pgxpool.Pool, error) {

	Dbp, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %s\n", err)
		return nil, fmt.Errorf("Unable to connect to DB")
	}
	//defer dbpool.Close()

	fmt.Println("Connected to PostgreSQL successfully!")

	fmt.Println("Successfully connected to DB")

	//getTableSchema(dbpool, "order")
	//createTable(dbpool)
	DbPool = Dbp
	return Dbp, nil
}

func AddRow(Dbpool *pgxpool.Pool, od ord.Order) (int, error) {

	insertQuery := `INSERT INTO "order" (user_id,product_id,quantity,total_price) VALUES ($1, $2, $3,$4) RETURNING user_id`
	var id int
	err := Dbpool.QueryRow(context.Background(), insertQuery, od.User_Id, od.Product_Id, od.Quantity, od.Total_Price).Scan(&id)
	if err != nil {
		log.Fatalf("Failed to insert row: %v\n", err)
		return 0, err
	}

	fmt.Printf("Inserted new order with ID: %d\n", id)
	return id, nil
}

func AddOrder(Dbpool *pgxpool.Pool, od ord.Order) (*Event.OrderCreatedResponse, error) {

	insertQuery := `INSERT INTO "ordercreated" (order_id,status) VALUES ($1, $2) RETURNING order_id`
	var id string
	fmt.Println(od.User_Id)
	var str string = "order_"
	str += strconv.Itoa(od.Product_Id)
	str += "_"
	str += strconv.Itoa(od.User_Id)
	var status string = "PENDING"
	fmt.Println(str)
	err := Dbpool.QueryRow(context.Background(), insertQuery, str, status).Scan(&id)
	if err != nil {
		log.Fatalf("Failed to insert row: %v\n", err)
		return &Event.OrderCreatedResponse{}, err
	}

	fmt.Printf("Inserted new order with Order-ID: %v\n", str)

	ord1 := Event.OrderCreatedResponse{str, status}
	Event.OCRList = append(Event.OCRList, ord1)
	fmt.Printf("OrderEvent Response has been created: %+v\n", ord1)
	return &ord1, nil
}

func AddPaymentStatus(Dbpool *pgxpool.Pool, od Event.ProcessingResponse) error {

	insertQuery := `INSERT INTO "paymentstatus" (order_id,payment_status) VALUES ($1, $2) RETURNING order_id`
	var id string

	//var status string = "PENDING"
	err := Dbpool.QueryRow(context.Background(), insertQuery, od.Order_id, od.Payment_status).Scan(&id)
	if err != nil {
		log.Fatalf("Failed to insert row: %v\n", err)
		return err
	}

	fmt.Printf("Inserted new payment status for Order-ID: %v\n", od.Order_id)
	return nil
}

func GetPrice(Dbpool *pgxpool.Pool, orderID string) (float64, error) {
	query := `SELECT total_price FROM "order" WHERE user_id = $1`
	var price float64
	//Converting the string to int
	input := orderID
	parts := strings.Split(input, "_") // Split by "_"
	lastPart := parts[len(parts)-1]

	err := Dbpool.QueryRow(context.Background(), query, lastPart).Scan(&price)
	if err != nil {
		log.Printf("Failed to fetch price for Order-ID %s: %v\n", orderID, err)
		return 0, err
	}

	fmt.Printf("Retrieved total price for Order-ID %s: %.2f\n", orderID, price)
	return price, nil
}

func GetAll(Dbpool *pgxpool.Pool) ([]ord.Order, error) {
	allrow := `SELECT * from "order"`

	rows, err := Dbpool.Query(context.Background(), allrow)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders: %w", err)
	}
	defer rows.Close() // Close the rows iterator

	var orders []ord.Order

	// Iterate over rows and scan into Order struct

	for rows.Next() {
		var o ord.Order

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

func CreateTable(Dbpool *pgxpool.Pool) {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS "order" (
		user_id SERIAL PRIMARY KEY,
		product_id INT NOT NULL,
		quantity INT NOT NULL,
		total_price FLOAT NOT NULL
	);
	`
	var err error
	_, err = Dbpool.Exec(context.Background(), createTableQuery)
	if err != nil {
		log.Fatalf("Failed to create table: %v\n", err)
	}

	fmt.Println("Table 'order' created successfully!")
}

func CreateOCD(Dbpool *pgxpool.Pool) {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS "ordercreated" (
		order_id VARCHAR(20) PRIMARY KEY,
		status VARCHAR(20)
	);
	`
	var err error
	_, err = Dbpool.Exec(context.Background(), createTableQuery)
	if err != nil {
		log.Fatalf("Failed to create table: %v\n", err)
	}

	fmt.Println("Table 'OrderCreated' created successfully!")
}

func CreatePaymentTable(Dbpool *pgxpool.Pool) {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS "paymentstatus" (
		order_id VARCHAR(20) PRIMARY KEY,
		payment_status VARCHAR(20)
	);
	`
	var err error
	_, err = Dbpool.Exec(context.Background(), createTableQuery)
	if err != nil {
		log.Fatalf("Failed to create table: %v\n", err)
	}

	fmt.Println("Table 'payment status' created successfully!")
}

/*func getTableSchema(dbpool *pgxpool.Pool, tableName string) {
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
*/
