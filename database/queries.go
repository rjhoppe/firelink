package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Example 1
func CreateTableQuery(p *pgxpool.Pool) {
	_, err := p.Exec(context.Background(), "CREATE TABLE users (id SERIAL PRIMARY KEY,name VARCHAR(255) NOT NULL,email VARCHAR(255) UNIQUE NOT NULL);")
	if err != nil {
		log.Fatal("Error while creating the table")
	}
}

// Example 2
func InsertQuery(p *pgxpool.Pool) {
	_, err := p.Exec(context.Background(), "insert into users(name, email) values($1, $2)", "John", "johnysinsj@astronaut.com")
	if err != nil {
		log.Fatal("Error while inserting value into the table")
	}
}

// Insert drink recipe
// Insert dinner recipe
// Backup DB
