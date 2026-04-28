package main

import (
	"github.com/nats-io/nats.go"
	"github.com/pythonsogood/ap-assignment1/notification/cmd/notification/config"
	"github.com/pythonsogood/ap-assignment1/notification/internal/subscriber"
)

func main() {
	conf, err := config.NewDefaultConfig()

	if err != nil {
		panic(err.Error())
	}

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
