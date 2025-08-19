package basic

import (
	"context"
	"log"

	echov1 "github.com/botchris/go-echopb/gen/github.com/botchris/go-echopb/testing/echo/v1"
	"google.golang.org/grpc"
)

// Args defines the command line arguments for the echo subcommand.
type Args struct {
	Message string `arg:"positional,required" help:"The message to send to the Echo service."`
}

func Run(ctx context.Context, conn *grpc.ClientConn, args Args) {
	client := echov1.NewEchoServiceClient(conn)

	res, err := client.Echo(ctx, &echov1.EchoRequest{Message: args.Message})
	if err != nil {
		log.Fatalf("Failed to call Echo service: %v", err)
	}

	log.Printf("#%d %s\n", res.GetMessageCount(), res.GetMessage())
}
