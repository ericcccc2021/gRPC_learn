package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"playground/greet_s/greet_s"
	"strconv"
)

type server struct{}

func (s server) LongGreet(stream greet_s.GreetService_LongGreetServer) error {
	allRes := ""
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		c := res.GetInput()
		println("receive a new packet: c")
		allRes += c
	}
	stream.SendAndClose(&greet_s.LongGreetResponse{
		Output: allRes,
	})
	println("printing all concatenated res: " + allRes)
	return nil
}

func (s server) GreetManyTimes(request *greet_s.GreetManyTimesRequest, timesServer greet_s.GreetService_GreetManyTimesServer) error {
	//TODO implement me

	firstN := request.GetGreeting().GetLastName()
	lastN := request.GetGreeting().GetLastName()
	for i := 0; i < 10; i++ {
		result := "hello " + firstN + " " + lastN + " " + strconv.Itoa(i)
		res := &greet_s.GreetManyTimesResponse{
			Result: result,
		}
		err := timesServer.Send(res)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s server) Greet(ctx context.Context, request *greet_s.GreetRequest) (*greet_s.GreetResponse, error) {
	return &greet_s.GreetResponse{
		Result: "hello, this is from server side greetings",
	}, nil
}

func main() {
	fmt.Println("hello, gRPC greet_s server")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("fail to listten", err)
	}
	s := grpc.NewServer()
	greet_s.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatal("fail to serve:", err)
	}
}
