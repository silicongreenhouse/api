package main

import (
	"flag"
	"fmt"
	"log"
)

var port int
var staticFolder string

func main() {
	portFlag := flag.Int("port", 80, "Defines the port where the server listen")
	staticFlag := flag.String("static", "/var/www/silicon_greenhouse/", "Defines where to check for static files to serve the control panel")
	flag.Parse()
	port = *portFlag
	staticFolder = *staticFlag

	fmt.Printf("Static folder: %s\n", staticFolder)

	App.Static("/", staticFolder)

	log.Printf("Listening on port %d\n", port)
	log.Fatal(App.Listen(fmt.Sprintf(":%d", port)))
}
