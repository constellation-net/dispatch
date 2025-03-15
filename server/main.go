package server

import (
	"context"
	"time"

	"github.com/constellation-net/dispatch/client"
	"github.com/constellation-net/dispatch/env"

	"github.com/emersion/go-smtp"
)

var (
	backend   *Backend
	server    *smtp.Server
	scheduler *client.Scheduler
)

func init() {
	// Prepare server
	backend = &Backend{}

	server = smtp.NewServer(backend)
	server.Addr = env.Vars.Relay.Address()
	server.Domain = env.Vars.Relay.Host
	server.WriteTimeout = 10 * time.Second
	server.ReadTimeout = 10 * time.Second
	server.MaxMessageBytes = 1024 * 1024
	server.MaxRecipients = 50
	server.AllowInsecureAuth = true
}

func Run(s *client.Scheduler) error {
	scheduler = s

	// Gracefully shut down the server when the main thread terminates
	defer func() {
		if err := server.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	}()

	// Start the SMTP server
	return server.ListenAndServe()
}
