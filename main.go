package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"os"

	grpcpeer "github.com/dioxine/gogrpcpeer/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type server struct {
	grpcpeer.UnimplementedUserServiceServer
}

func (s *server) User(ctx context.Context, req *grpcpeer.UserRequest) (*grpcpeer.UserResponse, error) {
	fmt.Printf("User function was called with %v\n", req)
	user := req.GetUser()
	res := &grpcpeer.UserResponse{User: user, MessageResponse: fmt.Sprintf("Hello, %v", user.FirstName)}
	return res, nil
}

func main() {

	cert, err := tls.LoadX509KeyPair("cert/server/public/server.crt", "cert/server/private/server.pem")
	if err != nil {
		log.Fatalf("failed to load key pair: %s", err)
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
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    ca,
	}

	s := grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig)))
	grpcpeer.RegisterUserServiceServer(s, &server{})
	fmt.Println("Server is starting...")
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if lis != nil {
		fmt.Println("Server is started")
	}
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
