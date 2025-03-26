package pubsub

import "sync"

// MessageBus represents a message that is used to send messages from producers to consumers.
type MessageBus struct {
	channels map[string]chan interface{}
	mu       sync.RWMutex
}

// NewMessageBus creates a new MessageBus instance.
func NewMessageBus() *MessageBus {
	return &MessageBus{
		channels: make(map[string]chan interface{}),
	}
}

// RegisterChannel registers a new channel with the given name and bufferSize.
func (mb *MessageBus) RegisterChannel(channel string, bufferSize int) {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	mb.channels[channel] = make(chan interface{}, bufferSize)
}

// GetChannel returns the channel with the given name.
func (mb *MessageBus) GetChannel(name string) (chan interface{}, bool) {
	mb.mu.RLock()
	defer mb.mu.RUnlock()

	ch, ok := mb.channels[name]
	return ch, ok
}
