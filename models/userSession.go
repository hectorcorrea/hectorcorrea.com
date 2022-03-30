package models

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

const sessionDuration time.Duration = time.Hour * 24

type UserSession struct {
	SessionId string `xml:"sessionId"`
	Login     string `xml:"login"`
	CreatedOn string `xml:"createdOn"`
	ExpiresOn string `xml:"expiresOn"`
}

func (userSession *UserSession) StillActive() bool {
	return userSession.ExpiresOn > nowUTC()
}

func (userSession *UserSession) Expired() bool {
	return !userSession.StillActive()
}

func GetUserSession(sessionId string) (UserSession, error) {
	if sessionId == "" {
		return UserSession{}, errors.New("No ID was received")
	}

	filename := filepath.Join(".", "session.xml")
	reader, err := os.Open(filename)
	if err != nil {
		log.Printf("Error reading session file: %s %s\n", filename, err)
	}
	defer reader.Close()

	// Read the bytes and unmarshall into our struct
	byteValue, _ := ioutil.ReadAll(reader)
	var userSession UserSession
	xml.Unmarshal(byteValue, &userSession)

	if userSession.SessionId != sessionId {
		return UserSession{}, errors.New(fmt.Sprintf("UserSession %s not found", sessionId))
	}

	if userSession.StillActive() {
		if userSession.slideExpiration() {
			// TODO: implement this
			// // ...update the expiration time in the DB
			// sqlUpdate := `UPDATE sessions SET expiresOn = ? WHERE id = ?`
			// _, err = db.Exec(sqlUpdate, s.ExpiresOn, s.SessionId)
			// if err != nil {
			// 	log.Printf("Error updating session: %s", err)
			// }
		}
		log.Printf("Session %s, OK: %s", userSession.Login, sessionId)
		return userSession, nil
	}

	log.Printf("Session %s expired: %s", userSession.Login, sessionId)
	return UserSession{}, errors.New(fmt.Sprintf("UserSession %s has already expired", sessionId))
}

func NewUserSession(login string) (UserSession, error) {
	// Create a new session...
	sessionId := newId()
	s := UserSession{
		SessionId: sessionId,
		Login:     login,
		CreatedOn: nowUTC(),
		ExpiresOn: expireUTC(),
	}

	// ...dump it to XML
	xmlDeclaration := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\r\n"
	buffer := bytes.NewBufferString(xmlDeclaration)
	encoder := xml.NewEncoder(buffer)
	encoder.Indent("  ", "    ")
	err := encoder.Encode(s)
	if err != nil {
		return UserSession{}, err
	}

	// ... and save it.
	filename := filepath.Join(".", "session.xml")
	return s, ioutil.WriteFile(filename, buffer.Bytes(), 0644)
}

// source: https://www.socketloop.com/tutorials/golang-how-to-generate-random-string
func newId() string {
	size := 32
	rb := make([]byte, size)
	_, err := rand.Read(rb)
	if err != nil {
		panic(fmt.Sprintf("Error generating new Id %s", err))
	}
	return base64.URLEncoding.EncodeToString(rb)
}

// SlideExpiration renews the expiration date of the session if it is
// past half its lifetime. Returns true if the session was renewed.
// (https://docs.microsoft.com/en-us/dotnet/api/system.web.security.formsauthentication.slidingexpiration?view=netframework-4.7.2)
func (s *UserSession) slideExpiration() bool {
	//TODO: implement this
	return false
	// now := time.Now().UTC()
	// halfLife := s.ExpiresOn.Add(-sessionDuration / 2)
	// if now.After(halfLife) {
	// 	s.ExpiresOn = s.ExpiresOn.Add(sessionDuration)
	// 	return true
	// }
	// return false
}

func nowUTC() string {
	const time_format_now string = "2006-01-02 15:04:05.000" // yyyy-mm-dd hh:mm:ss.xxx
	return time.Now().UTC().Format(time_format_now)
}

func expireUTC() string {
	const time_format_now string = "2006-01-02 15:04:05.000" // yyyy-mm-dd hh:mm:ss.xxx
	return time.Now().Add(sessionDuration).UTC().Format(time_format_now)
}
