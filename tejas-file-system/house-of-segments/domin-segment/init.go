package dominsegment

// protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative domin-segment\proto\storage.proto

import (
	"hosue-of-segments/domin-segment/proto"
)

type StorageServer struct {
	proto.UnimplementedStorageBoxServiceServer
}
