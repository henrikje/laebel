package main

import (
	"context"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
	"github.com/r3labs/sse/v2"
	"log"
)

func PublishStatusUpdates(dockerClient *client.Client, server *sse.Server) {
	ctx := context.Background()

	println("Registering for Docker events...")
	eventCh, errCh := dockerClient.Events(ctx, events.ListOptions{})

	for {
		select {
		case event := <-eventCh:
			switch {
			case event.Type == events.ContainerEventType:
				switch event.Action {
				case events.ActionCreate, events.ActionStart, events.ActionPause, events.ActionUnPause, events.ActionStop, events.ActionKill, events.ActionDie, events.ActionDestroy, events.ActionHealthStatus, events.ActionHealthStatusHealthy, events.ActionHealthStatusRunning, events.ActionHealthStatusUnhealthy, events.ActionRestart, events.ActionRename:
					server.Publish("refresh", &sse.Event{
						Data: []byte(event.Action),
					})
					log.Println("Published refresh event.")
				}
			}
		case err := <-errCh:
			log.Println("Error receiving Docker events.\nCause:", err)
		}
	}
}
