package models

// Code to produce the XML for an RSS feed
// Heavily based on: https://siongui.github.io/2015/02/27/go-parse-rss2/
//
// rss := NewRss("title", "desc", "url")
// rss.Add("t1", "d1", "u1")
// rss.Add("t2", "d2", "u2")
// rss.ToXml()
//

import (
	"bytes"
	"encoding/xml"
)

type ItemGuid struct {
	PermaLink bool   `xml:"isPermaLink,attr"`
	Link      string `xml:",chardata"`
}

type Item struct {
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	Guid        ItemGuid `xml:"guid"`
	PubDate     string   `xml:"pubDate"`
}

type AtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

type Channel struct {
	Title         string   `xml:"title"`
	Description   string   `xml:"description"`
	Link          string   `xml:"link"`
	Generator     string   `xml:"generator"`
	lastBuildDate string   `xml:"lastBuildDate"`
	AtomLink      AtomLink `xml:"atom:link"`
	ItemList      []Item   `xml:"item"`
}

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Dc      string   `xml:"xmlns:dc,attr"`
	Content string   `xml:"xmlns:content,attr"`
	Atom    string   `xml:"xmlns:atom,attr"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

func NewRss(title, description, link string) Rss {
	items := []Item{}

	atomLink := AtomLink{Href: link,
		Rel:  "self",
		Type: "application/rss+xml"}

	channel := Channel{Title: title,
		Description: description,
		Generator:   "Custom go code",
		AtomLink:    atomLink,
		ItemList:    items}

	rss := Rss{Channel: channel,
		Dc:      "http://purl.org/dc/elements/1.1/",
		Content: "http://purl.org/rss/1.0/modules/content/",
		Atom:    "http://www.w3.org/2005/Atom",
		Version: "2.0"}
	return rss
}

func (rss *Rss) Add(t, d, url, date string) {
	guid := ItemGuid{Link: url, PermaLink: true}
	item := Item{Title: t, Description: d, Guid: guid, PubDate: date}
	rss.Channel.ItemList = append(rss.Channel.ItemList, item)
}

func (rss Rss) ToXml() (string, error) {
	text := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\r\n"
	buffer := bytes.NewBufferString(text)
	enc := xml.NewEncoder(buffer)
	enc.Indent("  ", "    ")
	if err := enc.Encode(rss); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
