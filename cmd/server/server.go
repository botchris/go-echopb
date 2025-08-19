package main

import (
	"context"
	"time"

	echov1 "github.com/botchris/go-echopb/gen/github.com/botchris/go-echopb/testing/echo/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	echov1.UnimplementedEchoServiceServer
}

// NewServer returns a new EchoServiceServer implementation.
func NewServer() echov1.EchoServiceServer {
	return &server{}
}

func (a *server) Echo(_ context.Context, req *echov1.EchoRequest) (*echov1.EchoResponse, error) {
	return &echov1.EchoResponse{
		Message: req.Message,
	}, nil
}

func (a *server) EchoAbort(_ context.Context, req *echov1.EchoRequest) (*echov1.EchoResponse, error) {
	return &echov1.EchoResponse{
		Message: req.Message,
	}, status.Error(codes.Aborted, req.Message)
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

			if err := server.Send(&echov1.ServerStreamingEchoResponse{
				Message: req.Message,
			}); err != nil {
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
			return client.SendAndClose(&echov1.ClientStreamingEchoResponse{
				MessageCount: int32(count),
			})
		default:
			if _, err := client.Recv(); err != nil {
				return err
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

func (a *server) HalfDuplexEcho(server echov1.EchoService_HalfDuplexEchoServer) error {
	ctx := server.Context()
	done := ctx.Done()
	count := 0
	buffer := make([]*echov1.EchoRequest, 0)

	for {
		select {
		case <-done:
			return nil
		default:
			req, err := server.Recv()
			if err != nil {
				for _, bReq := range buffer {
					if sErr := server.Send(&echov1.EchoResponse{
						Message:      bReq.Message,
						MessageCount: int32(count),
					}); sErr != nil {
						return sErr
					}
				}

				return nil
			}

			buffer = append(buffer, req)
		}
	}
}
