package models

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

func ImportFromMySQL() {
	log.Printf("Connecting to MySQL")
	sqlDb, err := connectDB()
	if err != nil {
		panic(err)
	}
	defer sqlDb.Close()

	fmt.Printf("SELECT * FROM blogs")
	sqlSelect := `
			SELECT id, title, summary, slug, contentMd, createdOn, updatedOn, postedOn
			FROM blogs
			ORDER BY postedOn DESC`
	rows, err := sqlDb.Query(sqlSelect)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var id int64
	var title, summary, slug, content sql.NullString
	var createdOn, updatedOn, postedOn mysql.NullTime
	for rows.Next() {
		if err := rows.Scan(&id, &title, &summary, &slug, &content, &createdOn, &updatedOn, &postedOn); err != nil {
			log.Printf("Error on rows.Scan")
			panic(err)
		}

		createdOn := timeValue(createdOn)
		date := createdOn[0:10]
		time := createdOn[11:19]

		// Save a shell record in the database
		// (with the proper generated id based on the original)
		entry, err := textDb.NewEntryFor(date, time)
		if err != nil {
			log.Printf("Error on textDb.NewEntryFor")
			panic(err)
		}

		// Populate the rest of the fields and save through our model Blog class
		blog := newBlogFromEntry(entry)
		blog.Title = stringValue(title)
		blog.ContentMarkdown = stringValue(content)
		blog.Summary = stringValue(summary)
		blog.CreatedOn = date + " " + time
		blog.PostedOn = timeValue(postedOn)
		blog.UpdatedOn = timeValue(updatedOn)
		b2, err := blog.SaveForce()
		if err != nil {
			log.Printf("Error on blog.Save")
			panic(err)
		}
		fmt.Printf("Saved %s %s\n", b2.Slug, b2.Id)
	}
	log.Printf("Done")
	return
}
