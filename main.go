package main

import (
	"fmt"
	"log"
)

func main() {
	go func() {
		questions, _ := loadOpenSAT("OpenSAT.json")
		showStats(questions)
	}()

	port := ":8080"

	fmt.Println("starting server at ", port)
	log.Fatal(startServer(port))
}
