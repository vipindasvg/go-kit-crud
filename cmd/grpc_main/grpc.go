package main

import (
	"net"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitoc "github.com/go-kit/kit/tracing/opencensus"
	kitgrpc "github.com/go-kit/kit/transport/grpc"

	"google.golang.org/grpc"
	usersvc "github.com/vipindasvg/go-kit-crud/user"
	pb "github.com/vipindasvg/go-kit-crud/user/pb"
)

const (
	grpcport     = ":50051"
)

func main() {
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger,
			"svc", "user",
			"ts", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	// Create Order Service
	var svc usersvc.Service
	{
		var err error
		repository := usersvc.NewRepo(logger)
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
		svc = usersvc.NewService(repository, logger)
	}
	// Create Go kit endpoints for the Order Service
	var endpoints usersvc.Endpoints
	{
		fmt.Println("endpoints")
		endpoints = usersvc.MakeEndpoints(svc)
	}
	// set-up grpc transport
	var (
		ocTracingGrpc   = kitoc.GRPCServerTrace()
		serverOptions   = []kitgrpc.ServerOption{ocTracingGrpc}
		userService  	= usersvc.NewGRPCServer(endpoints, serverOptions, logger)
		grpcListener, _ = net.Listen("tcp", grpcport)
		grpcServer      = grpc.NewServer()
	)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "gRPC", "addr", grpcport)
		pb.RegisterUserServiceServer(grpcServer, userService)
		errs <- grpcServer.Serve(grpcListener)
	}()

	level.Error(logger).Log("exit", <-errs)
}
