package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

var globaldbpool *pgxpool.Pool

func ConnectToDatabase() *pgxpool.Pool {
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	dbpool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	globaldbpool = dbpool

	return dbpool
}

func GetUsers() {
	rows, err := globaldbpool.Query(context.Background(), "SELECT * from users")

	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	}

	fmt.Println(rows)

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
		  fmt.Println("error while iterating dataset")
		}

		fmt.Println(values)
	}
}