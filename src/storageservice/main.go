// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	pb "github.com/GoogleCloudPlatform/microservices-demo/src/storageservice/genproto"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	defaultPort = "50055"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.Level = logrus.DebugLevel
	log.Formatter = &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
		TimestampFormat: time.RFC3339Nano,
	}
	log.Out = os.Stdout
}
func mustMapEnv(target *string, envKey string) {
	v := os.Getenv(envKey)
	if v == "" {
		panic(fmt.Sprintf("environment variable %q not set", envKey))
	}
	*target = v
}

func main() {
	//go initTracing()
	//go initProfiling("storageservice", "1.0.0")

	port := defaultPort
	if value, ok := os.LookupEnv("PORT"); ok {
		port = value
	}
	port = fmt.Sprintf(":%s", port)

	accessKey := ""
	secretKey := ""
	region := ""
	svc := &server{}
	mustMapEnv(&accessKey, "AWS_ACCESS_KEY")
	mustMapEnv(&secretKey, "AWS_ACCESS_SECRET")
	mustMapEnv(&svc.bucketName, "AWS_BUCKET_NAME")
	mustMapEnv(&region, "AWS_BUCKET_REGION")
	/*e := gin.New()
	e.POST("/shipping",deployService)
	e.Run( port)*/
	s, err := createSession(accessKey, secretKey, region)
	if err != nil {
		panic(fmt.Sprint(err.Error()))

	}
	svc.session = s
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer(grpc.StatsHandler(&ocgrpc.ServerHandler{}))

	pb.RegisterStorageServiceServer(srv, svc)
	healthpb.RegisterHealthServer(srv, svc)
	log.Infof("storage Service listening on port %s", port)
	// Register reflection service on gRPC server.
	reflection.Register(srv)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

type server struct {
	session    *session.Session
	bucketName string
}

func (s *server) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}
func (s *server) Watch(req *healthpb.HealthCheckRequest, ws healthpb.Health_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "health check via Watch not implemented")
}

func (s *server) StoreOrder(ctx context.Context, in *pb.StorageRequest) (*pb.StorageResponse, error) {
	resp := new(pb.StorageResponse)
	service := s3.New(s.session)

	raw, err := json.Marshal(in)
	if err != nil {
		resp.Status = err.Error()
		return resp, err
	}
	_, err = service.PutObject(&s3.PutObjectInput{
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
		resp.Status = err.Error()
		return resp, err
	}
	resp.Status = " order data inserted successfully"
	return resp, nil
}
func createSession(accessKey, secretKey, region string) (*session.Session, error) {
	creds := credentials.NewStaticCredentials(accessKey, secretKey, "")
	return session.NewSession(&aws.Config{Region: aws.String(region), Credentials: creds})
}
