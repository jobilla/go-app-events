package pubsub

import (
	"context"
	"fmt"
	"gocloud.dev/pubsub"
	"gocloud.dev/pubsub/mempubsub"
)

func NewInMemoryTopic(ctx context.Context, topic string) (*pubsub.Topic, error) {
	return pubsub.OpenTopic(ctx, fmt.Sprintf("mem://%s", topic))
}
