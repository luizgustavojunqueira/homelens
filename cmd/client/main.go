package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"homelens/client"
	"homelens/shared"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	agentClient := client.NewAgentClient(log.Printf)
	conn := agentClient.Connect("localhost:6969")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	out := make(chan shared.SystemInfo)
	go client.Collect(ctx, time.Second, out)

	for {
		select {
		case <-sigChan:
			fmt.Println("Shutting down...")
			cancel()
			agentClient.Disconnect(conn)
			return
		case snapshot := <-out:
			if err := agentClient.SendSnapshot(conn, snapshot); err != nil {
				fmt.Println("send error:", err)
				return
			}
		}
	}
}
