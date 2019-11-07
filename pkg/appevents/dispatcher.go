package appevents

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/golang/protobuf/proto"
	app_pubsub "github.com/jobilla/go-app-events/internal/pkg/pubsub"
	"gocloud.dev/pubsub"
)

type Topic interface {
	Send(ctx context.Context, message *pubsub.Message) error
}

type Dispatcher struct {
	ctx   context.Context
	topic Topic
}

type Message interface {
	proto.Message
	// StringType should return a string representation of this
	// protobuf message, which will be used when sending app
	// events through this package.
	StringType() string
}

func (d *Dispatcher) Bootstrap(topicId string) error {
	topic, err := app_pubsub.OpenTopic(context.Background(), topic)
	d.topic = topic

	return err
}

func (d *Dispatcher) Dispatch(event string, message Message) error {
	m, err := proto.Marshal(message)

	if err != nil {
		return err
	}

	jsonString, err := json.Marshal(ProtoBody{
		Payload: m,
		Proto:   message.StringType(),
	})
	if err != nil {
		return err
	}

	body := []byte(base64.StdEncoding.EncodeToString(jsonString))

	return d.topic.Send(d.ctx, &pubsub.Message{
		Body: body,
		Metadata: map[string]string{
			"event": event,
		},
		BeforeSend: nil,
	})
}
