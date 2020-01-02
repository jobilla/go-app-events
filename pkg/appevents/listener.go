package appevents

import (
	"context"
	"encoding/json"
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
}

func (l *Listener) Bootstrap() Listener {
	l.ctx = context.Background()

	return l
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
		message, err := l.subscription.Receive(l.ctx)

		if err != nil {
			return err
		}

		if handler, ok := l.handlers[message.Metadata["event"]]; ok {
			body := &ProtoBody{}
			err := body.FromPubsubMessage(message)

			if err != nil {
				return err
			}

			handler(body.Payload, message.Metadata["event"])
		}

		message.Ack()
	}
}
