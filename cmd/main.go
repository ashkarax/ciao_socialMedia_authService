package main

import (
	"fmt"
	"log"
	"net"

	config_authSvc "github.com/ashkarax/ciao_socialMedia_authService/pkg/infrastructure/config"
	di_authSvc "github.com/ashkarax/ciao_socialMedia_authService/pkg/infrastructure/di"
	"github.com/ashkarax/ciao_socialMedia_authService/pkg/infrastructure/pb"
	"google.golang.org/grpc"
)

func main() {

	config, err := config_authSvc.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	server, serErr := di_authSvc.InitializeAuthServer(config)
	if serErr != nil {
		log.Fatalf("error with initializing di:%v", serErr)
	}

	lis, err := net.Listen("tcp", config.PortMngr.RunnerPort)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Auth Service started on:", config.PortMngr.RunnerPort)

	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, server)

	// Log every connection attempt to the server
	go func() {
		for {
			conn, err := lis.Accept()
			if err != nil {
				log.Println("Error accepting connection:", err)
				continue
			}
			log.Println("New connection from:", conn.RemoteAddr())

			// Optionally read from the connection and log data (for demonstration purposes)
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				log.Println("Error reading from connection:", err)
				return
			}
			log.Printf("Received data: %s", string(buf[:n]))
		}
	}()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to start Auth_service server:%v", err)
	}
}
