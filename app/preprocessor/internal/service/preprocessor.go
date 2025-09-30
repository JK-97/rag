package service

import (
	"context"

	pb "rag/api/preprocessor/v1"
)

type PreprocessorService struct {
	pb.UnimplementedPreprocessorServer
}

func NewPreprocessorService() *PreprocessorService {
	return &PreprocessorService{}
}

func (s *PreprocessorService) CreatePreprocessor(ctx context.Context, req *pb.CreatePreprocessorRequest) (*pb.CreatePreprocessorReply, error) {
	return &pb.CreatePreprocessorReply{}, nil
}
func (s *PreprocessorService) UpdatePreprocessor(ctx context.Context, req *pb.UpdatePreprocessorRequest) (*pb.UpdatePreprocessorReply, error) {
	return &pb.UpdatePreprocessorReply{}, nil
}
func (s *PreprocessorService) DeletePreprocessor(ctx context.Context, req *pb.DeletePreprocessorRequest) (*pb.DeletePreprocessorReply, error) {
	return &pb.DeletePreprocessorReply{}, nil
}
func (s *PreprocessorService) GetPreprocessor(ctx context.Context, req *pb.GetPreprocessorRequest) (*pb.GetPreprocessorReply, error) {
	return &pb.GetPreprocessorReply{}, nil
}
func (s *PreprocessorService) ListPreprocessor(ctx context.Context, req *pb.ListPreprocessorRequest) (*pb.ListPreprocessorReply, error) {
	return &pb.ListPreprocessorReply{}, nil
}
