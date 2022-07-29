package config

import (
	"go-test/database"
	"go-test/tokens"
)

type Env struct {
	DB *database.DB
	Token *tokens.Keys
}