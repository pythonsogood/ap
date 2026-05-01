package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/pythonsogood/ap-assignment1/notification/cmd/notification/config"
	"github.com/pythonsogood/ap-assignment1/notification/internal/subscriber"
)

type logWriter struct {
	time_format string
	file_path   string
}

func (lw *logWriter) Write(bs []byte) (int, error) {
	text := fmt.Sprint(time.Now().UTC().Format(lw.time_format), " | ", string(bs))

	f, err := os.OpenFile(lw.file_path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	if err == nil {
		defer f.Close()

		f.WriteString(text)
	}

	return fmt.Print(text)
}

func main() {
	conf, err := config.NewDefaultConfig()

	if err != nil {
		panic(err.Error())
	}

	log.SetFlags(0)
	log.SetOutput(&logWriter{
		time_format: time.RFC3339,
		file_path:   conf.Log.File,
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

	for {
	}
}
