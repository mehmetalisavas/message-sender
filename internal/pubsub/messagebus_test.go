package pubsub

import (
	"testing"
)

func TestMessageBus_RegisterAndGetChannel(t *testing.T) {
	messageBus := NewMessageBus()
	channelName := "test-channel"
	bufferSize := 5

	messageBus.RegisterChannel(channelName, bufferSize)
	channel, exists := messageBus.GetChannel(channelName)

	if !exists {
		t.Errorf("expected channel %s to exist, but it does not", channelName)
	}

	if channel == nil {
		t.Errorf("expected channel %s to be non-nil, but got nil", channelName)
	}

	select {
	case channel <- "test-message":
	default:
		t.Errorf("expected channel %s to have buffer size %d, but it is full", channelName, bufferSize)
	}
}

func TestMessageBus_GetNonExistentChannel(t *testing.T) {
	messageBus := NewMessageBus()
	nonExistentChannel := "non-existent-channel"

	channel, exists := messageBus.GetChannel(nonExistentChannel)

	if exists {
		t.Errorf("expected channel %s to not exist, but it does", nonExistentChannel)
	}

	if channel != nil {
		t.Errorf("expected channel %s to be nil, but got non-nil", nonExistentChannel)
	}
}
