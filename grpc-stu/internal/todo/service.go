package todo

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc-stu/proto"
)

type Service struct {
	proto.UnimplementedTodoServiceServer
	tasks map[string]string
	lock  sync.Mutex
}

func NewService() *Service {
	return &Service{
		tasks: make(map[string]string),
	}
}

func (s *Service) AddTask(ctx context.Context, request *proto.AddTaskRequest) (*proto.AddTaskResponse, error) {
	if request.GetTask() == "" {
		return nil, status.Error(codes.InvalidArgument, "task cannot be empty")
	}

	id := uuid.New().String()
	s.lock.Lock()
	s.tasks[id] = request.GetTask()
	s.lock.Unlock()

	return &proto.AddTaskResponse{
		Id: id,
	}, nil
}

func (s *Service) CompleteTask(ctx context.Context, request *proto.CompleteTaskRequest) (*proto.CompleteTaskResponse, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.tasks[request.GetId()]; !ok {
		return nil, status.Error(codes.NotFound, "task not found")
	}

	delete(s.tasks, request.GetId())

	return &proto.CompleteTaskResponse{Success: true}, nil
}

func (s *Service) ListTasks(ctx context.Context, request *proto.ListTasksRequest) (*proto.ListTasksResponse, error) {
	tasks := make([]*proto.Task, 0, len(s.tasks))

	for id, task := range s.tasks {
		tasks = append(tasks, &proto.Task{
			Id:   id,
			Task: task,
		})
	}

	return &proto.ListTasksResponse{
		Tasks: tasks,
	}, nil
}
