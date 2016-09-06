package models

import "log"

func CreateDefaultUser() error {
	db, err := ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	row := db.QueryRow("SELECT count(*) FROM users")
	var count int
	err = row.Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		log.Printf("Creating default blog user...")
		login := env("BLOG_USR", "user1")
		password := env("BLOG_PASS", "welcome1")
		// TODO: hash the password
		sqlInsert := `INSERT INTO users(login, name, password) VALUES(?, ?, ?)`
		_, err = db.Exec(sqlInsert, login, login, password)
		return err
	}
	return nil
}

func GetUserId(login string) (int64, error) {
	db, err := ConnectDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	row := db.QueryRow("SELECT id FROM users WHERE login = ?", login)
	var id int64
	err = row.Scan(&id)
	if err != nil {
		log.Printf("Error fething id for user: %s", login)
	}
	return id, err
}
