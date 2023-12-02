package grpc

import (
	"context"
	"net"
	"time"

	"github.com/74th/vscode-book-r2-golang/domain/entity"
	"github.com/74th/vscode-book-r2-golang/domain/usecase"
	"github.com/74th/vscode-book-r2-golang/gateway/server/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// サーバAPI
type Server struct {
	pb.UnimplementedTodoListServiceServer
	addr       string
	listener   net.Listener
	server     *grpc.Server
	interactor *usecase.Interactor
}

func taskToPB(task *entity.Task) *pb.Task {
	return &pb.Task{
		Id:   int64(task.ID),
		Text: task.Text,
		Done: task.Done,
	}
}

// AddTask implements pb.TodoListServiceServer.
func (s *Server) AddTask(ctx context.Context, req *pb.AddTaskRequest) (*pb.AddTaskResponse, error) {
	task, err := s.interactor.CreateTask(&entity.Task{
		Text: req.Text,
	})

	if err != nil {
		return nil, status.Errorf(500, err.Error())
	}

	return &pb.AddTaskResponse{
		Task: taskToPB(task),
	}, nil
}

// DoneTask implements pb.TodoListServiceServer.
func (s *Server) DoneTask(ctx context.Context, req *pb.DoneTaskRequest) (*pb.DoneTaskResponse, error) {
	task, err := s.interactor.DoneTask(int(req.Id))

	if err != nil {
		return nil, status.Errorf(500, err.Error())
	}

	return &pb.DoneTaskResponse{
		Task: taskToPB(task),
	}, nil
}

// GetTasks implements pb.TodoListServiceServer.
func (s *Server) GetTasks(ctx context.Context, req *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {
	tasks, err := s.interactor.ShowTasks()

	if err != nil {
		return nil, status.Errorf(500, err.Error())
	}

	pbTasks := make([]*pb.Task, len(tasks))

	for i := range tasks {
		pbTasks[i] = taskToPB(tasks[i])
	}

	return &pb.GetTaskResponse{
		Tasks: pbTasks,
	}, nil
}

func New(addr string, interactor *usecase.Interactor) *Server {
	s := &Server{
		server:     grpc.NewServer(),
		addr:       addr,
		interactor: interactor,
	}
	return s
}

func (s *Server) Serve() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.listener = listener

	grpcServer := grpc.NewServer()
	pb.RegisterTodoListServiceServer(grpcServer, s)

	errChan := make(chan error, 1)
	timer := time.NewTimer(100 * time.Millisecond)
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return err
	case <-timer.C:
	}

	return nil
}

func (s *Server) Close() error {
	s.server.Stop()
	return s.listener.Close()
}
