package noop

import (
	"context"
	"log"

	echov1 "github.com/botchris/go-echopb/gen/github.com/botchris/go-echopb/testing/echo/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Args defines the command line arguments for the echo subcommand.
type Args struct{}

func Run(ctx context.Context, conn *grpc.ClientConn, args Args) {
	client := echov1.NewEchoServiceClient(conn)

	_, err := client.NoOp(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("Failed to call Echo service: %v", err)
	}
}
