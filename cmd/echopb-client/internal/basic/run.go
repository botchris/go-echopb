package basic

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/botchris/go-echopb/cmd/echopb-client/internal/shared"
	echov1 "github.com/botchris/go-echopb/gen/github.com/botchris/go-echopb/testing/echo/v1"
)

// Args defines the command line arguments for the echo subcommand.
type Args struct {
	Message  string `arg:"positional,required" help:"The message to send to the Echo service. Use the pattern '@lorem(<number>)' to generate a random lorem ipsum sentence of the specified words count."`
	Count    int32  `arg:"--count" help:"How many times send the message." default:"1"`
	Interval string `arg:"--interval" help:"The interval between each message. Must be a valid duration string (e.g., '100ms', '2s', '1m')." default:"100ms"`
}

// Run executes the subcommand.
func Run(ctx context.Context, conn *shared.ConnectionPool, args Args) {
	messenger, gErr := shared.NewMessageGenerator(args.Message)
	if gErr != nil {
		log.Fatalf("Failed to create message generator: %v", gErr)
	}

	interval, dErr := time.ParseDuration(args.Interval)
	if dErr != nil {
		log.Fatalf("Failed to parse interval duration: %v", dErr)
	}

	if interval < time.Millisecond {
		interval = time.Millisecond
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	counter := int32(0)
	wg := sync.WaitGroup{}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if counter >= args.Count {
				wg.Wait()

				return
			}

			counter++

			wg.Add(1)

			go func(iteration int32) {
				defer wg.Done()

				client := echov1.NewEchoServiceClient(conn.Next())
				res, err := client.Echo(ctx, &echov1.EchoRequest{Message: messenger.Get()})

				if err != nil {
					log.Printf("#%d failed to call Echo service: %s\n", iteration+1, err)
				}

				log.Printf("#%d %s\n", iteration+1, res.GetMessage())
			}(counter)
		}
	}
}
