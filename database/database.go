package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main(requestType string) {

	// Create database connection
	connPool, err := pgxpool.NewWithConfig(context.Background(), DbConfig())
	if err != nil {
		log.Fatal("Error while creating connection to the database!!")
	}

	connection, err := connPool.Acquire(context.Background())
	if err != nil {
		log.Fatal("Error while acquiring connection from the database pool!!")
	}
	defer connection.Release()

	err = connection.Ping(context.Background())
	if err != nil {
		log.Fatal("Could not ping database")
	}

	fmt.Println("Connected to the database!!")

	// Database queries
	CreateTableQuery(connPool)
	InsertQuery(connPool)

	defer connPool.Close()

	// function here
}
