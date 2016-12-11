package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

	log.Printf("Importing\r\n%s", legacy)
	blog := Blog{}
	blog.Id = int64(legacy.Key)
	blog.Title = legacy.Title
	blog.Summary = legacy.Summary
	blog.Content = "TBD"
	blog.CreatedOn = "2016-12-12" // legacy.CreatedOn
	blog.UpdatedOn = "2016-12-12" // legacy.UpdatedOn
	blog.PostedOn = "2016-12-12"  // legacy.PostedOn
	err = blog.Import()
	if err != nil {
		log.Printf("ERROR: %s", err)
	}
	return err
}
