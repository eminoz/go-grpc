package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/eminoz/grpc-client/model"
	api "github.com/eminoz/grpc-client/proto"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var addr = flag.String("addr", "localhost:4040", "the address to connect to")

const (
	timestampFormat = time.StampNano // "Jan _2 15:04:05.000"
	streamingCount  = 100
)

func main() {
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client := api.NewEchoClient(conn)

	// g := gin.Default()
	app := fiber.New()

	app.Get("/creat", func(ct *fiber.Ctx) error {

		// Create metadata and context.
		md := metadata.Pairs("timestamp", time.Now().Format(timestampFormat))
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		// Make RPC using the context with the metadata.
		stream, err := client.ServerStreamingEcho(ctx, &api.EchoRequest{Message: "emin from server"})
		if err != nil {
			log.Fatalf("failed to call ServerStreamingEcho: %v", err)
		}

		// Read the header when the header arrives.
		header, err := stream.Header()
		if err != nil {
			log.Fatalf("failed to get header from stream: %v", err)
		}
		// Read metadata from server's header.
		if t, ok := header["timestamp"]; ok {
			fmt.Printf("timestamp from header:\n")
			for i, e := range t {
				fmt.Printf(" %d. %s\n", i, e)
			}
		} else {
			log.Fatal("timestamp expected but doesn't exist in header")
		}
		if l, ok := header["location"]; ok {
			fmt.Printf("location from header:\n")
			for i, e := range l {
				fmt.Printf(" %d. %s\n", i, e)
			}
		} else {
			log.Fatal("location expected but doesn't exist in header")
		}

		// Read all the responses.
		var rpcStatus error
		fmt.Printf("response:\n")
		for {
			r, err := stream.Recv()
			if err != nil {
				rpcStatus = err
				break
			}
			fmt.Printf(" - %s\n", r.Message)
			ct.JSON(r.Message)
		}
		if rpcStatus != io.EOF {
			log.Fatalf("failed to finish server streaming: %v", rpcStatus)
		}

		// Read the trailer after the RPC is finished.
		trailer := stream.Trailer()
		// Read metadata from server's trailer.
		if t, ok := trailer["timestamp"]; ok {
			fmt.Printf("timestamp from trailer:\n")
			for i, e := range t {
				fmt.Printf(" %d. %s\n", i, e)
			}
		} else {
			log.Fatal("timestamp expected but doesn't exist in trailer")
		}
		return nil
	})
	app.Post("/createPerson", func(c *fiber.Ctx) error {
		u := new(model.User)
		c.BodyParser(u)

		user := api.User{Username: u.Name, Lastname: u.Email}
		stream, err := client.ClientStreamingEcho(context.Background())
		if err != nil {
			return err
		}
		for i := 0; i < streamingCount; i++ {
			if err := stream.Send(&user); err != nil {
				log.Fatalf("failed to send streaming: %v\n", err)
			}
		} // Read the response.
		r, err := stream.CloseAndRecv()
		if err != nil {
			log.Fatalf("failed to CloseAndRecv: %v\n", err)
		}

		return c.JSON(r.Message)
	})
	app.Listen(":3131")
}
