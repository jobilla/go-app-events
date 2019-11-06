package appevents

import "testing"

func TestListener_Listen(t *testing.T) {
	listener := &Listener{
		subscription: nil,
		ctx:          nil,
		handlers:     nil,
	}

	handler := func(payload []byte, event string) {}

	listener.RegisterHandlers(map[string]AppEventHandler{
		"event.test": handler,
	})
}
