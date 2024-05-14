package aws_authSvc

import (
	"bytes"
	"fmt"

	config_authSvc "github.com/ashkarax/ciao_socialMedia_authService/pkg/infrastructure/config"
	interface_awss3_authSvc "github.com/ashkarax/ciao_socialMedia_authService/utils/aws_s3/interface"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
)

type awsS3Service struct {
	s3Credentials config_authSvc.AWS
}

func AWSS3FileUploaderSetup(s3cred config_authSvc.AWS) interface_awss3_authSvc.IAwsS3 {
	return &awsS3Service{s3Credentials: s3cred}
}

func (s *awsS3Service) AWSSessionInitializer() (*session.Session, error) {
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(s.s3Credentials.Region),
			Credentials: credentials.NewStaticCredentials(
				s.s3Credentials.AccessKey,
				s.s3Credentials.SecrectKey,
				"",
			),
			Endpoint: aws.String(s.s3Credentials.Endpoint),
		},
	)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func (s *awsS3Service) AWSS3MediaUploader(data *[]byte, contentType *string, sess *session.Session, bucketFolder *string) (*string, error) {

	reader := bytes.NewReader(*data)

	randomName := uuid.New().String()

	uploader := s3manager.NewUploader(sess)

	upload, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(*bucketFolder),
		Key:         aws.String(randomName),
		Body:        reader,
		ContentType: aws.String(*contentType),
		ACL:         aws.String("public-read"),
	})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &upload.Location, nil
}
