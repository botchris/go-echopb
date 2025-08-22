package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/alexflint/go-arg"
	echov1 "github.com/botchris/go-echopb/gen/github.com/botchris/go-echopb/testing/echo/v1"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type arguments struct {
	ListenAddr string `arg:"--listen" help:"The address to listen on for incoming connections." default:":443" placeholder:"ADDR"`
	Debug      bool   `arg:"--debug" help:"Enable debug logging."`
}

func main() {
	args := arguments{}
	arg.MustParse(&args)

	options := make([]grpc.ServerOption, 0)
	grpcServer := grpc.NewServer(options...)
	echov1.RegisterEchoServiceServer(grpcServer, NewServer())

	if args.Debug {
		grpclog.SetLoggerV2(grpclog.NewLoggerV2WithVerbosity(os.Stdout, os.Stdout, os.Stderr, 99))
	}

	stopped := make(chan struct{})

	fx.New(
		fx.Invoke(func(lc fx.Lifecycle) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					listener, err := net.Listen("tcp", args.ListenAddr)
					if err != nil {
						return err
					}

					go func() {
						defer close(stopped)

						// lock serving
						if sErr := grpcServer.Serve(listener); sErr != nil {
							fmt.Printf("GRPC server ended with error: %s\n", err)
						}
					}()

					return nil
				},
				OnStop: func(context.Context) error {
					grpcServer.GracefulStop()

					<-stopped

					return nil
				},
			})
		}),
	).
		Run()
}

type server struct {
	echov1.UnimplementedEchoServiceServer
}

// NewServer returns a new EchoServiceServer implementation.
func NewServer() echov1.EchoServiceServer {
	return &server{}
}

func (a *server) Echo(_ context.Context, req *echov1.EchoRequest) (*echov1.EchoResponse, error) {
	return &echov1.EchoResponse{Message: req.Message}, nil
}

func (a *server) EchoAbort(_ context.Context, req *echov1.EchoRequest) (*echov1.EchoResponse, error) {
	return &echov1.EchoResponse{Message: req.Message}, status.Error(codes.Aborted, req.Message)
}

func (a *server) NoOp(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (a *server) ServerStreamingEcho(req *echov1.ServerStreamingEchoRequest, server echov1.EchoService_ServerStreamingEchoServer) error {
	ctx := server.Context()
	done := ctx.Done()

	interval := time.Duration(req.MessageInterval) * time.Millisecond
	if interval < 100*time.Millisecond {
		interval = 100 * time.Millisecond
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	sent := 0

	for {
		select {
		case <-done:
			return nil
		case <-ticker.C:
			if sent >= int(req.MessageCount) {
				return nil
			}

			if err := server.Send(&echov1.ServerStreamingEchoResponse{Message: req.Message}); err != nil {
				return err
			}

			sent++
		}
	}
}

func (a *server) ServerStreamingEchoAbort(req *echov1.ServerStreamingEchoRequest, server echov1.EchoService_ServerStreamingEchoAbortServer) error {
	ctx := server.Context()
	done := ctx.Done()

	interval := time.Duration(req.MessageInterval) * time.Millisecond
	if interval < 100*time.Millisecond {
		interval = 100 * time.Millisecond
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	sent := 0

	for {
		select {
		case <-done:
			return status.Error(codes.Aborted, req.Message)
		case <-ticker.C:
			if sent >= int(req.MessageCount) {
				return status.Error(codes.Aborted, req.Message)
			}

			if err := server.Send(&echov1.ServerStreamingEchoResponse{
				Message: req.Message,
			}); err != nil {
				return status.Error(codes.Aborted, err.Error())
			}

			sent++
		}
	}
}

func (a *server) ClientStreamingEcho(client echov1.EchoService_ClientStreamingEchoServer) error {
	ctx := client.Context()
	done := ctx.Done()
	count := 0

	for {
		select {
		case <-done:
			return client.SendAndClose(&echov1.ClientStreamingEchoResponse{MessageCount: int32(count)})
		default:
			_, err := client.Recv()
			if err != nil {
				if !errors.Is(err, io.EOF) {
					return err
				}

				return client.SendAndClose(&echov1.ClientStreamingEchoResponse{MessageCount: int32(count)})
			}

			count++
		}
	}
}

func (a *server) FullDuplexEcho(server echov1.EchoService_FullDuplexEchoServer) error {
	ctx := server.Context()
	done := ctx.Done()
	count := 0

	for {
		select {
		case <-done:
			return nil
		default:
			req, err := server.Recv()
			if err != nil {
				return err
			}

			count++

			if sErr := server.Send(&echov1.EchoResponse{
				Message:      req.Message,
				MessageCount: int32(count),
			}); sErr != nil {
				return sErr
			}
		}
	}
}
