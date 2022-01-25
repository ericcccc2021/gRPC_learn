package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"math"
	"net"
	"playground/greet_s/greet_s"
	"strconv"
	"time"
)

type server struct{}

func (s server) GreetWithDeadLine(ctx context.Context, request *greet_s.GreetWithDeadLineRequest) (*greet_s.GreetWithDeadLineResponse, error) {
	time.Sleep(2 * time.Second)
	return &greet_s.GreetWithDeadLineResponse{
		Output: "output from server",
	}, nil
}

func (s server) SquareRoot(ctx context.Context, request *greet_s.SquareRootRequest) (*greet_s.SquareRootResponse, error) {
	number := request.GetNumber()
	if number < 0 {
		return nil, status.Errorf(codes.InvalidArgument,
			"received a negative number")
	}
	return &greet_s.SquareRootResponse{
		NumberRoot: math.Sqrt(number),
	}, nil
}

func (s server) EveryOneGreet(stream greet_s.GreetService_EveryOneGreetServer) error {

	count := 0
	for {
		count++
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		input := req.GetInput()
		fmt.Println("received request from client" + strconv.Itoa(count) + ": " + input)
		err = stream.Send(&greet_s.EveryOneGreetResponse{
			Output: "message from server to client" + strconv.Itoa(count),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

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
