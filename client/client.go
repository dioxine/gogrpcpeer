package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"time"

	gogrpcpeer "github.com/dioxine/gogrpcpeer/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	fmt.Println("Hello, World! This is client")

	cert, err := tls.LoadX509KeyPair("cert/client/public/client.crt", "cert/client/private/client.pem")
	if err != nil {
		log.Fatalf("failed to load client cert: %v", err)
	}

	ca := x509.NewCertPool()
	caFilePath := "cert/ca/public/ca.crt"
	caBytes, err := os.ReadFile(caFilePath)
	if err != nil {
		log.Fatalf("failed to read ca cert %q: %v", caFilePath, err)
	}
	if ok := ca.AppendCertsFromPEM(caBytes); !ok {
		log.Fatalf("failed to parse %q", caFilePath)
	}

	tlsConfig := &tls.Config{
		ServerName:   "localhost",
		Certificates: []tls.Certificate{cert},
		RootCAs:      ca,
	}

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	c := gogrpcpeer.NewUserServiceClient(conn)
	doUnary(c)

	//doServerStreaming(c)

	//doClientStreaming(c)
}

func doUnary(c gogrpcpeer.UserServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	fmt.Println("Sending request to server...")

	time.Sleep(1 * time.Second)

	req := &gogrpcpeer.UserRequest{
		User: &gogrpcpeer.User{
			Id:         1,
			FirstName:  "John",
			SecondName: "Doe",
			ThirdName:  "Doe",
		},
	}

	res, err := c.User(ctx, req)

	if err != nil {
		log.Fatalf("Error when calling User: %s", err)
	}
	log.Printf("Response from server: \nID: %d, \nFirst Name: %s, \nSecond Name: %s, \nThird Name: %s, \nResulting Message: %s", res.User.Id, res.User.FirstName, res.User.SecondName, res.User.ThirdName, res.MessageResponse)
}

/*
func doServerStreaming(c gogrpcpeer.UserServiceClient) {
	fmt.Println("Server streaming")
	req := &gogrpcpeer.UserMultipleRequest{User: &gogrpcpeer.User{Id: 1, FirstName: "John", SecondName: "Doe", ThirdName: "Belmondoe"}}
	resStream, err := c.UserMultiple(context.Background(), req)
	if err != nil {
		log.Fatalf("Error when calling UserMultiple: %s", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error when reading stream: %s", err)
		}
		log.Println(msg.GetMessageResponse())
	}
}



func doClientStreaming(c gogrpcpeer.UserServiceClient) {
	stream, err := c.LongUser(context.Background())
	if err != nil {
		log.Fatalf("Error when calling LongUser: %s", err)
	}

	requests := []*gogrpcpeer.LongUserRequest{
		{User: &gogrpcpeer.User{Id: 1, FirstName: "John", SecondName: "Doe", ThirdName: "Belmondoe"}},
		{User: &gogrpcpeer.User{Id: 2, FirstName: "Jane", SecondName: "Doe", ThirdName: "Belmondoe"}},
		{User: &gogrpcpeer.User{Id: 3, FirstName: "Joe", SecondName: "Doe", ThirdName: "Belmondoe"}},
		{User: &gogrpcpeer.User{Id: 4, FirstName: "Jim", SecondName: "Doe", ThirdName: "Belmondoe"}},
		{User: &gogrpcpeer.User{Id: 5, FirstName: "Jill", SecondName: "Doe", ThirdName: "Belmondoe"}},
		{User: &gogrpcpeer.User{Id: 6, FirstName: "Jack", SecondName: "Doe", ThirdName: "Belmondoe"}},
		{User: &gogrpcpeer.User{Id: 7, FirstName: "Jenny", SecondName: "Doe", ThirdName: "Belmondoe"}},
	}

	for _, req := range requests {
		fmt.Println("Sending request: ", req.User.GetFirstName())
		stream.Send(req)
		time.Sleep(200 * time.Millisecond)
	}

	for i := 0; i < 10; i++ {

	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error when receiving response: %s", err)
	}
	fmt.Println(res.GetMessageResponse())
}

*/
