package models

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Blog struct {
	Id        int64
	Title     string
	Summary   string
	Slug      string
	Content   string
	CreatedOn string
	UpdatedOn string
	PostedOn  string
}

func (b Blog) DebugString() string {
	str := fmt.Sprintf("Id: %d\nTitle: %s\nSummary: %s\nContent: %s\n",
		b.Id, b.Title, b.Summary, b.Content)
	return str
}

func BlogGetAll() ([]Blog, error) {
	blogs, err := getAll()
	return blogs, err
}

func BlogGetById(id int64) (Blog, error) {
	blog, err := getOne(id)
	return blog, err
}

func (b *Blog) beforeSave() error {
	b.calculateSlug()
	// b.UpdatedOn = time.Now().String()
	return nil
}

func (b *Blog) calculateSlug() {
	b.Slug = strings.Replace(b.Title, " ", "-", -1)
	log.Printf("calculated slug: %s", b.Slug)
	// TODO: handle other characters
}

func SaveNew() (int64, error) {
	db, err := ConnectDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	sqlInsert := `
		INSERT INTO blogs(title, summary, slug, content, createdOn)
		VALUES(?, ?, ?, ?, ?)`
	result, err := db.Exec(sqlInsert, "new blog", "", "new-blog", "", time.Now())
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (b *Blog) Save() error {
	db, err := ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()
	b.beforeSave()

	sqlUpdate := `
		UPDATE blogs
		SET title = ?, summary = ?, slug = ?, content = ?, updatedOn = ?
		WHERE id = ?`
	_, err = db.Exec(sqlUpdate, b.Title, b.Summary, b.Slug, b.Content, time.Now(), b.Id)
	return err
}

func getOne(id int64) (Blog, error) {
	db, err := ConnectDB()
	if err != nil {
		return Blog{}, err
	}
	defer db.Close()

	sqlSelect := `
		SELECT title, summary, slug, content,
			createdOn, updatedOn, postedOn
		FROM blogs
		WHERE id = ?`
	row := db.QueryRow(sqlSelect, id)

	var title, summary, slug, content sql.NullString
	var createdOn, updatedOn, postedOn mysql.NullTime
	err = row.Scan(&title, &summary, &slug, &content, &createdOn, &updatedOn, &postedOn)
	if err != nil {
		return Blog{}, err
	}

	var blog Blog
	blog.Id = id
	blog.Title = stringValue(title)
	blog.Summary = stringValue(summary)
	blog.Slug = stringValue(slug)
	blog.Content = stringValue(content)
	blog.CreatedOn = timeValue(createdOn)
	blog.UpdatedOn = timeValue(updatedOn)
	blog.PostedOn = timeValue(postedOn)
	return blog, nil
}

func getAll() ([]Blog, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlSelect := "SELECT id, title, summary, slug, postedOn FROM blogs order by postedOn desc"
	rows, err := db.Query(sqlSelect)
	if err != nil {
		return nil, err
	}

	var blogs []Blog
	var id int64
	var title, summary, slug sql.NullString
	var postedOn mysql.NullTime
	for rows.Next() {
		if err := rows.Scan(&id, &title, &summary, &slug, &postedOn); err != nil {
			return nil, err
		}
		blog := Blog{
			Id:       id,
			Title:    stringValue(title),
			Summary:  stringValue(summary),
			Slug:     stringValue(slug),
			PostedOn: timeValue(postedOn),
		}
		blogs = append(blogs, blog)
	}
	return blogs, nil
}
