# Microservices Project

This project is a demonstration of a microservices architecture using Go, gRPC, GraphQL, and PostgreSQL. It consists of four main services: `account`, `order`, `catalog`, and a `graphql-gateway`.

## Architecture

The architecture is designed to be scalable and maintainable. Each service is responsible for a specific domain and communicates with others via gRPC. The `graphql-gateway` acts as a single entry point for clients, aggregating data from the other services.

Here is a high-level overview of the architecture: 
[text](https://whimsical.com/microservices-Xgqs3Q7ZjExthBh7qLZUS3)

```
+-------------------+      +-------------------+      +-------------------+
|                   |      |                   |      |                   |
|  Account Service  | <--- |  GraphQL Gateway  | ---> |   Order Service   |
|                   |      |                   |      |                   |
+-------------------+      +-------------------+      +-------------------+
        ^                            |                            ^
        |                            |                            |
        v                            v                            v
+-------------------+      +-------------------+      +-------------------+
|  Account Database |      |   Catalog Service |      |   Order Database  |
+-------------------+      +-------------------+      +-------------------+
                                     ^
                                     |
                                     v
                             +-------------------+
                             | Catalog Database  |
                             +-------------------+
```

### Services

*   **Account Service:** Manages user accounts.
*   **Order Service:** Manages orders and their products.
*   **Catalog Service:** Manages the product catalog.
*   **GraphQL Gateway:** A GraphQL server that acts as an API gateway for the other services.

## Features

*   User account creation and retrieval.
*   Product creation and retrieval.
*   Order creation and retrieval.
*   GraphQL API for all services.
*   gRPC for inter-service communication.
*   Dockerized setup for easy development and deployment.

## API Documentation

The GraphQL API is the single entry point for all client-side applications. It provides a flexible and powerful way to query and mutate data.

### Queries

**`accounts(pagination: PaginationInput, id: String): [Account!]!`**

Fetches a list of accounts. Can be paginated and filtered by ID.

**`products(pagination: PaginationInput, query: String, id: String): [Product!]!`**

Fetches a list of products. Can be paginated and filtered by name or ID.

### Mutations

**`createAccount(account: AccountInput!): Account`**

Creates a new account.

**`createProduct(product: ProductInput!): Product`**

Creates a new product.

**`createOrder(order: OrderInput!): Order`**

Creates a new order.

## Getting Started

To get started with this project, you'll need to have Docker and Docker Compose installed on your machine.

### Prerequisites

*   [Docker](https://www.docker.com/get-started)
*   [Docker Compose](https://docs.docker.com/compose/install/)
*   [Go](https://golang.org/doc/install)
*   [protoc](https://grpc.io/docs/protoc-installation/)

### Installation

1.  Clone the repository:

    ```bash
    git clone https://github.com/saurabhkr78/Microservices.git
    cd Microservices
    ```

2.  Run the services using Docker Compose:

    ```bash
    docker-compose up -d --build
    ```

## Running the Services

The services will be available at the following ports:

*   **GraphQL Gateway:** `http://localhost:8080`
*   **Account Service:** `http://localhost:8001`
*   **Order Service:** `http://localhost:8002`
*   **Catalog Service:** `http://localhost:8003`

You can access the GraphQL Playground at `http://localhost:8080` to interact with the API.

## Database

Each service has its own PostgreSQL database, which is also run in a Docker container. The data is persisted in Docker volumes.

*   **Account Database:** `postgresql://postgres:postgres@localhost:5433/account`
*   **Order Database:** `postgresql://postgres:postgres@localhost:5434/order`
*   **Catalog Database:** `postgresql://postgres:postgres@localhost:5435/catalog`
