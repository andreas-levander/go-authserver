package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/georgysavva/scany/pgxscan"
)

type Database interface {
	GetUsers() []*User
	GetUser(string) []*User
	AddUser(User) error
}

type DB struct {
	pool *pgxpool.Pool
}

func Connect(DB_URL string) *DB {
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	dbpool, err := pgxpool.Connect(context.Background(), DB_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return &DB{pool: dbpool}
}

type User struct {
	User_id       int
	Username      string
	Password_hash string
	Roles         []string
}

func (db *DB) GetUsers() []*User {
	var users []*User
	err := pgxscan.Select(context.Background(), db.pool, &users, `SELECT * FROM users`)

	// for i, u := range users {
	// 	fmt.Println(i, *u)
	// }
	if err != nil {
		fmt.Fprintf(os.Stderr, "GetUsers failed: %v\n", err)
	}

	return users

}

func (db *DB) GetUser(username string) []*User {
	var user []*User
	err := pgxscan.Select(context.Background(), db.pool, &user, `SELECT * FROM users WHERE username=$1`, username)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetUsers failed: %v\n", err)
	}

	return user

}

func (db *DB) AddUser(user User) error {
	_, err := db.pool.Exec(context.Background(), `INSERT INTO users (username, password_hash, roles) VALUES ($1, $2, $3::text[])`, user.Username, user.Password_hash, user.Roles)

	if err != nil {
		fmt.Fprintf(os.Stderr, "AddUser failed: %v\n", err)
	}
	return err
}
