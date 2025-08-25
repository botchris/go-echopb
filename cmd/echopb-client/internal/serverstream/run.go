package serverstream

import (
	"context"
	"errors"
	"io"
	"log"

	"github.com/botchris/go-echopb/cmd/echopb-client/internal/serverstream/ssabort"
	echov1 "github.com/botchris/go-echopb/gen/github.com/botchris/go-echopb/testing/echo/v1"
	"google.golang.org/grpc"
)

// Args defines the command line arguments for the echo subcommand.
type Args struct {
	Message  string `arg:"positional,required" help:"The message to send to the Echo service."`
	Count    int32  `arg:"--count,required" help:"The total number of messages to be generated before the server closes the stream."`
	Interval int32  `arg:"--interval" help:"The interval in milliseconds between each message sent by the server." default:"100"`
	Abort    bool   `arg:"--abort" help:"Indicates the server to send an abort status when finishing the connection"`
}

func Run(ctx context.Context, conn *grpc.ClientConn, args Args) {
	if args.Abort {
		ssabort.Run(ctx, conn, ssabort.Args{
			Message:  args.Message,
			Count:    args.Count,
			Interval: args.Interval,
		})

		return
	}

	client := echov1.NewEchoServiceClient(conn)

	res, err := client.ServerStreamingEcho(ctx, &echov1.ServerStreamingEchoRequest{
		Message:         args.Message,
		MessageCount:    args.Count,
		MessageInterval: args.Interval,
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
