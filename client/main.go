package client

import (
	"sync"
	"time"

	"github.com/constellation-net/dispatch/env"

	"github.com/emersion/go-smtp"
)

// Scheduler is used to periodically dispatch emails to the upstream SMTP server
type Scheduler struct {
	Quitter chan bool

	queue Queue

	// This mutex is used to ensure multiple SendAll calls run synchronously
	mutex sync.Mutex
}

// Run starts a time.Ticker loop in a new goroutine to periodically dispatch the queue to the upstream server
func (s *Scheduler) Run() {
	t := time.NewTicker(env.Vars.DispatchInterval())

	go func() {
		for {
			select {
			case <-t.C:
				err := s.SendAll()
				if err != nil {
					panic(err)
				}

			case <-s.Quitter:
				t.Stop()
				return
			}
		}
	}()
}

// NewScheduler sets up a new scheduler
func NewScheduler() *Scheduler {
	return &Scheduler{
		Quitter: make(chan bool),
	}
}

// Enqueue adds a new message to the queue
func (s *Scheduler) Enqueue(_ string, o *smtp.MailOptions, r []Recipient, b []byte) {
	m := Message{
		Body:       string(b),
		Options:    o,
		Recipients: r,
	}
	s.queue.Enqueue(m)
}

// SendAll will attempt to dispatch all messages in the queue to the upstream server.
// This function can only run one at a time, any subsequent calls will block until the previous calls have completed.
// If an error occurs when sending a message, this call will stop, the failing message will not be sent again and all messages after it won't be sent until the next SendAll call.
func (s *Scheduler) SendAll() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Create a new SMTP client
	c, err := smtp.Dial(env.Vars.Upstream.Address())
	if err != nil {
		return err
	}
	defer func() {
		err = c.Quit()
		if err != nil {
			panic(err)
		}
	}()

	// Pull each message from the queue one-by-one
	for i := 0; i < s.queue.Len(); i++ {
		m := s.queue.Dequeue()
		if m != nil {
			err = m.Send(c)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
