package client_authSvc

import (
	"fmt"

	config_authSvc "github.com/ashkarax/ciao_socialMedia_authService/pkg/infrastructure/config"
	"github.com/ashkarax/ciao_socialMedia_authService/pkg/infrastructure/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitPostnrelServiceClient(config *config_authSvc.Config) (*pb.PostNrelServiceClient, error) {
	cc, err := grpc.Dial(config.PortMngr.PostNrelSvcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("-------", err)
		return nil, err
	}

	Client := pb.NewPostNrelServiceClient(cc)

	return &Client, nil
}
