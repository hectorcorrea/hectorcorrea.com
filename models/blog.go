package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/hectorcorrea/tbd/textdb"
	"hectorcorrea.com/markdown"
)

type Blog struct {
	Id              string
	Title           string
	Summary         string
	Slug            string
	ContentHtml     string
	ContentMarkdown string
	CreatedOn       string
	UpdatedOn       string
	PostedOn        string
}

func (b Blog) DebugString() string {
	str := fmt.Sprintf("Id: %s\nTitle: %s\nSummary: %s\n",
		b.Id, b.Title, b.Summary)
	return str
}

func (b Blog) IsDraft() bool {
	return false
	return b.PostedOn == ""
}

func (b Blog) URL(base string) string {
	return fmt.Sprintf("%s/blog/%s/%s", base, b.Slug, b.Id)
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

func BlogGetDrafts() ([]Blog, error) {
	// TODO: support drafts
	return []Blog{}, nil
}

func BlogGetById(id string) (Blog, error) {
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
	var parser markdown.Parser
	b.ContentHtml = parser.ToHtml(b.ContentMarkdown)
	return nil
}

func SaveNew() (string, error) {
	entry, err := textDb.NewEntry()
	return entry.Id, err
}

func (b *Blog) Save() error {
	b.beforeSave()
	metadata := textdb.Metadata{Title: b.Title}
	entry := textdb.TextEntry{
		Metadata: metadata,
		Content:  b.ContentMarkdown,
		Id:       b.Id,
	}
	entry, err := textDb.UpdateEntry(entry)
	return err
}

func (b *Blog) Import() error {
	db, err := connectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// Recalculate the slug value but not the updatedOn.
	// b.Slug = getSlug(b.Title)

	sqlUpdate := `
		INSERT INTO blogs(id, title, summary, slug, content, contentMd, createdOn, updatedOn, postedOn)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = db.Exec(sqlUpdate, b.Id, b.Title, b.Summary, b.Slug,
		b.ContentHtml, b.ContentMarkdown, b.CreatedOn, b.UpdatedOn, b.PostedOn)
	return err
}

func getOne(id string) (Blog, error) {
	entry, err := textDb.FindById(id)
	if err != nil {
		return Blog{}, err
	}

	var blog Blog
	blog.Id = entry.Id
	blog.Title = entry.Metadata.Title
	blog.Summary = "pending"
	blog.Slug = entry.Metadata.Slug
	blog.ContentMarkdown = entry.Content
	blog.CreatedOn = entry.Metadata.CreatedOn
	blog.UpdatedOn = entry.Metadata.UpdatedOn
	// blog.PostedOn = timeValue(postedOn)

	var parser markdown.Parser
	blog.ContentHtml = parser.ToHtml(entry.Content)
	return blog, nil
}

func getIdBySlug(slug string) (string, error) {
	entry, found := textDb.FindBySlug(slug)
	if !found {
		return "", errors.New(fmt.Sprintf("Slug not found: %s", slug))
	}
	return entry.Id, nil
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
	return getOne("TODO")
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
	return getOne("TODO")
}

func getAll(showDrafts bool) ([]Blog, error) {
	//TODO drafts
	return getMany("")
}

func getMany(sqlSelect string) ([]Blog, error) {
	var blogs []Blog
	for _, entry := range textDb.All() {
		blog := Blog{
			Id:       entry.Id,
			Title:    entry.Metadata.Title,
			Summary:  "pending",
			Slug:     entry.Metadata.Slug,
			PostedOn: entry.Metadata.CreatedOn,
		}
		blogs = append(blogs, blog)
	}
	return blogs, nil
}
