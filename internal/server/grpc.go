package server

import (
	"context"
	"fmt"
	"net"
	"strings"

	"google.golang.org/grpc"

	"config-analyzer/internal/analyzer"
	"config-analyzer/internal/models"

	pb "config-analyzer/api/proto/analyzerpb"
)

// GRPCServer предоставляет gRPC API для анализа конфигураций.
type GRPCServer struct {
	pb.UnimplementedAnalyzerServiceServer
	analyzer *analyzer.Analyzer
	addr     string
}

// NewGRPCServer создаёт новый gRPC-сервер.
func NewGRPCServer(addr string, a *analyzer.Analyzer) *GRPCServer {
	return &GRPCServer{
		analyzer: a,
		addr:     addr,
	}
}

// Start запускает gRPC-сервер.
func (s *GRPCServer) Start() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("не удалось начать слушать: %w", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAnalyzerServiceServer(grpcServer, s)

	fmt.Printf("gRPC-сервер запущен на %s\n", s.addr)
	return grpcServer.Serve(lis)
}

// Analyze реализует gRPC-метод анализа конфигурации.
func (s *GRPCServer) Analyze(ctx context.Context, req *pb.AnalyzeRequest) (*pb.AnalyzeResponse, error) {
	if req.Config == "" {
		return nil, fmt.Errorf("поле config обязательно")
	}

	result, err := s.analyzer.Analyze(ctx, models.AnalysisRequest{
		Reader: strings.NewReader(req.Config),
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка анализа: %w", err)
	}

	resp := &pb.AnalyzeResponse{
		Count: int32(len(result.Issues)),
	}

	for _, issue := range result.Issues {
		resp.Issues = append(resp.Issues, &pb.Issue{
			Severity:       issue.Severity.String(),
			Description:    issue.Description,
			Recommendation: issue.Recommendation,
			Path:           issue.Path,
		})
	}

	return resp, nil
}
