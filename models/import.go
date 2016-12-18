package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type LegacyBlog struct {
	Key       int    `json:"key"`
	Title     string `json:"title"`
	Summary   string `json:"summary"`
	Url       string `json:"url"`
	CreatedOn string `json:"createdOn"`
	UpdatedOn string `json:"updatedOn"`
	PostedOn  string `json:"postedOn"`
	MongoId   string `json:"_id"`
}

func (b LegacyBlog) String() string {
	str := ""
	str += fmt.Sprintf("%d\r\n", b.Key)
	str += b.Title + "\r\n"
	str += b.Url + "\r\n"
	return str
}

func ImportOne(fileName string) error {
	log.Printf("Importing %s", fileName)
	raw, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return err
	}

	var legacy LegacyBlog
	err = json.Unmarshal(raw, &legacy)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return err
	}

	log.Printf("\timporting metadata\r\n%s", legacy)
	blog := Blog{}
	blog.Id = int64(legacy.Key)
	blog.Title = legacy.Title
	blog.Summary = importSummary(legacy.Summary)
	blog.CreatedOn = fromZDate(legacy.CreatedOn)
	blog.UpdatedOn = fromZDate(legacy.UpdatedOn)
	blog.PostedOn = fromZDate(legacy.PostedOn)

	log.Printf("\timporting content")
	contentFile := strings.Replace(fileName, ".json", ".html", 1)
	bytes, err := ioutil.ReadFile(contentFile)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return err
	}
	blog.Content = string(bytes)

	err = blog.Import()
	if err != nil {
		log.Printf("ERROR: %s", err)
	}
	return err
}

func importSummary(text string) string {
	if len(text) <= 255 {
		return text
	}
	log.Printf("INFO: trimmed summary %s", text)
	return text[0:250] + "..."
}

// Converts a date in Z format (yyyy-MM-ddTHH:mm:ss.xxxZ)
// to a date that we can save in MySQL. Removes the "T"
// separator, drops the milliseconds and the "Z" marker.
func fromZDate(zdate string) string {
	date := zdate[0:10] + " " + zdate[11:19]
	return date
}
