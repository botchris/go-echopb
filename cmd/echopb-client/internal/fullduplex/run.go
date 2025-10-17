package fullduplex

import (
	"context"
	"log"
	"time"

	"github.com/botchris/go-echopb/cmd/echopb-client/internal/shared"
	echov1 "github.com/botchris/go-echopb/gen/github.com/botchris/go-echopb/testing/echo/v1"
)

// Args defines the command line arguments for the echo subcommand.
type Args struct {
	Message  string `arg:"positional,required" help:"The message to send to the Echo service."`
	Count    int32  `arg:"--count,required" help:"The total number of messages to be generated before the server closes the stream."`
	Interval string `arg:"--interval" help:"The interval between each message sent by the server. Must be a valid duration string (e.g., '100ms', '2s', '1m')." default:"100ms"`
}

// Run executes the subcommand.
func Run(ctx context.Context, conn *shared.ConnectionPool, args Args) {
	client := echov1.NewEchoServiceClient(conn.Next())

	bidi, err := client.FullDuplexEcho(ctx)
	if err != nil {
		log.Fatalf("Failed to call Echo service: %v", err)
	}

	interval, dErr := time.ParseDuration(args.Interval)
	if dErr != nil {
		log.Fatalf("Failed to parse interval duration: %v", dErr)
	}

	go func() {
		for {
			resp, rErr := bidi.Recv()
			if rErr != nil {
				return
			}

			log.Printf("Rcv: %s\n", resp.GetMessage())
		}
	}()

	for i := 0; i < int(args.Count); i++ {
		if sErr := bidi.Send(&echov1.EchoRequest{Message: args.Message}); sErr != nil {
			log.Fatalf("Failed to send message to Echo service: %v", sErr)
		}

		log.Printf("Sent: [#%d] %s\n", i+1, args.Message)

		if i < int(args.Count)-1 {
			time.Sleep(interval)
		}
	}
}
