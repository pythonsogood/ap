package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/pythonsogood/ap-assignment1/notification/cmd/notification/config"
	"github.com/pythonsogood/ap-assignment1/notification/internal/subscriber"
)

type logWriter struct {
	time_format string
}

func (lw *logWriter) Write(bs []byte) (int, error) {
	return fmt.Print(time.Now().UTC().Format(lw.time_format), " | ", string(bs))
}

func main() {
	conf, err := config.NewDefaultConfig()

	if err != nil {
		panic(err.Error())
	}

	log.SetFlags(0)
	log.SetOutput(&logWriter{
		time_format: time.RFC3339,
	})

	var event_subscriber subscriber.EventSubscriber

	switch conf.MessageBroker.Type {
	case config.MessageBrokerTypeNATS:
		nc, err := nats.Connect(conf.MessageBroker.Nats.ConnectionUrl)

		if err != nil {
			panic(err.Error())
		}

		event_subscriber = subscriber.NewNATSEventSubscriber(nc)
	default:
		panic("Unsupported message broker type!")
	}

	for _, subj := range conf.MessageBroker.LoggedSubjects {
		if err := event_subscriber.SubscribeNotification(subj); err != nil {
			panic(err.Error())
		}
	}
}
