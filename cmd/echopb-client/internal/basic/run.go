package basic

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
	Count    int32  `arg:"--count" help:"How many times send the message." default:"1"`
	Interval int32  `arg:"--interval" help:"The interval in milliseconds between each message." default:"100"`
}

func Run(ctx context.Context, conn *grpc.ClientConn, args Args) {
	client := echov1.NewEchoServiceClient(conn)

	for i := 0; i < int(args.Count); i++ {
		res, err := client.Echo(ctx, &echov1.EchoRequest{Message: args.Message})
		if err != nil {
			log.Fatalf("Failed to call Echo service: %v", err)
		}

		log.Printf("#%d %s\n", i+1, res.GetMessage())
		time.Sleep(time.Duration(args.Interval) * time.Millisecond)
	}
}
