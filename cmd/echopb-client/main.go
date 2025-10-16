package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"log"

	"github.com/alexflint/go-arg"
	"github.com/botchris/go-echopb/cmd/echopb-client/internal/abort"
	"github.com/botchris/go-echopb/cmd/echopb-client/internal/basic"
	"github.com/botchris/go-echopb/cmd/echopb-client/internal/clientstream"
	"github.com/botchris/go-echopb/cmd/echopb-client/internal/fullduplex"
	"github.com/botchris/go-echopb/cmd/echopb-client/internal/noop"
	"github.com/botchris/go-echopb/cmd/echopb-client/internal/serverstream"
	"github.com/botchris/go-echopb/cmd/echopb-client/internal/shared"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type arguments struct {
	ServerAddr  string `arg:"--host,required" help:"The address of the server, e.g. example.com:443" placeholder:"HOST:POST"`
	Insecure    bool   `arg:"--insecure" help:"Use an insecure connection (without TLS)"`
	Connections uint   `arg:"--connections" help:"Number of connections to establish to the server" default:"1"`

	Basic *basic.Args `arg:"subcommand:basic" help:"Sends a message to the service and waits for a response."`
	Abort *abort.Args `arg:"subcommand:abort" help:"Sends back abort status."`
	Noop  *noop.Args  `arg:"subcommand:no-op" help:"Sends an empty request to the server amd waits for an empty response."`

	ServerStreamBasic     *serverstream.Args `arg:"subcommand:server-stream" help:"(Server Stream) Sends a message to the service and waits for a stream of responses from the server."`
	ClientStreamBasic     *clientstream.Args `arg:"subcommand:client-stream" help:"(Client Stream) Sends a stream of messages to the server, and then waits for the count response from the server."`
	FullDuplexStreamBasic *fullduplex.Args   `arg:"subcommand:full-duplex" help:"(Full Duplex Stream) Sends a stream of messages to the server, and receives them back in realtime."`
}

func main() {
	args := arguments{}
	parser := arg.MustParse(&args)

	if args.Connections == 0 {
		log.Fatal("The number of connections must be greater than 0")
	}

	dialOptions := make([]grpc.DialOption, 0)
	if args.Insecure {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	}

	connections := make([]*grpc.ClientConn, 0)

	for i := uint(0); i < args.Connections; i++ {
		cc, err := grpc.NewClient(args.ServerAddr, dialOptions...)
		if err != nil {
			log.Fatal("Failed to create gRPC client:", err)
		}

		connections = append(connections, cc)
	}

	pool := shared.NewConnectionPool(connections)
	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		cancel()

		for i := range connections {
			if cErr := connections[i].Close(); cErr != nil {
				println("Failed to close connection:", cErr.Error())
			}
		}
	}()

	switch {
	case args.Basic != nil:
		basic.Run(ctx, pool, *args.Basic)
	case args.Abort != nil:
		abort.Run(ctx, pool, *args.Abort)
	case args.Noop != nil:
		noop.Run(ctx, pool, *args.Noop)
	case args.ServerStreamBasic != nil:
		serverstream.Run(ctx, pool, *args.ServerStreamBasic)
	case args.ClientStreamBasic != nil:
		clientstream.Run(ctx, pool, *args.ClientStreamBasic)
	case args.FullDuplexStreamBasic != nil:
		fullduplex.Run(ctx, pool, *args.FullDuplexStreamBasic)
	default:
		var help bytes.Buffer

		parser.WriteHelp(&help)
		println(help.String())
	}
}
