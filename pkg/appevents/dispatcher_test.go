package appevents

import (
	"context"
	Test "github.com/jobilla/go-app-events/test"
	"gocloud.dev/pubsub"
	"testing"
)

type TestTopic struct {
	t     *testing.T
	event string
}

func (t *TestTopic) Send(ctx context.Context, message *pubsub.Message) error {
	if message.Metadata["event_type"] != t.event {
		t.t.Errorf("expected event to be %s, got %s", t.event, message.Metadata["event"])
	}

	body := &ProtoBody{}
	err := body.FromPubsubMessage(message)

	if err != nil {
		return err
	}

	if body.Proto != "test" {
		t.t.Errorf("expected proto type to be test, got: %s", body.Proto)
	}

	return nil
}

func TestDispatcher_Dispatch(t *testing.T) {
	d := &Dispatcher{
		ctx: context.Background(),
		topic: &TestTopic{
			t:     t,
			event: "test.event",
		},
	}

	err := d.Dispatch("test.event", &Test.Test{
		Value: "foo",
	})

	if err != nil {
		t.Errorf("expected error to be nil, got: %s", err.Error())
	}
}
