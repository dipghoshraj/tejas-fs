package domain

import (
	"fmt"
	"io"
	"os"

	"github.com/dipghoshraj/media-service/file-service-nodes/domain/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func (s *StorageServer) IngressNode(stream grpc.ClientStreamingServer[proto.IngressStorageMessage, proto.IngressStorageResponse]) error {

	// file :=
	var file *os.File
	fileID := uuid.New().String()
	fmt.Println("stat processing")

	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("%v", err)
		}
	}()

	for {
		segment, err := stream.Recv()

		// md, ok := metadata.FromIncomingContext(stream.Context())
		// if !ok {
		// 	return fmt.Errorf("%v", err)
		// }
		filename := fileID

		if err != nil {
			return fmt.Errorf("%v", err)
		}

		if err == io.EOF {
			if file != nil {
				file.Close()
			}
			return stream.SendAndClose(&proto.IngressStorageResponse{ID: filename, Size: "3", Name: "fileName", IngressNodeId: "node"})
		}

		if file == nil {
			// TODO : Intialize Orbs here
			file, err = os.Create(filename)
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
