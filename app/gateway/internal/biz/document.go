package biz

import (
	"context"
	"time"

	commonv1 "rag/api/common/v1"
	v1 "rag/api/gateway/v1"

	"github.com/go-kratos/kratos/v2/log"
)

// Document represents document information
type Document struct {
	DocumentID  string
	Title       string
	FileType    string
	FileSize    int64
	TotalChunks int32
	Content     string
	Metadata    *commonv1.Metadata
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// DocumentProcessingProgress represents document processing progress
type DocumentProcessingProgress struct {
	ProgressPercentage  float32
	CurrentStage        string
	StatusMessage       string
	StartedAt           time.Time
	EstimatedCompletion time.Time
}

// DocumentStats represents document statistics
type DocumentStats struct {
	TotalChunks      int32
	TotalTokens      int32
	AverageChunkSize int32
	QueryCount       int32
	LastAccessed     time.Time
}

// DocumentRepo defines the data access interface for document management
type DocumentRepo interface {
	// 存储文档
	SaveDocument(ctx context.Context, doc *Document) error
	// 获取文档信息
	GetDocument(ctx context.Context, documentID string, includeFields []string) (*Document, error)
	// 删除文档
	DeleteDocument(ctx context.Context, documentID string, options *v1.DeleteOptions) error
	// 列出文档
	ListDocuments(ctx context.Context, pagination *commonv1.PaginationRequest, filters []*commonv1.Filter) ([]*Document, *commonv1.PaginationResponse, error)
	// 更新文档元数据
	UpdateDocumentMetadata(ctx context.Context, documentID string, metadata *commonv1.Metadata, updateFields []string) error
	// 获取文档统计信息
	GetDocumentStats(ctx context.Context, documentID string) (*DocumentStats, error)
	// 调用文档存储服务
	CallDocstoreService(ctx context.Context, operation string, data interface{}) (interface{}, error)
	// 调用预处理服务
	CallPreprocessorService(ctx context.Context, content []byte, options *v1.DocumentProcessingOptions) (*DocumentProcessingProgress, error)
	// 调用嵌入服务
	CallEmbeddingService(ctx context.Context, documentID string, chunks []string) error
}

// DocumentUsecase handles document management business logic
type DocumentUsecase struct {
	repo DocumentRepo
	log  *log.Helper
}

// NewDocumentUsecase creates a new document usecase
func NewDocumentUsecase(repo DocumentRepo, logger log.Logger) *DocumentUsecase {
	return &DocumentUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// UploadDocument handles document upload and processing
func (uc *DocumentUsecase) UploadDocument(ctx context.Context, req *v1.UploadDocumentRequest) (*v1.UploadDocumentResponse, error) {
	uc.log.WithContext(ctx).Infof("Uploading document: %s", req.Title)

	// 生成文档 ID
	documentID := generateDocumentID()

	// 创建文档对象
	doc := &Document{
		DocumentID: documentID,
		Title:      req.Title,
		FileType:   req.FileType,
		FileSize:   int64(len(req.FileContent)),
		Metadata:   req.Metadata,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// 保存文档到存储
	if err := uc.repo.SaveDocument(ctx, doc); err != nil {
		uc.log.WithContext(ctx).Errorf("Failed to save document: %v", err)
		return nil, err
	}

	// 调用预处理服务处理文档
	progress, err := uc.repo.CallPreprocessorService(ctx, req.FileContent, req.ProcessingOptions)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("Failed to process document: %v", err)
		return &v1.UploadDocumentResponse{
			DocumentId:   documentID,
			UploadStatus: "failed",
			DocumentInfo: &commonv1.DocumentInfo{
				DocumentId: documentID,
				Title:      req.Title,
				FileType:   req.FileType,
				FileSize:   int64(len(req.FileContent)),
				Metadata:   req.Metadata,
			},
		}, err
	}

	response := &v1.UploadDocumentResponse{
		DocumentId:   documentID,
		UploadStatus: "processing",
		Progress: &v1.DocumentProcessingProgress{
			ProgressPercentage: progress.ProgressPercentage,
			CurrentStage:       progress.CurrentStage,
			StatusMessage:      progress.StatusMessage,
		},
		DocumentInfo: &commonv1.DocumentInfo{
			DocumentId: documentID,
			Title:      req.Title,
			FileType:   req.FileType,
			FileSize:   int64(len(req.FileContent)),
			Metadata:   req.Metadata,
		},
	}

	uc.log.WithContext(ctx).Infof("Document upload initiated: %s", documentID)
	return response, nil
}

// GetDocument retrieves document information
func (uc *DocumentUsecase) GetDocument(ctx context.Context, req *v1.GetDocumentRequest) (*v1.GetDocumentResponse, error) {
	uc.log.WithContext(ctx).Infof("Getting document: %s", req.DocumentId)

	doc, err := uc.repo.GetDocument(ctx, req.DocumentId, req.IncludeFields)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("Failed to get document: %v", err)
		return nil, err
	}

	stats, err := uc.repo.GetDocumentStats(ctx, req.DocumentId)
	if err != nil {
		uc.log.WithContext(ctx).Warnf("Failed to get document stats: %v", err)
		stats = &DocumentStats{}
	}

	response := &v1.GetDocumentResponse{
		DocumentInfo: &commonv1.DocumentInfo{
			DocumentId:  doc.DocumentID,
			Title:       doc.Title,
			FileType:    doc.FileType,
			FileSize:    doc.FileSize,
			TotalChunks: doc.TotalChunks,
			Metadata:    doc.Metadata,
		},
		Content: doc.Content,
		Stats: &v1.DocumentStats{
			TotalChunks:      stats.TotalChunks,
			TotalTokens:      stats.TotalTokens,
			AverageChunkSize: stats.AverageChunkSize,
			QueryCount:       stats.QueryCount,
		},
	}

	uc.log.WithContext(ctx).Infof("Document retrieved: %s", req.DocumentId)
	return response, nil
}

// DeleteDocument handles document deletion
func (uc *DocumentUsecase) DeleteDocument(ctx context.Context, req *v1.DeleteDocumentRequest) (*v1.DeleteDocumentResponse, error) {
	uc.log.WithContext(ctx).Infof("Deleting document: %s", req.DocumentId)

	err := uc.repo.DeleteDocument(ctx, req.DocumentId, req.Options)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("Failed to delete document: %v", err)
		return &v1.DeleteDocumentResponse{
			Status: "failed",
		}, err
	}

	response := &v1.DeleteDocumentResponse{
		Status: "deleted",
		CleanupInfo: &v1.CleanupInfo{
			DeletionCompletedAt: nil, // 设置为当前时间的 protobuf timestamp
		},
	}

	uc.log.WithContext(ctx).Infof("Document deleted: %s", req.DocumentId)
	return response, nil
}

