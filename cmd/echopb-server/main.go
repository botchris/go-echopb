package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/botchris/go-echopb/cmd/echopb-server/internal/server"
	echov1 "github.com/botchris/go-echopb/gen/github.com/botchris/go-echopb/testing/echo/v1"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

type arguments struct {
	ListenAddr  string `arg:"--listen" help:"The address to listen on for incoming connections." default:":443" placeholder:"ADDR"`
	MetricsAddr string `arg:"--metrics-listen" help:"If provided, the address where to serve prometheus metrics under the HTTP path '/metrics'." placeholder:"ADDR"`
	Debug       bool   `arg:"--debug" help:"Enable debug logging."`
}

func main() {
	args := arguments{}
	arg.MustParse(&args)

	serverOptions := make([]grpc.ServerOption, 0)

	if args.MetricsAddr != "" {
		serverOptions = append(serverOptions,
			grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
			grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		)
	}

	grpcServer := grpc.NewServer(serverOptions...)
	echov1.RegisterEchoServiceServer(grpcServer, server.New())

	if args.Debug {
		grpclog.SetLoggerV2(grpclog.NewLoggerV2WithVerbosity(os.Stdout, os.Stdout, os.Stderr, 99))
	}

	if args.MetricsAddr != "" {
		grpc_prometheus.Register(grpcServer)
	}

	stopped := make(chan struct{})

	fx.New(
		fx.Invoke(func(lc fx.Lifecycle) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					grpcListener, err := (&net.ListenConfig{}).Listen(ctx, "tcp", args.ListenAddr)
					if err != nil {
						return err
					}

					var metricsListener net.Listener

					if args.MetricsAddr != "" {
						ml, mErr := (&net.ListenConfig{}).Listen(ctx, "tcp", args.MetricsAddr)
						if mErr != nil {
							return mErr
						}

						metricsListener = ml
					}

					go func() {
						defer func() {
							close(stopped)

							if cErr := grpcListener.Close(); cErr != nil {
								fmt.Printf("Failed to close grpcListener: %s\n", cErr)
							}
						}()

						fmt.Printf("gRPC server listening on %s\n", args.ListenAddr)

						// lock serving
						if sErr := grpcServer.Serve(grpcListener); sErr != nil {
							fmt.Printf("GRPC server ended with error: %s\n", err)
						}
					}()

					go func() {
						if args.MetricsAddr == "" {
							return
						}

						defer func() {
							if cErr := metricsListener.Close(); cErr != nil {
								fmt.Printf("Failed to close metricsListener: %s\n", cErr)
							}
						}()

						fmt.Printf("Metrics server serving on %s\n", args.MetricsAddr)

						http.Handle("/metrics", promhttp.Handler())
						if sErr := http.Serve(metricsListener, nil); sErr != nil {
							fmt.Printf("Metrics server ended with error: %s\n", sErr)
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
