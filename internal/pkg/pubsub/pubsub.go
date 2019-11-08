package pubsub

import (
	"context"
	"errors"
	"gocloud.dev/pubsub"
	"os"
)

func provideDriver() string {
	return os.Getenv("PUBSUB_DRIVER")
}

func OpenTopic(ctx context.Context, topic string) (*pubsub.Topic, error) {
	switch provideDriver() {
	case "gcloud":
		return NewGCloudTopic(ctx, os.Getenv("GCP_PROJECT"), topic)
	case "rabbit":
		return NewRabbitTopic(ctx, topic)
	case "memory":
		return NewInMemoryTopic(ctx, topic)
	}

	return nil, errors.New("invalid driver supplied")
}
