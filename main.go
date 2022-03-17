package main

import (
	"flag"
	"log"

	"hectorcorrea.com/models"
	"hectorcorrea.com/web"
)

func main() {
	var address = flag.String("address", "localhost:9001", "Address where server will listen for connections")
	var importer = flag.Bool("import", false, "True to import MySQL data")
	var resave = flag.String("resave", "", "Resaves all blog posts to force recalculate of HTML content")
	flag.Parse()
	if *importer == true {
		importData()
		return
	}

	if *resave != "" {
		resaveAll()
		return
	}

	web.StartWebServer(*address)
}

func importData() {
	if err := models.InitDB(); err != nil {
		log.Print("ERROR: Failed to initialize database: ", err)
	}
	log.Printf("Database: %s", models.DbConnStringSafe())
	models.ImportFromMySQL()
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
