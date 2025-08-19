package abort

import (
	"context"
	"log"

	echov1 "github.com/botchris/go-echopb/gen/github.com/botchris/go-echopb/testing/echo/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Args defines the command line arguments for the echo subcommand.
type Args struct {
	Message string `arg:"positional,required" help:"The message to send to the Echo service."`
}

func Run(ctx context.Context, conn *grpc.ClientConn, args Args) {
	client := echov1.NewEchoServiceClient(conn)

	res, err := client.EchoAbort(ctx, &echov1.EchoRequest{Message: args.Message})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			log.Fatalf("Failed to call EchoAbort service: %v", err)
		}

		log.Printf("Response Status (%d): %s\n", st.Code(), st.Message())
	}

	log.Printf("#%d %s\n", res.GetMessageCount(), res.GetMessage())
}
