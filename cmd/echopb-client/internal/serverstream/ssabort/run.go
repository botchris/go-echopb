package ssabort

import (
	"context"
	"errors"
	"io"
	"log"

	echov1 "github.com/botchris/go-echopb/gen/github.com/botchris/go-echopb/testing/echo/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Args defines the command line arguments for the echo subcommand.
type Args struct {
	Message  string `arg:"positional,required" help:"The message to send to the Echo service."`
	Count    int32  `arg:"--count,required" help:"The total number of messages to be generated before the server closes the stream."`
	Interval int32  `arg:"--interval" help:"The interval in milliseconds between each message sent by the server." default:"100"`
}

func Run(ctx context.Context, conn *grpc.ClientConn, args Args) {
	client := echov1.NewEchoServiceClient(conn)

	res, err := client.ServerStreamingEchoAbort(ctx, &echov1.ServerStreamingEchoRequest{
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
				log.Fatalf("Echo service returned EOF, but was an abort code")
			}

			st, ok := status.FromError(rErr)
			if !ok {
				log.Fatalf("Echo service returned an error, but it was not a gRPC status error: %v", rErr)
			}

			log.Printf("Response Status (%s): %s\n", st.Code().String(), st.Message())

			break
		}

		log.Printf("#%d %s", i, resp.GetMessage())

		i++
	}
}
