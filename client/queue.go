package client

import "sync"

// Queue is a data structure that follows the first-in-first-out principle.
// Here, it is used to store a queue of messages that need to be sent to the upstream SMTP server
type Queue struct {
	data  []Message
	mutex sync.Mutex
}

// Empty returns true if the queue is empty
func (s *Queue) Empty() bool {
	return len(s.data) == 0
}

// Enqueue adds a new item to the back of the queue.
func (s *Queue) Enqueue(m Message) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data = append(s.data, m)
}

// Dequeue takes the item from the front of the queue and removes it
func (s *Queue) Dequeue() *Message {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Return nil if queue is empty
	if s.Empty() {
		return nil
	}

	// Extract first element and shuffle everything down
	m := s.data[0]
	s.data = s.data[1:]
	return &m
}

// DequeueAll empties the queue
func (s *Queue) DequeueAll() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data = []Message{}
}

// Len returns the length of the queue
func (s *Queue) Len() int {
	return len(s.data)
}

// GetAll returns the currently queued messages
// To avoid unnecessary overhead, this does not use mutex at all, which means it's not necessarily accurate
func (s *Queue) GetAll() []Message {
	return s.data
}
