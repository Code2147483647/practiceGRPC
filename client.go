package main

import (
	"log"

	"github.com/testGRPC/chat"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func main() {
	grpcAuth := &BasicAuthCreds{
		username: "scaramoucheX2",
		password: "Can-you-do-the-fandango?",
	}

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9090",
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(grpcAuth),
	)

	if err != nil {
		log.Fatalf("failed to connect, %v", err)
	}
	defer conn.Close()

	client := chat.NewChatServiceClient(conn)
	message := chat.Message{
		Body:  "hello from client",
		Count: &wrapperspb.Int32Value{Value: 2},
	}

	response, err := client.SayHello(context.Background(), &message)
	if err != nil {
		log.Fatalf("error when calling SayHello: %v", err)
	}

	log.Printf("Response from server: %s, %d", response.Body, response.GetCount())
}
