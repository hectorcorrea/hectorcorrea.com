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
	var blogs []Blog
	for _, entry := range textDb.All() {
		if entry.IsDraft() {
			blog := newBlogFromEntry(entry)
			blogs = append(blogs, blog)
		}
	}
	return blogs, nil
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

func SaveNew() (string, error) {
	entry, err := textDb.NewEntry()
	return entry.Id, err
}

func (b *Blog) Save() error {
	// b.beforeSave()
	metadata := textdb.Metadata{
		Title:   b.Title,
		Summary: b.Summary,
	}
	entry := textdb.TextEntry{
		Metadata: metadata,
		Content:  b.ContentMarkdown,
		Id:       b.Id,
	}
	entry, err := textDb.UpdateEntry(entry)
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
	blog.Summary = entry.Metadata.Summary
	blog.Slug = entry.Metadata.Slug
	blog.ContentMarkdown = entry.Content
	blog.CreatedOn = entry.Metadata.CreatedOn
	blog.UpdatedOn = entry.Metadata.UpdatedOn
	blog.PostedOn = entry.Metadata.PostedOn

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

func MarkAsPosted(id string) (Blog, error) {
	entry, err := textDb.FindById(id)
	if err != nil {
		return Blog{}, err
	}

	entry.MarkAsPosted()
	_, err = textDb.UpdateEntry(entry)
	if err != nil {
		return Blog{}, err
	}
	return getOne(id)
}

func MarkAsDraft(id string) (Blog, error) {
	entry, err := textDb.FindById(id)
	if err != nil {
		return Blog{}, err
	}

	entry.MarkAsDraft()
	_, err = textDb.UpdateEntry(entry)
	if err != nil {
		return Blog{}, err
	}
	return getOne(id)
}

func getAll(showDrafts bool) ([]Blog, error) {
	var blogs []Blog
	for _, entry := range textDb.All() {
		if showDrafts || !entry.IsDraft() {
			blog := newBlogFromEntry(entry)
			blogs = append(blogs, blog)
		}
	}
	return blogs, nil
}

func newBlogFromEntry(entry textdb.TextEntry) Blog {
	blog := Blog{
		Id:        entry.Id,
		Title:     entry.Metadata.Title,
		Summary:   entry.Metadata.Summary,
		Slug:      entry.Metadata.Slug,
		CreatedOn: entry.Metadata.CreatedOn,
		UpdatedOn: entry.Metadata.UpdatedOn,
		PostedOn:  entry.Metadata.PostedOn,
	}
	return blog
}
