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

const sessionDuration time.Duration = time.Hour * 24

type UserSession struct {
	SessionId string
	CreatedOn time.Time
	ExpiresOn time.Time
	Login     string
}

func GetUserSession(sessionId string) (UserSession, error) {
	if sessionId == "" {
		return UserSession{}, errors.New("No ID was received")
	}

	db, err := connectDB()
	if err != nil {
		return UserSession{}, err
	}
	defer db.Close()

	sqlSelect := `
		SELECT expiresOn, createdOn, users.login
		FROM sessions
		 	INNER JOIN users ON sessions.userId = users.id
		WHERE sessions.id = ?`
	row := db.QueryRow(sqlSelect, sessionId)
	var expiresOn mysql.NullTime
	var createdOn mysql.NullTime
	var login sql.NullString
	err = row.Scan(&expiresOn, &createdOn, &login)
	if err != nil {
		log.Printf("Error on scan: %s", err)
		return UserSession{}, err
	}

	if expiresOn.Valid && expiresOn.Time.After(time.Now().UTC()) {
		// Session is valid on the database
		s := UserSession{
			SessionId: sessionId,
			ExpiresOn: expiresOn.Time,
			CreatedOn: createdOn.Time,
			Login:     stringValue(login),
		}
		if s.slideExpiration() {
			// ...update the expiration time in the DB
			sqlUpdate := `UPDATE sessions SET expiresOn = ? WHERE id = ?`
			_, err = db.Exec(sqlUpdate, s.ExpiresOn, s.SessionId)
			if err != nil {
				log.Printf("Error updating session: %s", err)
			}
		}
		return s, nil
	}

	return UserSession{}, errors.New("UserSession has already expired")
}

func NewUserSession(login string) (UserSession, error) {
	db, err := connectDB()
	if err != nil {
		return UserSession{}, err
	}
	defer db.Close()

	sessionId, err := newId()
	if err != nil {
		return UserSession{}, err
	}

	now := time.Now().UTC()
	s := UserSession{
		SessionId: sessionId,
		Login:     login,
		CreatedOn: now,
		ExpiresOn: now.Add(sessionDuration),
	}

	userId, err := GetUserId(login)
	if err != nil {
		return UserSession{}, err
	}

	err = cleanSessions(db, userId)
	if err != nil {
		log.Printf("Error cleaning older sessions for user %s, %s", login, err)
	}

	sqlInsert := `INSERT INTO sessions(id, userId, expiresOn, createdOn) VALUES(?, ?, ?, ?)`
	_, err = db.Exec(sqlInsert, s.SessionId, userId, s.ExpiresOn, s.CreatedOn)
	if err != nil {
		log.Printf("Error in SQL INSERT INTO sessions: %s", err)
	}
	return s, err
}

func DeleteUserSession(sessionId string) {
	db, err := connectDB()
	if err != nil {
		return
	}
	defer db.Close()

	sqlDelete := `DELETE FROM sessions WHERE id = ?`
	_, err = db.Exec(sqlDelete, sessionId)
}

func cleanSessions(db *sql.DB, userId int64) error {
	// All sessions for this user (regardless of expiration date)
	sqlDelete := "DELETE FROM sessions WHERE userId = ?"
	_, err := db.Exec(sqlDelete, userId)
	if err != nil {
		return err
	}

	// All expired sessions (regardless of the user)
	sqlDelete = "DELETE FROM sessions WHERE expiresOn < utc_timestamp()"
	_, err = db.Exec(sqlDelete)
	return err
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

// SlideExpiration renews the expiration date of the session if it is
// past half its lifetime. Returns true if the session was renewed.
// (https://docs.microsoft.com/en-us/dotnet/api/system.web.security.formsauthentication.slidingexpiration?view=netframework-4.7.2)
func (s *UserSession) slideExpiration() bool {
	now := time.Now().UTC()
	halfLife := s.ExpiresOn.Add(-sessionDuration / 2)
	if now.After(halfLife) {
		s.ExpiresOn = s.ExpiresOn.Add(sessionDuration)
		return true
	}
	return false
}
