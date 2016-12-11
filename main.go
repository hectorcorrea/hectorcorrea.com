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
	flag.Parse()
	if *importer != "" {
		if err := models.InitDB(); err != nil {
			log.Fatal("Failed to initialize database: ", err)
		}
		log.Printf("Database: %s", models.DbConnStringSafe())

		models.ImportOne(*importer)
		return
	}
	web.StartWebServer(*address)
}
