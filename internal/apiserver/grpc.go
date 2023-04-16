package apiserver

import (
	"net"

	"github.com/cuizhaoyue/iams/pkg/log"
	"google.golang.org/grpc"
)

// 提供apiserver的grpc服务.
type grpcAPIServer struct {
	*grpc.Server
	address string
}

// Run 启动grpc服务
func (s *grpcAPIServer) Run() {
	// 创建监听器
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatalf("failed to listen: %s", err.Error())
	}

	// 启动grpc服务
	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to start grpc server: %s", err.Error())
		}
	}()

	log.Infof("start grpc server at %s", s.address)
}

// Close 优雅关闭grpc服务
func (s *grpcAPIServer) Close() {
	s.GracefulStop()
	log.Infof("GRPC server on %s stopped", s.address)
}
