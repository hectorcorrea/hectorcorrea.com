package main

import (
	"flag"
	"log"

	"hectorcorrea.com/models"
	"hectorcorrea.com/web"
)

func main() {
	var address = flag.String("address", "localhost:9001", "Address where server will listen for connections")
	var importer = flag.String("import", "", "Name of file to import legacy blog from")
	var resave = flag.String("resave", "", "Resaves all blog posts to force recalculate of HTML content")
	flag.Parse()
	if *importer != "" {
		panic("TODO: re-implement the import feature")
	}

	if *resave != "" {
		resaveAll()
		return
	}

	web.StartWebServer(*address)
}

func resaveAll() {
	if err := models.InitDB(); err != nil {
		log.Fatal("Failed to initialize database: ", err)
	}
	log.Printf("Database: %s", models.DbConnStringSafe())
	blogs, _ := models.BlogGetAll(true)
	for _, b := range blogs {
		log.Printf("re-saving %s - %s", b.Id, b.Title)
		blog, _ := models.BlogGetById(b.Id)
		blog.Save()
	}
}