// ListDocuments retrieves a list of documents
func (uc *DocumentUsecase) ListDocuments(ctx context.Context, req *v1.ListDocumentsRequest) (*v1.ListDocumentsResponse, error) {
	uc.log.WithContext(ctx).Info("Listing documents")

	docs, pagination, err := uc.repo.ListDocuments(ctx, req.Pagination, req.Filters)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("Failed to list documents: %v", err)
		return nil, err
	}

	// 转换为 protobuf 格式
	docInfos := make([]*commonv1.DocumentInfo, len(docs))
	for i, doc := range docs {
		docInfos[i] = &commonv1.DocumentInfo{
			DocumentId:  doc.DocumentID,
			Title:       doc.Title,
			FileType:    doc.FileType,
			FileSize:    doc.FileSize,
			TotalChunks: doc.TotalChunks,
			Metadata:    doc.Metadata,
		}
	}

	response := &v1.ListDocumentsResponse{
		Documents:  docInfos,
		Pagination: pagination,
		Stats: &v1.DocumentListStats{
			TotalDocuments: pagination.Total,
		},
	}

	uc.log.WithContext(ctx).Infof("Listed %d documents", len(docs))
	return response, nil
}

// UpdateDocumentMetadata updates document metadata
func (uc *DocumentUsecase) UpdateDocumentMetadata(ctx context.Context, req *v1.UpdateDocumentMetadataRequest) (*v1.UpdateDocumentMetadataResponse, error) {
	uc.log.WithContext(ctx).Infof("Updating document metadata: %s", req.DocumentId)

	err := uc.repo.UpdateDocumentMetadata(ctx, req.DocumentId, req.Metadata, req.UpdateFields)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("Failed to update document metadata: %v", err)
		return &v1.UpdateDocumentMetadataResponse{
			Status: "failed",
		}, err
	}

	response := &v1.UpdateDocumentMetadataResponse{
		Status:          "updated",
		UpdatedMetadata: req.Metadata,
		MetadataVersion: 1,
	}

	uc.log.WithContext(ctx).Infof("Document metadata updated: %s", req.DocumentId)
	return response, nil
}

// generateDocumentID generates a unique document ID
func generateDocumentID() string {
	return "doc-" + time.Now().Format("20060102150405") + "-" + randomString(8)
}
