package clientstream

import (
	"context"
	"log"
	"time"

	echov1 "github.com/botchris/go-echopb/gen/github.com/botchris/go-echopb/testing/echo/v1"
	"google.golang.org/grpc"
)

// Args defines the command line arguments for the echo subcommand.
type Args struct {
	Message  string `arg:"positional,required" help:"The message to send to the Echo service."`
	Count    int32  `arg:"--count,required" help:"The total number of messages to be generated before the server closes the stream."`
	Interval int32  `arg:"--interval" help:"The interval in milliseconds between each message sent by the server." default:"100"`
}

func Run(ctx context.Context, conn *grpc.ClientConn, args Args) {
	client := echov1.NewEchoServiceClient(conn)

	stream, err := client.ClientStreamingEcho(ctx)
	if err != nil {
		log.Fatalf("Failed to call Echo service: %v", err)
	}

	for i := 0; i < int(args.Count); i++ {
		sErr := stream.Send(&echov1.ClientStreamingEchoRequest{Message: args.Message})
		if sErr != nil {
			log.Fatalf("Failed to send message to Echo service: %v", sErr)
		}

		log.Printf("#%d %s... sent!\n", i+1, args.Message)
		time.Sleep(time.Duration(args.Interval) * time.Millisecond)
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
