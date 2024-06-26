package di_authSvc

import (
	"fmt"

	client_authSvc "github.com/ashkarax/ciao_socialMedia_authService/pkg/infrastructure/client"
	config_authSvc "github.com/ashkarax/ciao_socialMedia_authService/pkg/infrastructure/config"
	db_authSvc "github.com/ashkarax/ciao_socialMedia_authService/pkg/infrastructure/db"
	server_authSvc "github.com/ashkarax/ciao_socialMedia_authService/pkg/infrastructure/server"
	repository_authSvc "github.com/ashkarax/ciao_socialMedia_authService/pkg/repository"
	usecase_authSvc "github.com/ashkarax/ciao_socialMedia_authService/pkg/usecase"
	aws_authSvc "github.com/ashkarax/ciao_socialMedia_authService/utils/aws_s3"
	gosmtp_authSvc "github.com/ashkarax/ciao_socialMedia_authService/utils/go_smtp"
	hashpassword_authSvc "github.com/ashkarax/ciao_socialMedia_authService/utils/hash_password"
	jwttoken_authSvc "github.com/ashkarax/ciao_socialMedia_authService/utils/jwt.go"
	randnumgene_authSvc "github.com/ashkarax/ciao_socialMedia_authService/utils/random_number"
	regex_authSvc "github.com/ashkarax/ciao_socialMedia_authService/utils/regex"
)

func InitializeAuthServer(config *config_authSvc.Config) (*server_authSvc.AuthService, error) {

	hashUtil := hashpassword_authSvc.NewHashUtil()

	DB, err := db_authSvc.ConnectDatabase(&config.DB, hashUtil)
		if err != nil {
			fmt.Println("ERROR CONNECTING DB FROM DI.GO")
			return nil, err
		}

	smtpUtil := gosmtp_authSvc.NewSmtpUtils(&config.Smtp)
	jwtUtil := jwttoken_authSvc.NewJwtUtil()
	randNumUtil := randnumgene_authSvc.NewRandomNumUtil()
	regexUtli := regex_authSvc.NewRegexUtil()
	awsS3util := aws_authSvc.AWSS3FileUploaderSetup(config.AwsS3)

	postNrelClient, err := client_authSvc.InitPostnrelServiceClient(config)
	if err != nil {
		return nil, err
	}

	userRepo := repository_authSvc.NewUserRepo(DB)
	userUseCase := usecase_authSvc.NewUserUseCase(userRepo, smtpUtil, jwtUtil, randNumUtil, regexUtli, &config.Token, hashUtil, postNrelClient, awsS3util)

	jwtUseCase := usecase_authSvc.NewJwtUseCase(&config.Token, jwtUtil, userRepo)

	embeddingStruct := server_authSvc.NewAuthServiceServer(userUseCase, jwtUseCase)

	return embeddingStruct, nil
}
