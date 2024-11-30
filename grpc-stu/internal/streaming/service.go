package streaming

import (
	"errors"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"grpc-stu/proto"
)

type Service struct {
	proto.UnimplementedStreamingServiceServer
}

func (s *Service) StreamServerTime(request *proto.StreamServerTimeRequest, server proto.StreamingService_StreamServerTimeServer) error {
	if request.GetIntervalSeconds() == 0 {
		return status.Error(codes.InvalidArgument, "IntervalSeconds must be set")
	}

	interval := time.Duration(request.GetIntervalSeconds()) * time.Second
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-server.Context().Done():
			return nil
		case <-ticker.C:
			currentTime := time.Now()
			resp := proto.StreamServerTimeResponse{CurrentTime: timestamppb.New(currentTime)}

			if err := server.Send(&resp); err != nil {
				return err
			}
		}
	}
}

func (s *Service) LogStream(stream proto.StreamingService_LogStreamServer) error {
	count := 0

	for {
		logEntry, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return stream.SendAndClose(&proto.LogStreamResponse{EntriesLogged: int32(count)})
			}
			return err
		}

		log.Printf("Received log [%s]: %s - %s", logEntry.GetTimestamp().AsTime(), logEntry.GetLevel().String(), logEntry.GetMsg())
		count++
	}
}

func (s *Service) Echo(stream proto.StreamingService_EchoServer) error {
	var (
		req  *proto.EchoRequest
		resp *proto.EchoResponse
		err  error
	)
	for {
		req, err = stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}

		log.Printf("message recevied: %s", req.GetMessage())

		resp = &proto.EchoResponse{Message: req.GetMessage()}
		if err := stream.Send(resp); err != nil {

		}
	}
}
