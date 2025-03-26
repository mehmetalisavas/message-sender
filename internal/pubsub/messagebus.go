package pubsub

import "sync"

type MessageBus struct {
	channels map[string]chan interface{}
	mu       sync.RWMutex
}

func NewMessageBus() *MessageBus {
	return &MessageBus{
		channels: make(map[string]chan interface{}),
	}
}

func (mb *MessageBus) RegisterChannel(channel string, bufferSize int) {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	mb.channels[channel] = make(chan interface{}, bufferSize)
}

func (mb *MessageBus) GetChannel(name string) (chan interface{}, bool) {
	mb.mu.RLock()
	defer mb.mu.RUnlock()

	ch, ok := mb.channels[name]
	return ch, ok
}
