package markdown

import (
	"io/ioutil"
	"testing"
)

func TestText(t *testing.T) {
	md := loadFromDisk()
	var parser Parser
	html := parser.ToHtml(md)
	saveToDisk(html)
	t.Errorf("%s", html)
}

func loadFromDisk() string {
	bytes, _ := ioutil.ReadFile("test.md")
	return string(bytes)
}

func saveToDisk(html string) {
	ioutil.WriteFile("test.html", []byte(html), 0644)
}
