package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"runtime/debug"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/testGRPC/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	PORT = ":9090"
)

// logging
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("before : %s, %v", info.FullMethod, req)
	resp, err := handler(ctx, req)
	log.Printf("after : %s, %v", info.FullMethod, resp)
	return resp, err
}

// recover
func RecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	defer func() {
		if e := recover(); e != nil {
			debug.PrintStack()
			err := status.Errorf(codes.Internal, "Panic err: %v", e)
			log.Printf("err: %v", err)
		}
	}()

	return handler(ctx, req)
}

func authorize(ctx context.Context) error {
	// warning: this is only for illustration purposes - don't implement authorization that is hardcoded!
	var authList = map[string]bool{
		base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", "scaramoucheX2", "Can-you-do-the-fandango?"))): true,
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.InvalidArgument, "retrieving metadata failed")
	}

	elem, ok := md["authorization"]
	if !ok {
		return status.Errorf(codes.InvalidArgument, "no auth details supplied")
	}

	authorization := elem[0][len("Basic "):] //extract base64 basic auth value (similar to HTTP Basic Auth)
	valid, ok := authList[authorization]
	if !ok {
		return status.Errorf(codes.NotFound, "auth not found")
	}

	if !valid {
		return status.Errorf(codes.Unauthenticated, "auth failed")
	}

	return nil
}

// auth
func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	// if info.FullMethod != "/proto.GoReleaseService/ListReleases" {
	// 	if err := authorize(ctx); err != nil {
	// 		return nil, err
	// 	}
	// }
	if err := authorize(ctx); err != nil {
		return nil, err
	}

	h, err := handler(ctx, req)

	//logging
	log.Printf("request - Method:%s\tDuration:%s\tError:%v\n",
		info.FullMethod,
		time.Since(start),
		err)

	return h, err

}

func main() {
	fmt.Println("vim-go")
	opts := []grpc.ServerOption{
		grpc_middleware.WithUnaryServerChain(
			LoggingInterceptor,
			RecoveryInterceptor,
			AuthInterceptor,
		),
	}

	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("failed to listen on port %s", PORT)
	}

	s := chat.Server{}
	rpcServer := grpc.NewServer(opts...)

	chat.RegisterChatServiceServer(rpcServer, &s)
	if err := rpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC server over port, err: %v", err)
	}

}
