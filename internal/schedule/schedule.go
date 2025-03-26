package schedule

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/mehmetalisavas/message-sender/internal/pubsub"
	"github.com/mehmetalisavas/message-sender/internal/service"
)

type Schedule struct {
	storageService service.Storage
	messageBus     *pubsub.MessageBus
	producers      []pubsub.Producer
	consumers      []pubsub.Consumer
	wg             sync.WaitGroup
}

// NewScheduler creates a new Schedule instance.
func NewScheduler(storageService service.Storage) *Schedule {
	bus := pubsub.NewMessageBus()
	bus.RegisterChannel(pubsub.MessageSenderTopic, 2)
	return &Schedule{
		storageService: storageService,
		messageBus:     bus,
		producers:      make([]pubsub.Producer, 0),
		consumers:      make([]pubsub.Consumer, 0),
	}
}

func (s *Schedule) AddProducer(producer pubsub.Producer) {
	s.producers = append(s.producers, producer)
}

func (s *Schedule) AddConsumer(consumer pubsub.Consumer) {
	s.consumers = append(s.consumers, consumer)
}

func (s *Schedule) Start(ctx context.Context, workerCount int) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // this is just a fallback in case the function exits earlier due to an error.

	for _, producer := range s.producers {
		s.wg.Add(1)
		go func(p pubsub.Producer) {
			defer s.wg.Done()
			p.Produce(ctx)
		}(producer)
	}

	for _, consumer := range s.consumers {
		s.wg.Add(1)
		go func(c pubsub.Consumer) {
			defer s.wg.Done()
			c.Consume(ctx, workerCount)
		}(consumer)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	cancel()    // cancel the context before doing cleanup
	s.wg.Wait() // wait for all producers and consumers to finish
}

func (s *Schedule) MessageBus() *pubsub.MessageBus {
	return s.messageBus
}

func (s *Schedule) RegisterChannelToMessageBus(topic string, workerCount int) {
	s.messageBus.RegisterChannel(topic, workerCount)
}
