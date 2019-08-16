package storage

import (
	"bytes"
	"encoding/json"
	pb "github.com/GoogleCloudPlatform/microservices-demo/src/storageservice/genproto"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"net/http"
)

type S3 struct {
	session    *s3.S3
	bucketName string
}

var s3client *S3

func NewS3Connection(accessKey, secretKey, region, bucketName string) (*S3, error) {
	if s3client != nil {
		return s3client, nil
	}
	s, err := createAWSSession(accessKey, secretKey, region)
	if err != nil {
		return nil, err

	}
	s3client = new(S3)
	s3client.session = s3.New(s)
	s3client.bucketName = bucketName
	return s3client, nil
}
func createAWSSession(accessKey, secretKey, region string) (*session.Session, error) {
	creds := credentials.NewStaticCredentials(accessKey, secretKey, "")
	return session.NewSession(&aws.Config{Region: aws.String(region), Credentials: creds})
}

func (s *S3) Store(in *pb.StorageRequest) error {
	raw, err := json.Marshal(in)
	if err != nil {

		return err
	}
	_, err = s.session.PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(s.bucketName),
		Key:                  aws.String(in.TrackingId),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(raw),
		ContentLength:        aws.Int64(int64(len(raw))),
		ContentType:          aws.String(http.DetectContentType(raw)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {

		return err
	}
	return nil
}
