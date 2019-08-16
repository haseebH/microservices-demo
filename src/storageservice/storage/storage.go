package storage

import pb "github.com/GoogleCloudPlatform/microservices-demo/src/storageservice/genproto"

type Storage interface {
	Store(in *pb.StorageRequest) error
}
