package dominsegment

import (
	"context"
	"hosue-of-segments/domin-segment/proto"
)

func (s *StorageServer) HealthCheck(ctx context.Context, req *proto.HelloWroldMessage) (*proto.HelloWroldResponse, error) {
	return &proto.HelloWroldResponse{
		Status:    "inactive",
		UsedSpace: 0,
	}, nil
}
