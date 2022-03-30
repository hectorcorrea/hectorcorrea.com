package models

import (
	"os"

	"github.com/hectorcorrea/textodb"
)

type DbSettings struct {
	driver     string
	user       string
	password   string
	database   string
	connString string
}

var textDb textodb.TextoDb

func InitDB() error {
	rootDir := env("DB_ROOT_DIR", "./data")
	textDb = textodb.InitTextDb(rootDir)

	return CreateDefaultUser()
}

func DbInfo() string {
	return textDb.RootDir
}

func env(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}
