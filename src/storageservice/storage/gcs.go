package storage

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"fmt"
	pb "github.com/GoogleCloudPlatform/microservices-demo/src/storageservice/genproto"
)

type GCSClient struct {
	client     *storage.Client
	bucketName string
}

var gcsClient *GCSClient

func NewGcsConnection(bucketName string) (*GCSClient, error) {
	if gcsClient != nil {
		return gcsClient, nil
	}
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client: %v", err)
	}
	gcsClient = new(GCSClient)
	gcsClient.client = client
	gcsClient.bucketName = bucketName
	return gcsClient, nil
}

func (g *GCSClient) Store(in *pb.StorageRequest) error {
	bucketHandler := g.client.Bucket(g.bucketName)
	writer := bucketHandler.Object(in.TrackingId).NewWriter(context.Background())
	raw, err := json.Marshal(in)
	if err != nil {

		return err
	}
	defer writer.Close()
	if _, err := writer.Write(raw); err != nil {
		return err
	}
	return nil
}
