package server

import (
	"github.com/emersion/go-smtp"
)

type Backend struct{}

// NewSession is called after client greeting (EHLO, HELO).
func (b *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &Session{
		conn: c,
	}, nil
}
