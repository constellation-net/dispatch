package server

import (
	"github.com/constellation-net/dispatch/client"
	"io"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

// Session handles an active connection from a client
type Session struct {
	conn *smtp.Conn

	From       string
	Options    *smtp.MailOptions
	Recipients []client.Recipient
	Body       []byte
}

// AuthMechanisms returns an array of valid auth mechanisms
func (s *Session) AuthMechanisms() []string {
	return []string{sasl.Plain}
}

// Auth is the handler for supported authenticators
func (s *Session) Auth(_ string) (sasl.Server, error) {
	return sasl.NewPlainServer(PlainAuth), nil
}

// Mail is called when a client wishes to send a new message. It contains the sender address.
func (s *Session) Mail(f string, o *smtp.MailOptions) error {
	s.From = f
	s.Options = o
	return nil
}

// Rcpt adds a new recipient to the current message. A client will call this once for every recipient.
func (s *Session) Rcpt(to string, o *smtp.RcptOptions) error {
	r := client.Recipient{
		Address: to,
		Options: o,
	}

	s.Recipients = append(s.Recipients, r)
	return nil
}

// Data marks the beginning of the message's body. This function will block until the client has finished sending the body.
func (s *Session) Data(r io.Reader) error {
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	s.Body = b
	s.Flush()
	return nil
}

// Reset is called whenever a client wishes to reset the session to a default state.
// This allows the client to send multiple messages in one connection.
func (s *Session) Reset() {
	s.Flush()
	s.From = ""
	s.Recipients = []client.Recipient{}
	s.Body = []byte{}
}

// Logout is sent by a client when it is ready to terminate the connection completely.
func (s *Session) Logout() error {
	s.Reset()
	return nil
}

// Flush copies the current message information over to the scheduler to be dispatched upstream.
// This should be called after Data completes.
// If Session.Recipients is empty when this is called, the function will return without doing anything
func (s *Session) Flush() {
	if len(s.Recipients) > 0 {
		scheduler.Enqueue(s.From, s.Options, s.Recipients, s.Body)
	}
}
