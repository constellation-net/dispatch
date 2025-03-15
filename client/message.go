package client

import (
	"io"

	"github.com/constellation-net/dispatch/env"

	"github.com/emersion/go-smtp"
)

// Message represents the details of a message to be passed on to the upstream server
type Message struct {
	Body       string
	Recipients []Recipient
	Options    *smtp.MailOptions
}

// Prepare sends the source and recipient addresses to the server, and returns an io.WriteCloser for writing the body to
func (m *Message) Prepare(c *smtp.Client) (io.WriteCloser, error) {
	// Add the source address
	err := c.Mail(env.Vars.Upstream.From, m.Options)
	if err != nil {
		return nil, err
	}

	// Add all recipients to the message
	err = m.AddRecipients(c)
	if err != nil {
		return nil, err
	}

	return c.Data()
}

// Send will prepare a message and then send it to the SMTP server to deliver to the recipients
func (m *Message) Send(c *smtp.Client) error {
	// Prep the message for sending
	wc, err := m.Prepare(c)
	if err != nil {
		return err
	}

	// Send the message body
	_, err = wc.Write([]byte(m.Body))
	if err != nil {
		return err
	}

	return wc.Close()
}

// AddRecipients sends each recipient to the SMTP server
func (m *Message) AddRecipients(c *smtp.Client) error {
	for _, r := range m.Recipients {
		err := c.Rcpt(r.Address, r.Options)
		if err != nil {
			return err
		}
	}

	return nil
}
