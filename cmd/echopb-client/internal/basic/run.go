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
	Message  string `arg:"positional,required" help:"The message to send to the Echo service."`
	Count    int32  `arg:"--count" help:"How many times send the message." default:"1"`
	Interval int32  `arg:"--interval" help:"The interval in milliseconds between each message." default:"100"`
}

func Run(ctx context.Context, conn *shared.ConnectionPool, args Args) {
	ticker := time.NewTicker(time.Duration(args.Interval) * time.Millisecond)
	defer ticker.Stop()

	counter := int32(0)

	var wg sync.WaitGroup

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
				res, err := client.Echo(ctx, &echov1.EchoRequest{Message: args.Message})

				if err != nil {
					log.Printf("#%d failed to call Echo service: %s\n", counter+1, err)
				}

				log.Printf("#%d %s\n", counter+1, res.GetMessage())
			}(counter)
		}
	}
}
