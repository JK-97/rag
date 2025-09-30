package service

import (
	"context"

	pb "rag/api/docstore/v1"
	"rag/app/docstore/internal/biz"
)

type DocstoreService struct {
	pb.UnimplementedDocStoreServer
	uc *biz.GreeterUsecase
}

func NewDocstoreService(uc *biz.GreeterUsecase) *DocstoreService {
	return &DocstoreService{}
}

func (s *DocstoreService) CreateDocstore(ctx context.Context, req *pb.UploadDocumentRequest) (*pb.UploadDocumentResponse, error) {
	return &pb.UploadDocumentResponse{}, nil
}
