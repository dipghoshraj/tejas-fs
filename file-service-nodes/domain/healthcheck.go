package domain

import (
	"context"

	"github.com/dipghoshraj/media-service/file-service-nodes/diskstorage"
	"github.com/dipghoshraj/media-service/file-service-nodes/domain/proto"
)

func (s *StorageServer) HealthCheck(ctx context.Context, req *proto.HealthCheckMessage) (*proto.HelloReply, error) {
	usedSapce, err := diskstorage.GetUsedSpace()
	if err != nil {
		return &proto.HelloReply{
			Status:    "inactive",
			UsedSpace: 0,
		}, nil
	}
	return &proto.HelloReply{
		Status:    "active",
		UsedSpace: usedSapce,
	}, nil
}
