package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"playground/greet_s/greet_s"
	"strconv"
	"time"
)

func main() {
	fmt.Println("hello, greet gPRC client")
	cnn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cannot connect %v", err)
	}

	defer cnn.Close()

	c := greet_s.NewGreetServiceClient(cnn)
	//doServiceStreaming(c)
	//doClientStreaming(c)
	//doBiDirectionStreaming(c)
	doError(c)
}

func doError(c greet_s.GreetServiceClient) {
	res, err := c.SquareRoot(context.Background(), &greet_s.SquareRootRequest{Number: -80})
	if err != nil {
		resErr, ok := status.FromError(err)
		if ok {
			// actual error from gRPC, a purposed error
			if resErr.Code() == codes.InvalidArgument {
				fmt.Println("we probably sent a negative number")
			}
		} else {
			log.Fatalf("error when calling SquareRoot: %v", err)
		}
		fmt.Println("errors when calling SquareRoot: " + err.Error())
	}
	if res != nil {
		fmt.Printf("result of square root -80 is %v", res.GetNumberRoot())
	}

	res, err = c.SquareRoot(context.Background(), &greet_s.SquareRootRequest{Number: 80})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("result of square root 80 is %v", res.GetNumberRoot())

}

func doBiDirectionStreaming(c greet_s.GreetServiceClient) {
	stream, err := c.EveryOneGreet(context.Background())
	if err != nil {
		log.Fatalf("doBiDirectionStreaming cannot connect %v", err)
	}

	waitc := make(chan struct{})

	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println("message from c to s, count " + strconv.Itoa(i))
			err := stream.Send(&greet_s.EveryOneGreetRequest{
				Input: "message from c to s, count " + strconv.Itoa(i),
			})
			if err != nil {
				log.Fatalf("err at line 40 %v", err)
			}
			time.Sleep(1 * time.Second)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("cient side get EOF error%v", err)
				close(waitc)
				break
			}
			if err != nil {
				log.Fatalf("err at line 54 %v", err)
			}
			fmt.Println("received from server: " + res.GetOutput())
		}
	}()

	<-waitc
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
