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
				case events.ActionCreate, events.ActionDestroy:
					println("Received event:", event.Action)
					server.Publish("updates", &sse.Event{
						Event: []byte("main"),
						Data:  []byte("refresh"), // TODO Can we avoid having to send data?
					})
					log.Println("Published main event.")
				case events.ActionStart, events.ActionPause, events.ActionUnPause, events.ActionStop, events.ActionKill, events.ActionDie, events.ActionHealthStatus, events.ActionHealthStatusHealthy, events.ActionHealthStatusRunning, events.ActionHealthStatusUnhealthy, events.ActionRestart, events.ActionRename:
					println("Received event:", event.Action)
					serviceName := event.Actor.Attributes["com.docker.compose.service"]
					server.Publish("updates", &sse.Event{
						Event: []byte(serviceName),
						Data:  []byte("refresh"),
					})
					log.Println("Published service event.")
				}
			}
		case err := <-errCh:
			log.Println("Error receiving Docker events.\nCause:", err)
		}
	}
}
