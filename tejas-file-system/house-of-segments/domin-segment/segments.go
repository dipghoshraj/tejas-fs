package dominsegment

import (
	"fmt"
	"hosue-of-segments/domin-segment/proto"
	"io"
	"os"

	"google.golang.org/grpc"
)

func (s *StorageServer) StoreSegment(stream grpc.ClientStreamingServer[proto.SegmentMessage, proto.SegmentResponse]) error {

	var file *os.File
	segments, err := stream.Recv()
	if err != nil {
		return stream.SendAndClose(&proto.SegmentResponse{Success: false})
	}

	segName := segments.GetSegmentName()
	size := int64(0)

	for {
		segment, err := stream.Recv()

		if err != nil {
			if file != nil {
				file.Close()
			}
			if err == io.EOF {
				fmt.Printf("error for %v", err)
				return stream.SendAndClose(&proto.SegmentResponse{ID: segName, Size: size, Name: segName, Success: true})
			}
			return stream.SendAndClose(&proto.SegmentResponse{Success: false})
		}

		if file == nil {
			// TODO : Intialize Orbs here
			file, err = os.Create(segName)
			if err != nil {
				file.Close()
				return fmt.Errorf("%v", err)
			}
		}

		_, err = file.Write(segment.GetOrb())
		fmt.Printf("file chunk size %d \n", len(segment.GetOrb()))
		size += int64(len(segment.GetOrb()))
		if err != nil {
			file.Close()
			return fmt.Errorf("%v", err)
		}
	}

}
