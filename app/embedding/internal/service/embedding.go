package service

import (
	"context"

	pb "rag/api/embedding/v1"
)

type EmbeddingService struct {
	pb.UnimplementedEmbeddingServer
}

func NewEmbeddingService() *EmbeddingService {
	return &EmbeddingService{}
}

func (s *EmbeddingService) CreateEmbedding(ctx context.Context, req *pb.CreateEmbeddingRequest) (*pb.CreateEmbeddingReply, error) {
	return &pb.CreateEmbeddingReply{}, nil
}
func (s *EmbeddingService) UpdateEmbedding(ctx context.Context, req *pb.UpdateEmbeddingRequest) (*pb.UpdateEmbeddingReply, error) {
	return &pb.UpdateEmbeddingReply{}, nil
}
func (s *EmbeddingService) DeleteEmbedding(ctx context.Context, req *pb.DeleteEmbeddingRequest) (*pb.DeleteEmbeddingReply, error) {
	return &pb.DeleteEmbeddingReply{}, nil
}
func (s *EmbeddingService) GetEmbedding(ctx context.Context, req *pb.GetEmbeddingRequest) (*pb.GetEmbeddingReply, error) {
	return &pb.GetEmbeddingReply{}, nil
}
func (s *EmbeddingService) ListEmbedding(ctx context.Context, req *pb.ListEmbeddingRequest) (*pb.ListEmbeddingReply, error) {
	return &pb.ListEmbeddingReply{}, nil
}
