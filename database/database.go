package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/georgysavva/scany/pgxscan"
)

//var dbpool *pgxpool.Pool

type DB struct {
	pool *pgxpool.Pool
}

func Connect() *DB{
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	dbpool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return &DB{pool: dbpool}
}

type User struct {
	User_id int
	Name string
	Password_hash string
}

func (db *DB) GetUsers() []*User {
	var users []*User
	err := pgxscan.Select(context.Background(), db.pool, &users, `SELECT * FROM users`)
	
	for i, u := range users {
		fmt.Println(i, *u)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "SELECT failed: %v\n", err)
	}

	return users

}