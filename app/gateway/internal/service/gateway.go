package service

import (
	"context"

	commonv1 "rag/api/common/v1"
	pb "rag/api/gateway/v1"
	"rag/app/gateway/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GatewayService struct {
	pb.UnimplementedGatewayServer

	authUc  *biz.AuthUsecase
	queryUc *biz.QueryUsecase
	docUc   *biz.DocumentUsecase
	log     *log.Helper
}

func NewGatewayService(authUc *biz.AuthUsecase, queryUc *biz.QueryUsecase, docUc *biz.DocumentUsecase, logger log.Logger) *GatewayService {
	return &GatewayService{
		authUc:  authUc,
		queryUc: queryUc,
		docUc:   docUc,
		log:     log.NewHelper(logger),
	}
}

// Login handles user authentication
func (s *GatewayService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	s.log.WithContext(ctx).Info("Login request received")
	return s.authUc.Login(ctx, req)
}

// RefreshToken handles token refresh
func (s *GatewayService) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	s.log.WithContext(ctx).Info("RefreshToken request received")
	return s.authUc.RefreshToken(ctx, req)
}

// Query handles intelligent query processing
func (s *GatewayService) Query(ctx context.Context, req *pb.QueryRequest) (*pb.QueryResponse, error) {
	s.log.WithContext(ctx).Info("Query request received")
	return s.queryUc.ProcessQuery(ctx, req)
}

// UploadDocument handles document upload
func (s *GatewayService) UploadDocument(ctx context.Context, req *pb.UploadDocumentRequest) (*pb.UploadDocumentResponse, error) {
	s.log.WithContext(ctx).Info("UploadDocument request received")
	return s.docUc.UploadDocument(ctx, req)
}

// GetDocument retrieves document information
func (s *GatewayService) GetDocument(ctx context.Context, req *pb.GetDocumentRequest) (*pb.GetDocumentResponse, error) {
	s.log.WithContext(ctx).Info("GetDocument request received")
	return s.docUc.GetDocument(ctx, req)
}

// DeleteDocument handles document deletion
func (s *GatewayService) DeleteDocument(ctx context.Context, req *pb.DeleteDocumentRequest) (*pb.DeleteDocumentResponse, error) {
	s.log.WithContext(ctx).Info("DeleteDocument request received")
	return s.docUc.DeleteDocument(ctx, req)
}

// ListDocuments retrieves a list of documents
func (s *GatewayService) ListDocuments(ctx context.Context, req *pb.ListDocumentsRequest) (*pb.ListDocumentsResponse, error) {
	s.log.WithContext(ctx).Info("ListDocuments request received")
	return s.docUc.ListDocuments(ctx, req)
}

// UpdateDocumentMetadata updates document metadata
func (s *GatewayService) UpdateDocumentMetadata(ctx context.Context, req *pb.UpdateDocumentMetadataRequest) (*pb.UpdateDocumentMetadataResponse, error) {
	s.log.WithContext(ctx).Info("UpdateDocumentMetadata request received")
	return s.docUc.UpdateDocumentMetadata(ctx, req)
}

// HealthCheck performs health check
func (s *GatewayService) HealthCheck(ctx context.Context, req *emptypb.Empty) (*commonv1.HealthCheckResponse, error) {
	s.log.WithContext(ctx).Info("HealthCheck request received")

	return &commonv1.HealthCheckResponse{
		Status:    "SERVING",
		Service:   "gateway",
		Version:   "v1.0.0",
		Timestamp: timestamppb.Now(),
		Details: map[string]string{
			"uptime": "running",
			"status": "healthy",
		},
	}, nil
}
