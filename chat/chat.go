package chat

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Server struct {
}

// SayHello test error situation
func (s *Server) SayHello(ctx context.Context, message *Message) (*Message, error) {
	log.Printf("Received message body from client: %s, count: %v", message.Body, message.Count)
	return &Message{Body: "hello back", Count: &wrapperspb.Int32Value{Value: 1}}, status.Error(codes.NotFound, "qqqqq!")
}
