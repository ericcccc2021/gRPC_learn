package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"playground/greet_s/greet_s"
	"strconv"
)

func main() {
	fmt.Println("hello, greet gPRC client")
	cnn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cannot connect %v", err)
	}

	defer cnn.Close()

	c := greet_s.NewGreetServiceClient(cnn)
	doServiceStreaming(c)
	doClientStreaming(c)

}

func doClientStreaming(c greet_s.GreetServiceClient) {
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("cannot send stream to server", err)
	}
	r := []*greet_s.LongGreetRequest{}
	for i := 0; i < 10; i++ {
		r = append(r, &greet_s.LongGreetRequest{
			Input: "this is message" + strconv.Itoa(i),
		})
	}
	for _, req := range r {
		err := stream.Send(req)
		if err != nil {
			log.Println("error when sending streams")
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error when receiving response from server", err)

	}
	fmt.Println("this is what I received from LongGreet: %v", res)
}

func doServiceStreaming(c greet_s.GreetServiceClient) {
	stream, err := c.GreetManyTimes(context.Background(), &greet_s.GreetManyTimesRequest{
		Greeting: &greet_s.Greeting{
			FirstName: "hanxiong",
			LastName:  "zhang",
		},
	})
	if err != nil {
		log.Fatalf("cannot fetch from server", err)
	}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("message over")
			break
		}
		if err != nil {
			log.Fatalf("error reading stream", err)
		}
		log.Println("log response: " + msg.GetResult())
	}
}
