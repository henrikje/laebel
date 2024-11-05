package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
	"github.com/r3labs/sse/v2"
	"log"
)

func PublishStatusUpdates(dockerClient *client.Client, server *sse.Server) {
	ctx := context.Background()
	eventCh, errCh := dockerClient.Events(ctx, events.ListOptions{})

	for {
		select {
		case event := <-eventCh:
			switch {
			case event.Type == events.ContainerEventType:
				switch event.Action {
				case events.ActionCreate, events.ActionDestroy,
					events.ActionStart, events.ActionPause, events.ActionUnPause, events.ActionStop, events.ActionDie, events.ActionHealthStatus, events.ActionHealthStatusHealthy, events.ActionHealthStatusRunning, events.ActionHealthStatusUnhealthy, events.ActionRestart, events.ActionRename:
					server.Publish("updates", &sse.Event{
						Event: []byte("reload"),
						Data:  eventData(event),
					})
				}
			}
		case err := <-errCh:
			log.Println("Error receiving Docker events.\nCause:", err)
		}
	}
}

func eventData(event events.Message) []byte {
	return []byte(fmt.Sprintf("%s:%s", event.Actor.ID, event.Action))
}
