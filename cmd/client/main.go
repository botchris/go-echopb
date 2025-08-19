package main

import (
	"bytes"
	"context"
	"log"

	"github.com/alexflint/go-arg"
	"github.com/botchris/go-echopb/cmd/client/internal/abort"
	"github.com/botchris/go-echopb/cmd/client/internal/basic"
	"github.com/botchris/go-echopb/cmd/client/internal/noop"
	"github.com/botchris/go-echopb/cmd/client/internal/ssbasic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type arguments struct {
	ServerAddr string `arg:"--host,required" help:"The address of the server"`
	Insecure   bool   `arg:"--insecure" help:"Use an insecure connection (without TLS)"`

	Basic *basic.Args `arg:"subcommand:basic" help:"Sends a message to the service and waits for a response."`
	Abort *abort.Args `arg:"subcommand:abort" help:"Sends back abort status."`
	Noop  *noop.Args  `arg:"subcommand:no-op" help:"Sends an empty request to the server amd waits for an empty response."`

	ServerStreamBasic *ssbasic.Args `arg:"subcommand:server-streaming-basic" help:"Sends a message to the service and waits for a stream of responses."`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	args := arguments{}
	parser := arg.MustParse(&args)

	dialOptions := make([]grpc.DialOption, 0)
	if args.Insecure {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	cc, err := grpc.NewClient(args.ServerAddr, dialOptions...)
	if err != nil {
		log.Fatal("Failed to create gRPC client:", err)
	}

	defer func() {
		if cErr := cc.Close(); cErr != nil {
			println("Failed to close connection:", cErr.Error())
		}
	}()

	switch {
	case args.Basic != nil:
		basic.Run(ctx, cc, *args.Basic)
	case args.Abort != nil:
		abort.Run(ctx, cc, *args.Abort)
	case args.Noop != nil:
		noop.Run(ctx, cc, *args.Noop)
	case args.ServerStreamBasic != nil:
		ssbasic.Run(ctx, cc, *args.ServerStreamBasic)
	default:
		var help bytes.Buffer

		parser.WriteHelp(&help)
		println(help.String())
	}
}
