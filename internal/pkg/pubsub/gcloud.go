package pubsub

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gocloud.dev/gcp"
	"gocloud.dev/pubsub"
	"gocloud.dev/pubsub/gcppubsub"
	pubsub2 "google.golang.org/genproto/googleapis/pubsub/v1"
)

func NewGCloudTopic(ctx context.Context, project string, topic string) (*pubsub.Topic, error) {
	return pubsub.OpenTopic(ctx, fmt.Sprintf("gcppubsub://projects/%s/topics/%s", project, topic))
}

func NewGCloudSubscription(ctx context.Context, project string, topic string, subscription string) (*pubsub.Subscription, func (), error) {
	creds, err := gcp.DefaultCredentials(ctx)

	if err != nil {
		return nil, func() {}, err
	}

	conn, cleanup, err := gcppubsub.Dial(ctx, creds.TokenSource)

	if err != nil {
		return nil, cleanup, err
	}

	subClient, err := gcppubsub.SubscriberClient(ctx, conn)

	if err != nil {
		return nil, cleanup, err
	}

	_, err = subClient.GetSubscription(ctx, &pubsub2.GetSubscriptionRequest{Subscription: subscription})

	// sub doesn't exist
	if err != nil {
		log.WithField("subscription", subscription).Debug("subscription not found")
		_, err = subClient.CreateSubscription(ctx, &pubsub2.Subscription{
			Name:  fmt.Sprintf("projects/%s/subscriptions/%s", project, subscription),
			Topic: fmt.Sprintf("projects/%s/topics/%s", project, topic),
		})
		if err != nil {
			log.
				WithField("subscription", subscription).
				WithField("topic", topic).
				WithError(err).
				Error("failed to create subscription")

			return nil, cleanup, err
		}
		log.
			WithField("subscription", subscription).
			WithField("topic", topic).
			Debug("new subscription created")
	}


	return gcppubsub.OpenSubscription(
		subClient,
		gcp.ProjectID(project),
		subscription,
		nil,
	), cleanup, nil
}
