package main

import "log"

func main() {
    log.Println("Listening on port 3000")
    log.Fatal(App.Listen(":3000"))
}
