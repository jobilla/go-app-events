package pubsub

import (
	"context"
	"fmt"
	"gocloud.dev/pubsub"
	_ "gocloud.dev/pubsub/gcppubsub"
)

func NewGCloudTopic(ctx context.Context, project string, topic string) (*pubsub.Topic, error) {
	return pubsub.OpenTopic(ctx, fmt.Sprintf("gcppubsub://projects/%s/topics/%s", project, topic))
}
