package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
	"github.com/r3labs/sse/v2"
	"log"
)

func PublishStatusUpdates(cli *client.Client, server *sse.Server) {
	ctx := context.Background()

	println("Registering for Docker events...")
	// TODO Specify the time we got the initial full response as "since"?
	eventCh, errCh := cli.Events(ctx, events.ListOptions{})

	for {
		select {
		case event := <-eventCh:
			if event.Type == events.ContainerEventType {
				fmt.Println("Received event, type:", event.Type, "action:", event.Action, "actorID:", event.Actor.ID)
				// TODO Generate proper event data
				data := fmt.Sprintf("%s: %s", event.Actor.ID, event.Action)
				server.Publish("status-updates", &sse.Event{
					Data: []byte(data),
				})
			}
		case err := <-errCh:
			log.Println("Error receiving Docker events:", err)
		}
	}
}
