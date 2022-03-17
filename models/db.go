package models

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/hectorcorrea/tbd/textdb"
)

type DbSettings struct {
	driver     string
	user       string
	password   string
	database   string
	connString string
}

var dbSettings DbSettings
var textDb textdb.TextDb

func InitDB() error {
	rootDir := env("DB_ROOT_DIR", "./textdb")
	textDb = textdb.InitTextDb(rootDir)

	dbSettings = DbSettings{
		driver:   env("DB_DRIVER", "mysql"),
		user:     env("DB_USER", "root"),
		password: env("DB_PASSWORD", ""),
		database: env("DB_NAME", "blogdb"),
	}
	dbSettings.connString = fmt.Sprintf("%s:%s@/%s?parseTime=true", dbSettings.user, dbSettings.password, dbSettings.database)
	return CreateDefaultUser()
}

func DbConnStringSafe() string {
	return fmt.Sprintf("%s:%s@/%s", dbSettings.user, "***", dbSettings.database)
}

func connectDB() (*sql.DB, error) {
	return sql.Open(dbSettings.driver, dbSettings.connString)
}

func env(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}

func timeValue(t mysql.NullTime) string {
	if t.Valid {
		return t.Time.String()
	}
	return ""
}

func stringValue(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}
