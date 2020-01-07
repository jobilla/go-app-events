package appevents

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type AppEventHandler func(payload []byte)
type ProtoBody struct {
	Payload []byte `json:"payload"`
	Proto   string `json:"proto"`
}

func (p *ProtoBody) FromPubsubMessage(message *pubsub.Message) error {
	return json.Unmarshal(message.Data, p)
}

type Listener struct {
	subscription *pubsub.Subscription
	ctx          context.Context
	handlers     map[string]AppEventHandler
}

func (l *Listener) Bootstrap(projectID string, subID string) error {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}

	subscription := client.Subscription(subID)
	l.subscription = subscription
	l.ctx = context.Background()

	return err
}

func (l *Listener) RegisterHandlers(handlers map[string]AppEventHandler) {
	if l.handlers == nil {
		l.handlers = handlers
	} else {
		for event, handler := range handlers {
			l.handlers[event] = handler
		}
	}
}

func (l *Listener) Listen() error {
	log.Debug("Listening for messages...")
	for {
		err := l.subscription.Receive(l.ctx, func(ctx context.Context, message *pubsub.Message) {
			log.
				WithContext(ctx).
				WithField("event_type", message.Attributes["event_type"]).
				Debug("Message received")

			if handler, ok := l.handlers[message.Attributes["event_type"]]; ok {
				body := &ProtoBody{}
				body.FromPubsubMessage(message)
				handler(body.Payload)

				log.
					WithContext(ctx).
					Debug("Message handled")
			}

			message.Ack()
		})

		if err != nil {
			return err
		}
	}
}
