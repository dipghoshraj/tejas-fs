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

	segName := segments.GetFileid()
	size := int64(0)

	for {
		segment, err := stream.Recv()

		if err != nil {
			return fmt.Errorf("%v", err)
		}

		if err == io.EOF {
			if file != nil {
				file.Close()
			}
			return stream.SendAndClose(&proto.SegmentResponse{ID: segName, Size: size, Name: segName, Success: true})
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
		fmt.Printf("file chunk size %d", len(segment.GetOrb()))

		if err != nil {
			file.Close()
			return fmt.Errorf("%v", err)
		}

	}

}
