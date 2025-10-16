package noop

import (
	"context"
	"log"

	"github.com/botchris/go-echopb/cmd/echopb-client/internal/shared"
	echov1 "github.com/botchris/go-echopb/gen/github.com/botchris/go-echopb/testing/echo/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Args defines the command line arguments for the echo subcommand.
type Args struct{}

func Run(ctx context.Context, conn *shared.ConnectionPool, args Args) {
	client := echov1.NewEchoServiceClient(conn.Next())

	_, err := client.NoOp(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("Failed to call Echo service: %v", err)
	}
}
