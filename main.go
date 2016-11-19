package main

import (
	"flag"
	"hectorcorrea.com/web"
)

func main() {
	var address = flag.String("address", "localhost:9001", "Address where server will listen for connections")
	flag.Parse()
	web.StartWebServer(*address)
}
