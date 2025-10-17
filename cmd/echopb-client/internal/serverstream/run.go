package serverstream

import (
	"context"
	"errors"
	"io"
	"log"
	"time"

	"github.com/botchris/go-echopb/cmd/echopb-client/internal/serverstream/ssabort"
	"github.com/botchris/go-echopb/cmd/echopb-client/internal/shared"
	echov1 "github.com/botchris/go-echopb/gen/github.com/botchris/go-echopb/testing/echo/v1"
)

// Args defines the command line arguments for the echo subcommand.
type Args struct {
	Message  string `arg:"positional,required" help:"The message to send to the Echo service."`
	Count    int32  `arg:"--count,required" help:"The total number of messages to be generated before the server closes the stream."`
	Interval string `arg:"--interval" help:"The interval between each message sent by the server. Must be a valid duration string (e.g., '100ms', '2s', '1m')." default:"100ms"`
	Abort    bool   `arg:"--abort" help:"Indicates the server to send an abort status when finishing the connection"`
}

// Run executes the subcommand.
func Run(ctx context.Context, conn *shared.ConnectionPool, args Args) {
	if args.Abort {
		ssabort.Run(ctx, conn, ssabort.Args{
			Message:  args.Message,
			Count:    args.Count,
			Interval: args.Interval,
		})

		return
	}

	interval, dErr := time.ParseDuration(args.Interval)
	if dErr != nil {
		log.Fatalf("Failed to parse interval duration: %v", dErr)
	}

	if interval < time.Millisecond {
		interval = time.Millisecond
	}

	client := echov1.NewEchoServiceClient(conn.Next())

	res, err := client.ServerStreamingEcho(ctx, &echov1.ServerStreamingEchoRequest{
		Message:         args.Message,
		MessageCount:    args.Count,
		MessageInterval: int32(interval.Milliseconds()),
	})

	if err != nil {
		log.Fatalf("Failed to call Echo service: %v", err)
	}

	i := 1

	for {
		resp, rErr := res.Recv()
		if rErr != nil {
			if errors.Is(rErr, io.EOF) {
				break
			}

			log.Fatalf("Failed to receive response: %v", rErr)
		}

		log.Printf("#%d %s", i, resp.GetMessage())

		i++
	}
}
