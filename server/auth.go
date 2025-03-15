package server

import (
	"errors"

	"github.com/constellation-net/dispatch/env"

	"golang.org/x/crypto/bcrypt"
)

// ValidCredentials return true if the given username and password match the configured credentials
func ValidCredentials(u string, p string) bool {
	return u == env.Vars.Relay.Username && ValidPassword([]byte(p))
}

// ValidPassword validates just the password given against the hash provided in env.Vars
func ValidPassword(p []byte) bool {
	err := bcrypt.CompareHashAndPassword(env.Vars.Relay.Password, p)
	return err == nil
}

func PlainAuth(_ string, u string, p string) error {
	if !ValidCredentials(u, p) {
		return errors.New("invalid username or password")
	}

	return nil
}
