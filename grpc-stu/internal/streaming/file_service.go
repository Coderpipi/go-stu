package streaming

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc-stu/proto"
)

type FileService struct {
	proto.UnimplementedFileUploadServiceServer
}

func (f *FileService) DownloadFile(request *proto.DownloadFileRequest, server proto.FileUploadService_DownloadFileServer) error {
	if request.GetName() == "" {
		return status.Error(codes.InvalidArgument, "file name cannot be empty")
	}

	// open file
	file, err := os.Open(request.GetName())
	if err != nil {
		if os.IsNotExist(err) {
			return status.Error(codes.NotFound, "file not found")
		}
		return err
	}

	const bufferSize = 5 * 1024

	buf := make([]byte, bufferSize)

	for {
		bytes, err := file.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return status.Error(codes.Internal, "error reading file")
		}

		err = server.Send(&proto.DownloadFileResponse{Content: buf[:bytes]})
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *FileService) UploadFile(server proto.FileUploadService_UploadFileServer) error {
	fileName := fmt.Sprintf("%s.png", uuid.New().String())

	file, err := os.Create(fileName)
	if err != nil {
		return status.Error(codes.Internal, "error create file")
	}
	defer file.Close()

	for {
		res, err := server.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return server.SendAndClose(&proto.UploadFileResponse{Name: fileName})
			}
			return err
		}

		if _, err := file.Write(res.Content); err != nil {
			return status.Error(codes.Internal, "error writing to file")
		}
	}
}
