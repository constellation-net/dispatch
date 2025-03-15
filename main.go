package main

import (
	"log"

	"github.com/constellation-net/dispatch/client"
	"github.com/constellation-net/dispatch/server"
)

var scheduler = client.NewScheduler()

func main() {
	// Recover from panics
	defer func() {
		if err := recover(); err != nil {
			log.Println("[ERROR] Recovered from panic:", err)
		}
	}()

	// Start the SMTP server
	err := server.Run(scheduler)
	if err != nil {
		panic(err)
	}
}
