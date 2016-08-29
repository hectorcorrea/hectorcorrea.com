package models

import "database/sql"

type Blog struct {
	Id        int
	Title     string
	Summary   string
	Url       string
	Content   string
	CreatedOn string
	PostedOn  string
}

func BlogGetAll() ([]Blog, error) {
	blogs, err := getAll()
	return blogs, err
}

func BlogGetById(id int) (Blog, error) {
	blog, err := getOne(id)
	return blog, err
}

func getOne(id int) (Blog, error) {
	db, err := ConnectDB()
	if err != nil {
		return Blog{}, err
	}

	var title, summary, url sql.NullString
	sqlSelect := "SELECT title, summary, url FROM blogs where id = ?"
	err = db.QueryRow(sqlSelect, id).Scan(&title, &summary, &url)
	if err != nil {
		return Blog{}, err
	}
	var blog Blog
	blog.Id = id
	blog.Title = stringValue(title)
	blog.Summary = stringValue(summary)
	blog.Url = stringValue(url)
	return blog, nil
}

func getAll() ([]Blog, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, err
	}

	sqlSelect := "SELECT id, title, summary, url FROM blogs order by postedOn desc"
	rows, err := db.Query(sqlSelect)
	if err != nil {
		return nil, err
	}

	var blogs []Blog
	var id int
	var title, summary, url sql.NullString
	for rows.Next() {
		if err := rows.Scan(&id, &title, &summary, &url); err != nil {
			return nil, err
		}
		blog := Blog{
			Id:      id,
			Title:   stringValue(title),
			Summary: stringValue(summary),
			Url:     stringValue(url),
		}
		blogs = append(blogs, blog)
	}
	return blogs, nil
}

func stringValue(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}
