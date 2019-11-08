package pubsub

import (
	"context"
	"fmt"
	"gocloud.dev/pubsub"
	_ "gocloud.dev/pubsub/rabbitpubsub"
)

// Create a new topic on RabbitMQ. This requires the environment
// variable `RABBIT_SERVER_URL` to be set.
func NewRabbitTopic(ctx context.Context, topic string) (*pubsub.Topic, error) {
	return pubsub.OpenTopic(ctx, fmt.Sprintf("rabbit://%s", topic))
}
