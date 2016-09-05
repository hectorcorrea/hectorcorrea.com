package models

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type DbSettings struct {
	driver     string
	user       string
	password   string
	database   string
	connString string
}

var dbSettings DbSettings

func InitDB() {
	dbSettings = DbSettings{
		driver:   env("DB_DRIVER", "mysql"),
		user:     env("DB_USER", "root"),
		password: env("DB_PASSWORD", ""),
		database: env("DB_NAME", "blogdb"),
	}
	dbSettings.connString = fmt.Sprintf("%s:%s@/%s?parseTime=true", dbSettings.user, dbSettings.password, dbSettings.database)
}

func SafeConnString() string {
	return fmt.Sprintf("%s:%s@/%s", dbSettings.user, "***", dbSettings.database)
}

func ConnectDB() (*sql.DB, error) {
	return sql.Open(dbSettings.driver, dbSettings.connString)
}

func env(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}
