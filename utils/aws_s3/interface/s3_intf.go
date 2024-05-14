package interface_awss3_authSvc

import "github.com/aws/aws-sdk-go/aws/session"

type IAwsS3 interface {
	AWSSessionInitializer() (*session.Session, error)
	AWSS3MediaUploader(data *[]byte, contentType *string, sess *session.Session, bucketFolder *string) (*string, error)
}
