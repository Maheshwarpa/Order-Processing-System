# Order-Processing-System
Project - 2

### Overview

This project is an Order Processing System that integrates with Apache Kafka and PostgreSQL to handle order management efficiently. It consists of multiple components, including an API server, Kafka consumers and producers, and database operations.

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
