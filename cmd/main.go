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

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to start Auth_service server:%v", err)
	}
}
