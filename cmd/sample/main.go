package main

import (
	"context"
	"github.com/jobilla/go-app-events/pkg/appevents"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func main() {
	log.SetLevel(log.DebugLevel)
	l := &appevents.Listener{}
	err := l.Bootstrap("jobillaweb-staging", "app-events", "test-sub")

	if err != nil {
		panic(err)
	}

	g, _ := errgroup.WithContext(context.Background())
	g.Go(func () error {
		return l.Listen()
	})

	err = g.Wait()

	if err != nil {
		panic(err)
	}
}
