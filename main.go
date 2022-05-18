package main

import (
	"flag"
	"fmt"
	"log"
)

var port int

func main() {
	portFlag := flag.Int("port", 80, "Defines the port where the server listen")
	flag.Parse()
	port = *portFlag

	log.Println("Listening on port 3000")
	log.Fatal(App.Listen(fmt.Sprintf(":%d", port)))
}
