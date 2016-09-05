package models

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Session struct {
	Id        string
	ExpiresOn time.Time
}

func IsValidSession(id string) bool {
	if id == "" {
		return false
	}

	db, err := ConnectDB()
	if err != nil {
		return false
	}
	defer db.Close()

	sqlSelect := `SELECT expiresOn FROM sessions WHERE id = ?`
	row := db.QueryRow(sqlSelect, id)
	var expiresOn mysql.NullTime
	err = row.Scan(&expiresOn)
	if err != nil {
		return false
	}

	// TODO: delete expired session
	return expiresOn.Valid && expiresOn.Time.After(time.Now())
}

func AddSession() (Session, error) {
	db, err := ConnectDB()
	if err != nil {
		return Session{}, err
	}
	defer db.Close()

	id, err := newId()
	if err != nil {
		return Session{}, err
	}

	s := Session{
		Id:        id,
		ExpiresOn: time.Now().Add(time.Hour * 2),
	}

	sqlInsert := `INSERT INTO sessions(id, expiresOn) VALUES(?, ?)`
	_, err = db.Exec(sqlInsert, s.Id, s.ExpiresOn)
	return s, err
}

func DeleteSession(id string) {
	db, err := ConnectDB()
	if err != nil {
		return
	}
	defer db.Close()

	sqlDelete := `DELETE FROM sessions WHERE id = ?`
	_, err = db.Exec(sqlDelete, id)
}

// source: https://www.socketloop.com/tutorials/golang-how-to-generate-random-string
func newId() (string, error) {
	size := 32
	rb := make([]byte, size)
	_, err := rand.Read(rb)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(rb), nil
}
