package models

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/hectorcorrea/textodb"
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

// https://procrypt.github.io/post/2017-06-01-sorting-structs-in-golang/
type BlogSort []Blog

func (blogs BlogSort) Len() int {
	return len(blogs)
}

func (blogs BlogSort) Less(i, j int) bool {
	return blogs[i].PostedOn > blogs[j].PostedOn
}

func (blogs BlogSort) Swap(i, j int) {
	blogs[i], blogs[j] = blogs[j], blogs[i]
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

	var sorted BlogSort = blogs
	sort.Sort(sorted)

	return sorted, err
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

func BlogGetByOldId(oldId string) (Blog, error) {
	id := getNewIdForOldId(oldId)
	if id == "" {
		return Blog{}, errors.New(fmt.Sprintf("Old id (%s) not found", oldId))
	}
	return getOne(id)
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

func (b *Blog) Save() (Blog, error) {
	entry, err := textDb.FindById(b.Id)
	if err != nil {
		return Blog{}, err
	}

	entry.Title = b.Title
	entry.Summary = b.Summary
	entry.SetContent(b.ContentMarkdown)
	entry, err = textDb.UpdateEntry(entry)
	if err != nil {
		return Blog{}, err
	}

	blog := newBlogFromEntry(entry)
	return blog, nil
}

func (b *Blog) SaveForce(oldId string) (Blog, error) {
	entry, err := textDb.FindById(b.Id)
	if err != nil {
		return Blog{}, err
	}

	entry.SetContent(b.ContentMarkdown)
	entry.Title = b.Title
	entry.Summary = b.Summary
	entry.CreatedOn = b.CreatedOn
	entry.PostedOn = b.PostedOn
	entry.UpdatedOn = b.UpdatedOn
	entry.SetField("oldId", oldId)
	entry, err = textDb.UpdateEntryHonorDates(entry)
	if err != nil {
		return Blog{}, err
	}

	return getOne(b.Id)
}

func getOne(id string) (Blog, error) {
	entry, err := textDb.FindById(id)
	if err != nil {
		return Blog{}, err
	}

	var blog Blog
	blog.Id = entry.Id
	blog.ContentMarkdown = entry.Content()
	blog.Title = entry.Title
	blog.Summary = entry.Summary
	blog.Slug = entry.Slug
	blog.CreatedOn = entry.CreatedOn
	blog.UpdatedOn = entry.UpdatedOn
	blog.PostedOn = entry.PostedOn

	var parser markdown.Parser
	blog.ContentHtml = parser.ToHtml(entry.Content())
	return blog, nil
}

func getNewIdForOldId(oldId string) string {
	number, _ := strconv.Atoi(oldId)
	if number == 0 {
		return ""
	}
	entry, _ := textDb.FindBy("oldId", oldId)
	return entry.Id
}

func getIdBySlug(slug string) (string, error) {
	entry, found := textDb.FindBySlug(slug)
	if !found {
		return "", errors.New(fmt.Sprintf("Slug not found: %s", slug))
	}
	return entry.Id, nil
}

// Notice that this method reads the entry, updates on field, and resaves it.
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

// Notice that this method reads the entry, updates on field, and resaves it.
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

func newBlogFromEntry(entry textodb.TextoEntry) Blog {
	blog := Blog{
		Id:        entry.Id,
		Title:     entry.Title,
		Summary:   entry.Summary,
		Slug:      entry.Slug,
		CreatedOn: entry.CreatedOn,
		UpdatedOn: entry.UpdatedOn,
		PostedOn:  entry.PostedOn,
	}
	return blog
}
