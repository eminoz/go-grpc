package controller

import (
	"context"
	"fmt"
	"io"
	"time"

	api "github.com/eminoz/grpc-api/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Strm struct{}

const (
	timestampFormat = time.StampNano
	streamingCount  = 10
)

// UnaryEcho is unary echo.
func (s Strm) UnaryEcho(_ context.Context, _ *api.EchoRequest) (*api.EchoResponse, error) {
	panic("not implemented") // TODO: Implement
}

// ServerStreamingEcho is server side streaming.
func (s Strm) ServerStreamingEcho(in *api.EchoRequest, stream api.Echo_ServerStreamingEchoServer) error {
	fmt.Print(in.Message)
	defer func() {
		trailer := metadata.Pairs("timestamp", time.Now().Format(timestampFormat))
		stream.SetTrailer(trailer)
	}()

	// Read metadata from client.
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return status.Errorf(codes.DataLoss, "ServerStreamingEcho: failed to get metadata")
	}
	if t, ok := md["timestamp"]; ok {
		fmt.Printf("timestamp from metadata:\n")
		for i, e := range t {
			fmt.Printf(" %d. %s\n", i, e)
		}
	}

	// Create and send header.
	header := metadata.New(map[string]string{"location": "MTV", "timestamp": time.Now().Format(timestampFormat)})
	stream.SendHeader(header)

	fmt.Printf("request received: %v\n", in)

	// Read requests and send responses.
	for i := 0; i < streamingCount; i++ {
		fmt.Printf("echo message %v\n", in.Message)
		err := stream.Send(&api.EchoResponse{Message: in.Message})
		if err != nil {
			return err
		}
	}
	return nil
}

// ClientStreamingEcho is client side streaming.
func (s Strm) ClientStreamingEcho(stream api.Echo_ClientStreamingEchoServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("echo last received message\n")
			return stream.SendAndClose(&api.EchoResponse{Message: "all user saved"})
		}
		fmt.Println(in.Username + " " + in.Lastname)
		if err != nil {
			return err
		}
	}
}

// BidirectionalStreamingEcho is bidi streaming.
func (s Strm) BidirectionalStreamingEcho(_ api.Echo_BidirectionalStreamingEchoServer) error {
	panic("not implemented") // TODO: Implement
}
