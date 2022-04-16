package models

import (
	"errors"
	"fmt"
	"sort"
	"strconv"

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
	Type            string // blog or page
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

// Records that are blog post entries published
func BlogGetPosts() []Blog {
	var blogs []Blog
	for _, entry := range textDb.All() {
		blog := newBlogFromEntry(entry)
		if !blog.IsDraft() && blog.Type == "blog" {
			blogs = append(blogs, blog)
		}
	}

	var sorted BlogSort = blogs
	sort.Sort(sorted)
	return sorted
}

// Records that are blog post entries not published
func BlogGetDrafts() []Blog {
	var blogs []Blog
	for _, entry := range textDb.All() {
		blog := newBlogFromEntry(entry)
		if blog.IsDraft() {
			blogs = append(blogs, blog)
		}
	}
	return blogs
}

// Records that are system pages (not blog posts) like about and home.
func BlogGetPages() []Blog {
	var blogs []Blog
	for _, entry := range textDb.All() {
		blog := newBlogFromEntry(entry)
		if blog.Type == "page" {
			blogs = append(blogs, blog)
		}
	}
	return blogs
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
	if b.Type == "page" {
		entry.SetField("type", "page")
	} else {
		entry.SetField("type", "blog")
	}
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

	blog := newBlogFromEntry(entry)
	blog.ContentMarkdown = entry.Content()

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
	if entry.GetField("type") == "page" {
		blog.Type = "page"
	} else {
		blog.Type = "blog"
	}
	return blog
}
