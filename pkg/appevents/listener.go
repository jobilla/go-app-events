package appevents

import (
	"context"
	"encoding/json"
	pubsub2 "github.com/jobilla/go-app-events/internal/pkg/pubsub"
	log "github.com/sirupsen/logrus"
	"gocloud.dev/pubsub"
)

type AppEventHandler func(payload []byte, event string)
type ProtoBody struct {
	Payload []byte `json:"payload"`
	Proto   string `json:"proto"`
}

func (p *ProtoBody) FromPubsubMessage(message *pubsub.Message) error {
	return json.Unmarshal(message.Body, p)
}

type Listener struct {
	subscription *pubsub.Subscription
	ctx          context.Context
	handlers     map[string]AppEventHandler
	cleanup      func()
}

func (l *Listener) Bootstrap(projectID string, topicID string, subID string) error {
	l.ctx = context.Background()
	var err error

	l.subscription, l.cleanup, err = pubsub2.OpenSubscription(l.ctx, subID, topicID)

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
	for {
		log.Debug("receiving message")
		message, err := l.subscription.Receive(l.ctx)

		if err != nil {
			return err
		}
		log.WithField("message", string(message.Body)).Debug("received message")

		if handler, ok := l.handlers[message.Metadata["event"]]; ok {
			body := &ProtoBody{}
			err = body.FromPubsubMessage(message)

			handler(body.Payload, message.Metadata["event"])
		}

		if err != nil {
			return err
		}

		message.Ack()
	}
}
