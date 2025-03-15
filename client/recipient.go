package client

import "github.com/emersion/go-smtp"

type Recipient struct {
	Address string
	Options *smtp.RcptOptions
}
