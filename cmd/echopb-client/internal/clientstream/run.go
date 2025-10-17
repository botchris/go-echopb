package clientstream

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

	stream, err := client.ClientStreamingEcho(ctx)
	if err != nil {
		log.Fatalf("Failed to call Echo service: %v", err)
	}

	interval, dErr := time.ParseDuration(args.Interval)
	if dErr != nil {
		log.Fatalf("Failed to parse interval duration: %v", dErr)
	}

	for i := 0; i < int(args.Count); i++ {
		sErr := stream.Send(&echov1.ClientStreamingEchoRequest{Message: args.Message})
		if sErr != nil {
			log.Fatalf("Failed to send message to Echo service: %v", sErr)
		}

		log.Printf("#%d %s... sent!\n", i+1, args.Message)
		time.Sleep(interval)
	}

	println()
	log.Printf("All messages sent, waiting for response...\n")
	println()

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Failed to receive response from Echo service: %v", err)
	}

	log.Printf("Server receive count: %d\n", res.GetMessageCount())
}
