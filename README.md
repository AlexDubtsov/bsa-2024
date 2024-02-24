# bsa-2024

## Author

- [Alex Dubtsov](https://www.linkedin.com/in/alex-dubtsov-191b2a114/)

## Objective

The objective of is to create a simple Bitcoin Wallet REST API based on a simplified Bitcoin transaction model.

The programming language used is [Go](https://go.dev/).

## How To Launch

### Notes

SQLite is used as Data Base.<br>
To install Data Base Driver use command:

```console
go get github.com/mattn/go-sqlite3
```

Gin Web Framework is used for processing HTTP requests.

- [gin-gonic](https://github.com/gin-gonic/gin)<br>
To install Framework use command:

```console
go get -u github.com/gin-gonic/gin
```

### Commands

1. Clone the repository

```code
git clone https://github.com/AlexDubtsov/bsa-2024
```

2. Build or run the project

```code
go run main.go
```

## Feature list

1. Check not spent funds amount: [localhost:8080/balance](http://localhost:8080/balance)

2. List of all transactions: [localhost:8080/transactions](http://localhost:8080/transactions)

3. Endpoint for the transfer request: [http://localhost:8080/transfer](http://localhost:8080/transfer)

JSON payload example:
- {
- "requested_amount\": \"0.001\"
- }

4. Endpoint for the supply request: [http://localhost:8080/supply](http://localhost:8080/supply)

JSON payload example:
- {
- "supplied_amount\": \"0.001\"
- }