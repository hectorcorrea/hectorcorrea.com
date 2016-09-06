package models

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
)

type UserSession struct {
	SessionId string
	ExpiresOn time.Time
	Login     string
}

func GetUserSession(sessionId string) (UserSession, error) {
	if sessionId == "" {
		return UserSession{}, errors.New("No ID was received")
	}

	db, err := ConnectDB()
	if err != nil {
		return UserSession{}, err
	}
	defer db.Close()

	sqlSelect := `
		SELECT expiresOn, users.login
		FROM sessions
		 	INNER JOIN users ON sessions.userId = users.id
		WHERE sessions.id = ?`
	row := db.QueryRow(sqlSelect, sessionId)
	var expiresOn mysql.NullTime
	var login sql.NullString
	err = row.Scan(&expiresOn, &login)
	if err != nil {
		log.Printf("Error on scan: %s", err)
		return UserSession{}, err
	}

	if expiresOn.Valid && expiresOn.Time.After(time.Now()) {
		s := UserSession{SessionId: sessionId, ExpiresOn: expiresOn.Time, Login: stringValue(login)}
		return s, nil
	}

	return UserSession{}, errors.New("UserSession has already expired")
}

func NewUserSession(login string) (UserSession, error) {
	db, err := ConnectDB()
	if err != nil {
		return UserSession{}, err
	}
	defer db.Close()

	sessionId, err := newId()
	if err != nil {
		return UserSession{}, err
	}

	s := UserSession{
		SessionId: sessionId,
		Login:     login,
		ExpiresOn: time.Now().Add(time.Hour * 2),
	}

	userId, err := GetUserId(login)
	if err != nil {
		return UserSession{}, err
	}

	sqlInsert := `INSERT INTO sessions(id, userId, expiresOn) VALUES(?, ?, ?)`
	_, err = db.Exec(sqlInsert, s.SessionId, userId, s.ExpiresOn)
	if err != nil {
		log.Printf("Error in SQL INSERT INTO sessions: %s", err)
	}
	return s, err
}

func DeleteUserSession(sessionId string) {
	db, err := ConnectDB()
	if err != nil {
		return
	}
	defer db.Close()

	sqlDelete := `DELETE FROM sessions WHERE id = ?`
	_, err = db.Exec(sqlDelete, sessionId)
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
