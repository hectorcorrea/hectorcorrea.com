package models

import (
	"database/sql"
	"fmt"
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

func (b Blog) URL(base string) string {
	return fmt.Sprintf("%s/blog/%s/%d", base, b.Slug, b.Id)
}

// RFC 1123Z looks like "Mon, 02 Jan 2006 15:04:05 -0700"
// https://golang.org/pkg/time/
func (b Blog) PostedOnRFC1123Z() string {
	layout := "2006-01-02 15:04:05 -0700 MST"
	t, err := time.Parse(layout, b.PostedOn)
	if err != nil {
		return ""
	}
	return t.Format(time.RFC1123Z)
}

func BlogGetAll(showDrafts bool) ([]Blog, error) {
	blogs, err := getAll(showDrafts)
	return blogs, err
}

func BlogGetById(id int64) (Blog, error) {
	blog, err := getOne(id)
	return blog, err
}

func BlogGetBySlug(slug string) (Blog, error) {
	id, err := getIdBySlug(slug)
	if err != nil {
		return Blog{}, err
	}
	return getOne(id)
}

func (b *Blog) beforeSave() error {
	b.Slug = getSlug(b.Title)
	b.UpdatedOn = dbUtcNow()
	return nil
}

func getSlug(title string) string {
	slug := strings.Trim(title, " ")
	slug = strings.ToLower(slug)
	slug = strings.Replace(slug, "c#", "c-sharp", -1)
	var chars []rune
	for _, c := range slug {
		isAlpha := c >= 'a' && c <= 'z'
		isDigit := c >= '0' && c <= '9'
		if isAlpha || isDigit {
			chars = append(chars, c)
		} else {
			chars = append(chars, '-')
		}
	}
	slug = string(chars)

	// remove double dashes
	for strings.Index(slug, "--") > -1 {
		slug = strings.Replace(slug, "--", "-", -1)
	}

	if len(slug) == 0 || slug == "-" {
		return ""
	}

	// make sure we don't end with a dash
	if slug[len(slug)-1] == '-' {
		return slug[0 : len(slug)-1]
	}

	return slug
}

func SaveNew() (int64, error) {
	db, err := connectDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	sqlInsert := `
		INSERT INTO blogs(title, summary, slug, content, createdOn)
		VALUES(?, ?, ?, ?, ?)`
	result, err := db.Exec(sqlInsert, "new blog", "", "new-blog", "", dbUtcNow())
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (b *Blog) Save() error {
	db, err := connectDB()
	if err != nil {
		return err
	}
	defer db.Close()
	b.beforeSave()

	sqlUpdate := `
		UPDATE blogs
		SET title = ?, summary = ?, slug = ?, content = ?, updatedOn = ?
		WHERE id = ?`
	_, err = db.Exec(sqlUpdate, b.Title, b.Summary, b.Slug, b.Content, dbUtcNow(), b.Id)
	return err
}

func (b *Blog) Import() error {
	db, err := connectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// Recalculate the slug value but not the updatedOn.
	b.Slug = getSlug(b.Title)

	sqlUpdate := `
		INSERT INTO blogs(id, title, summary, slug, content, createdOn, updatedOn, postedOn)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = db.Exec(sqlUpdate, b.Id, b.Title, b.Summary, b.Slug, b.Content,
		b.CreatedOn, b.UpdatedOn, b.PostedOn)
	return err
}

func getOne(id int64) (Blog, error) {
	db, err := connectDB()
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

func getIdBySlug(slug string) (int64, error) {
	db, err := connectDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()
	var id int64
	sqlSelect := "SELECT id FROM blogs WHERE slug = ? LIMIT 1"
	row := db.QueryRow(sqlSelect, slug)
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func MarkAsPosted(id int64) (Blog, error) {
	db, err := connectDB()
	if err != nil {
		return Blog{}, err
	}
	defer db.Close()

	now := time.Now().UTC()
	sqlUpdate := "UPDATE blogs SET postedOn = ? WHERE id = ?"
	_, err = db.Exec(sqlUpdate, now, id)
	if err != nil {
		return Blog{}, err
	}
	return getOne(id)
}

func MarkAsDraft(id int64) (Blog, error) {
	db, err := connectDB()
	if err != nil {
		return Blog{}, err
	}
	defer db.Close()

	sqlUpdate := "UPDATE blogs SET postedOn = NULL WHERE id = ?"
	_, err = db.Exec(sqlUpdate, id)
	if err != nil {
		return Blog{}, err
	}
	return getOne(id)
}

func getAll(showDrafts bool) ([]Blog, error) {
	db, err := connectDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlSelect := ""
	if showDrafts {
		sqlSelect = `
			SELECT id, title, summary, slug, postedOn
			FROM blogs
			ORDER BY postedOn DESC`
	} else {
		sqlSelect = `
			SELECT id, title, summary, slug, postedOn
			FROM blogs
			WHERE postedOn IS NOT null
			ORDER BY postedOn DESC`
	}
	rows, err := db.Query(sqlSelect)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
