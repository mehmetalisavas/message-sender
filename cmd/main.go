package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/mehmetalisavas/message-sender/config"
	_ "github.com/mehmetalisavas/message-sender/docs"
	"github.com/mehmetalisavas/message-sender/internal/api"
	"github.com/mehmetalisavas/message-sender/internal/db/mysql"
	"github.com/mehmetalisavas/message-sender/internal/db/redis"
	"github.com/mehmetalisavas/message-sender/internal/pubsub"
	"github.com/mehmetalisavas/message-sender/internal/route"
	"github.com/mehmetalisavas/message-sender/internal/schedule"
	"github.com/mehmetalisavas/message-sender/pkg/services/notification"

	"github.com/sethvargo/go-envconfig"
)

const defaultRequestTimeout = 10  // seconds
const defaultTickerInterval = 120 // seconds

func main() {
	ctx := context.Background()

	c := config.New()
	if err := envconfig.Process(ctx, &c); err != nil {
		log.Fatal(err)
	}
	mysqlDB, err := mysql.NewClient(c)
	if err != nil {
		log.Fatal(err)
	}
	defer mysqlDB.Close()

	sqlStorage := mysql.NewSqlStore(mysqlDB)

	cacheService, err := redis.NewRedisCacheStore(ctx, c)
	if err != nil {
		log.Fatalf("error while starting cache service: %v \n", err)
	}

	notificationService := notification.NewNotificationService(c.NotificationServiceURL, time.Duration(defaultRequestTimeout)*time.Second)

	scheduler := schedule.NewScheduler(sqlStorage)
	messageProducer := pubsub.NewMessageProducer(&c, sqlStorage, scheduler.MessageBus(), defaultTickerInterval)
	scheduler.AddProducer(messageProducer)
	messageConsumer := pubsub.NewMessageConsumer(sqlStorage, scheduler.MessageBus(), notificationService, cacheService)
	scheduler.AddConsumer(messageConsumer)

	go scheduler.Start(ctx, 2) // start with 2 workers

	api := api.New(&c, sqlStorage)

	routers := route.Routers(api)

	log.Printf("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", routers))

}
