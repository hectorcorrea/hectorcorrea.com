package models

import (
	"bytes"
	"crypto/sha256"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type User struct {
	Login    string `xml:"login"`
	Password string `xml:"passwordHashed"`
}

func CreateDefaultUser() error {
	user, err := ReadUser()
	if err != nil {
		user = User{
			Login:    defaultUser(),
			Password: defaultPassword(),
		}
		log.Printf(fmt.Sprintf("Creating initial blog user: %s", user.Login))
		return SaveUser(user)
	}
	return nil
}

func ReadUser() (User, error) {
	filename := filepath.Join(".", "user.xml")
	reader, err := os.Open(filename)
	if err != nil {
		log.Printf("Error reading user file: %s %s\n", filename, err)
		return User{}, err
	}
	defer reader.Close()

	// Read the bytes and unmarshall into our struct
	byteValue, _ := ioutil.ReadAll(reader)
	var user User
	xml.Unmarshal(byteValue, &user)
	return user, nil
}

func SaveUser(user User) error {
	xmlDeclaration := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\r\n"
	buffer := bytes.NewBufferString(xmlDeclaration)
	encoder := xml.NewEncoder(buffer)
	encoder.Indent("  ", "    ")
	err := encoder.Encode(user)
	if err != nil {
		return err
	}

	// ... and save it.
	filename := filepath.Join(".", "user.xml")
	return ioutil.WriteFile(filename, buffer.Bytes(), 0644)
}

func SetPassword(login, newPassword string) error {
	user := User{
		Login:    login,
		Password: hashPassword(newPassword),
	}
	return SaveUser(user)
}

func defaultUser() string {
	return env("BLOG_USR", "user1")
}

func defaultPassword() string {
	password := env("BLOG_PASS", "welcome1")
	return hashPassword(password)
}

func hashPassword(password string) string {
	salt := env("BLOG_SALT", "")
	salted := password + salt
	data := []byte(salted)
	hashed := sha256.Sum256(data)
	return fmt.Sprintf("%x", hashed)
}

func LoginUser(login, password string) (bool, error) {
	user, err := ReadUser()
	if err != nil {
		return false, err
	}

	if login != user.Login {
		log.Printf("Unknown user: %s", login)
		return false, errors.New("Unknown user")
	}

	hashedPassword := hashPassword(password)
	if hashedPassword != user.Password {
		log.Printf("Wrong password for user: %s", login)
		return false, errors.New("Wrong password")
	}

	return true, nil
}
