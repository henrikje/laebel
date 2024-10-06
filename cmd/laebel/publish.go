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
				case events.ActionCreate, events.ActionDestroy:
					server.Publish("updates", &sse.Event{
						Event: []byte("reload"),
						Data:  []byte(fmt.Sprintf("%s:%s", event.Actor.ID, event.Action)),
					})
					server.Publish("updates", &sse.Event{
						Event: []byte("close"),
						Data:  []byte("No further updates will be sent."),
					})

				case events.ActionStart, events.ActionPause, events.ActionUnPause, events.ActionStop, events.ActionKill, events.ActionDie, events.ActionHealthStatus, events.ActionHealthStatusHealthy, events.ActionHealthStatusRunning, events.ActionHealthStatusUnhealthy, events.ActionRestart, events.ActionRename:
					log.Println(server)
					server.Publish("updates", &sse.Event{
						Event: []byte("status"),
						Data:  []byte(fmt.Sprintf("%s:%s", event.Actor.ID, event.Action)),
					})
				}
			}
		case err := <-errCh:
			log.Println("Error receiving Docker events.\nCause:", err)
		}
	}
}
