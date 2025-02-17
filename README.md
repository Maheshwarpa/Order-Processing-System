# Order-Processing-System
Project - 2

### Overview

This project is an Order Processing System that integrates with Apache Kafka and PostgreSQL to handle order management efficiently. It consists of multiple components, including an API server, Kafka consumers and producers, and database operations.

### Note on Program Functionality
This program is designed to handle only POST requests to the /Orders endpoint. Specifically, the program expects a JSON payload containing order data, which will be processed and saved into the database.

- POST Request: Sends order data to the server for storage and processing.
-- The payload should contain the order information in JSON format.
-- Upon receiving a valid POST request, the server will:
--- Save the order details in the database.
--- Trigger Kafka consumers and producers to handle related messaging and processing.

### Features

1. REST API to handle order requests

2. Kafka producer for publishing order events

3. Kafka consumer for processing order events

4. PostgreSQL integration for order storage

5. Payment status determination

6. Notification service integration

### Technologies Used

1. Golang

2. Gin Web Framework

3. Apache Kafka (Producer/Consumer)

4. PostgreSQL (pgx/v5)

5. IBM Sarama (Kafka client for Go)

### Project Structure

```
├── API
│   ├── main.go (API Server)
├── module
│   ├── Consumer
│   │   ├── consumer.go (Kafka Consumer Logic)
│   ├── Producer
│   │   ├── producer.go (Kafka Producer Logic)
│   ├── DatabaseConn
│   │   ├── database.go (Database Connection & Queries)
│   ├── Orders
│   │   ├── orders.go (Order Model)
│   ├── Event
│   │   ├── event.go (Event Processing Logic)
│   ├── Notification
│   │   ├── notification.go (Notification Service)

```

### Detailed Function Descriptions

#### API

- main.go - This file contains the main entry point for the API server. Functions include:

- setupRouter(): Configures API routes and handlers.

- handleOrderRequest(): Processes incoming order requests.

- startServer(): Starts the API server on the specified port.

#### Consumer

- consumer.go - Implements Kafka consumer logic. Functions include:

- consumeOrders(): Listens to the Kafka topic for new order messages.

- processOrderMessage(): Parses and processes the consumed Kafka messages.

#### Producer

- producer.go - Implements Kafka producer logic. Functions include:

- publishOrderEvent(): Publishes order events to the Kafka topic.

- createKafkaProducer(): Initializes a Kafka producer instance.

#### DatabaseConn

- database.go - Manages database connections and queries. Functions include:

- connectDB(): Establishes a connection to the PostgreSQL database.

- insertOrder(): Inserts a new order record into the database.

- fetchOrderStatus(): Retrieves the status of a given order.

#### Orders

- orders.go - Defines the order model and related functions. Functions include:

- createOrder(): Constructs an order object.

- validateOrder(): Ensures order details are correct before processing.

#### Event

- event.go - Handles event processing logic. Functions include:

- processOrderEvent(): Processes order-related events received from Kafka.

- updateOrderStatus(): Updates the order status based on event outcomes.

#### Notification

- notification.go - Manages notifications. Functions include:

- sendNotification(): Sends an order confirmation or status update notification.

- notifyPaymentStatus(): Notifies users about payment status updates.

### Setup Instructions

#### Prerequisites

- Install Go

- Install PostgreSQL

- Install and configure Apache Kafka

### Docker Compose Setup for PostgreSQL, Kafka, Zookeeper, and Redis
This repository includes a Docker Compose configuration to set up a development environment with the following services:

- PostgreSQL: A relational database used for storing application data.
- Kafka: A distributed streaming platform used for building real-time data pipelines and streaming applications.
- Zookeeper: A centralized service for maintaining configuration information, naming, and providing distributed synchronization.
- Redis: An in-memory data structure store used as a cache, message broker, and for various other purposes.

#### Run Docker Compose
Use the following command to start all the services defined in the docker-compose.yml file:
```
docker-compose up

```
This command will pull the necessary Docker images (PostgreSQL, Kafka, Zookeeper, and Redis), create the containers, and start them up.

#### Running the Program
Once your environment is set up and all dependencies are installed, navigate to the directory containing main.go and run the following command

```
go run main.go
```

This will start the Gin HTTP server on localhost:8080.

#### Testing the POST Endpoint
Once the server is running, you can test the POST endpoint /Orders by sending a request with curl or using a tool like Postman.

- Using curl:

```
curl -X POST http://localhost:8080/Orders -H "Content-Type: application/json" -d '{"user_id":25, "Product_Id": 12345, "Quantity": 10, "Total_Price": 100.00}'
```

#### Access the Services
The services will be available on the following ports:

- PostgreSQL: localhost:5432
- Kafka: localhost:9092
- Zookeeper: localhost:2181
- Redis: localhost:6379
You can connect to these services using any PostgreSQL client, Kafka client, Redis client, or any other related tool.

####  add an order, make a POST request to the following endpoint

```
POST http://localhost:8080/Orders

```

#### Request Body Example

```
{
        "user_id": 9999,
        "product_id": 25,
        "quantity": 49,
        "total_price": 1500.33
        
}

```

The server will respond with:

- 200 OK if the order is successfully added and processed.
- 400 Bad Request if there is an error with the request, such as invalid JSON or an issue with adding the row to the database.

#### Stopping the Services
To stop the services and remove the containers, use the following command:
```
docker-compose down

```
