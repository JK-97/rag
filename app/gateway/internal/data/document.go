package data

import (
	"context"
	"fmt"
	"time"

	commonv1 "rag/api/common/v1"
	v1 "rag/api/gateway/v1"
	"rag/app/gateway/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// documentRepo implements biz.DocumentRepo interface
type documentRepo struct {
	data *Data
	log  *log.Helper
}

// NewDocumentRepo creates a new document repository
func NewDocumentRepo(data *Data, logger log.Logger) biz.DocumentRepo {
	return &documentRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// SaveDocument saves document to storage
func (r *documentRepo) SaveDocument(ctx context.Context, doc *biz.Document) error {
	r.log.WithContext(ctx).Infof("Saving document: %s", doc.DocumentID)

	// TODO: 调用 docstore 服务保存文档
	// 这里是模拟实现，实际应该通过 gRPC 调用 docstore 服务

	r.log.WithContext(ctx).Infof("Document saved successfully: %s", doc.DocumentID)
	return nil
}

// GetDocument retrieves document information
func (r *documentRepo) GetDocument(ctx context.Context, documentID string, includeFields []string) (*biz.Document, error) {
	r.log.WithContext(ctx).Infof("Getting document: %s", documentID)

	// TODO: 调用 docstore 服务获取文档
	// 这里是模拟实现，实际应该通过 gRPC 调用 docstore 服务

	// 模拟文档数据
	doc := &biz.Document{
		DocumentID:  documentID,
		Title:       "示例文档",
		FileType:    "pdf",
		FileSize:    1024000,
		TotalChunks: 10,
		Content:     "这是一个示例文档的内容...",
		Metadata: &commonv1.Metadata{
			Data: map[string]string{
				"author":   "示例作者",
				"category": "技术文档",
			},
			CreatedAt: nil, // 应该设置为当前时间的 protobuf timestamp
			UpdatedAt: nil, // 应该设置为当前时间的 protobuf timestamp
		},
		CreatedAt: time.Now().AddDate(0, -1, 0), // 一个月前创建
		UpdatedAt: time.Now(),
	}

	r.log.WithContext(ctx).Infof("Document retrieved: %s", documentID)
	return doc, nil
}

// DeleteDocument deletes document
func (r *documentRepo) DeleteDocument(ctx context.Context, documentID string, options *v1.DeleteOptions) error {
	r.log.WithContext(ctx).Infof("Deleting document: %s", documentID)

	// TODO: 调用 docstore 服务删除文档
	// 这里是模拟实现，实际应该通过 gRPC 调用 docstore 服务

	if options != nil {
		if options.DeleteRelatedChunks {
			r.log.WithContext(ctx).Info("Deleting related chunks")
		}
		if options.DeleteEmbeddings {
			r.log.WithContext(ctx).Info("Deleting embeddings")
		}
	}

	r.log.WithContext(ctx).Infof("Document deleted successfully: %s", documentID)
	return nil
}

// ListDocuments retrieves a list of documents
func (r *documentRepo) ListDocuments(ctx context.Context, pagination *commonv1.PaginationRequest, filters []*commonv1.Filter) ([]*biz.Document, *commonv1.PaginationResponse, error) {
	r.log.WithContext(ctx).Info("Listing documents")

	// TODO: 调用 docstore 服务获取文档列表
	// 这里是模拟实现，实际应该通过 gRPC 调用 docstore 服务

	// 模拟文档列表
	docs := []*biz.Document{
		{
			DocumentID:  "doc-1",
			Title:       "机器学习基础",
			FileType:    "pdf",
			FileSize:    2048000,
			TotalChunks: 15,
			Metadata: &commonv1.Metadata{
				Data: map[string]string{
					"author":   "张三",
					"category": "AI",
				},
			},
			CreatedAt: time.Now().AddDate(0, -2, 0),
			UpdatedAt: time.Now().AddDate(0, -1, 0),
		},
		{
			DocumentID:  "doc-2",
			Title:       "深度学习实践",
			FileType:    "pdf",
			FileSize:    3072000,
			TotalChunks: 20,
			Metadata: &commonv1.Metadata{
				Data: map[string]string{
					"author":   "李四",
					"category": "AI",
				},
			},
			CreatedAt: time.Now().AddDate(0, -1, -15),
			UpdatedAt: time.Now().AddDate(0, 0, -5),
		},
	}

	// 应用过滤器
	if len(filters) > 0 {
		var filteredDocs []*biz.Document
		for _, doc := range docs {
			passAllFilters := true
			for _, filter := range filters {
				if !r.applyFilter(doc, filter) {
					passAllFilters = false
					break
				}
			}
			if passAllFilters {
				filteredDocs = append(filteredDocs, doc)
			}
		}
		docs = filteredDocs
	}

	// 应用分页
	page := int32(1)
	pageSize := int32(10)
	if pagination != nil {
		if pagination.Page > 0 {
			page = pagination.Page
		}
		if pagination.PageSize > 0 {
			pageSize = pagination.PageSize
		}
	}

	total := int64(len(docs))
	totalPages := int32((total + int64(pageSize) - 1) / int64(pageSize))

	// 计算分页范围
	start := (page - 1) * pageSize
	end := start + pageSize
	if end > int32(len(docs)) {
		end = int32(len(docs))
	}
	if start > int32(len(docs)) {
		start = int32(len(docs))
	}

	pagedDocs := docs[start:end]

	paginationResp := &commonv1.PaginationResponse{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}

	r.log.WithContext(ctx).Infof("Listed %d documents (page %d/%d)", len(pagedDocs), page, totalPages)
	return pagedDocs, paginationResp, nil
}

// applyFilter applies a filter to a document
func (r *documentRepo) applyFilter(doc *biz.Document, filter *commonv1.Filter) bool {
	switch filter.Field {
	case "title":
		return r.matchStringFilter(doc.Title, filter)
	case "file_type":
		return r.matchStringFilter(doc.FileType, filter)
	case "author":
		if doc.Metadata != nil && doc.Metadata.Data != nil {
			if author, exists := doc.Metadata.Data["author"]; exists {
				return r.matchStringFilter(author, filter)
			}
		}
		return false
	case "category":
		if doc.Metadata != nil && doc.Metadata.Data != nil {
			if category, exists := doc.Metadata.Data["category"]; exists {
				return r.matchStringFilter(category, filter)
			}
		}
		return false
	default:
		return true // 未知字段默认通过
	}
}

// matchStringFilter matches string value against filter
func (r *documentRepo) matchStringFilter(value string, filter *commonv1.Filter) bool {
	if len(filter.Values) == 0 {
		return true
	}

	switch filter.Operator {
	case "eq":
		for _, filterValue := range filter.Values {
			if value == filterValue {
				return true
			}
		}
		return false
	case "ne":
		for _, filterValue := range filter.Values {
			if value == filterValue {
				return false
			}
		}
		return true
	case "like":
		// 简单的包含匹配
		for _, filterValue := range filter.Values {
			if len(filterValue) > 0 && len(value) > 0 {
				// 简单实现：检查是否包含子字符串
				for i := 0; i <= len(value)-len(filterValue); i++ {
					if value[i:i+len(filterValue)] == filterValue {
						return true
					}
				}
			}
		}
		return false
	case "in":
		for _, filterValue := range filter.Values {
			if value == filterValue {
				return true
			}
		}
		return false
	default:
		return true
	}
}

// UpdateDocumentMetadata updates document metadata
func (r *documentRepo) UpdateDocumentMetadata(ctx context.Context, documentID string, metadata *commonv1.Metadata, updateFields []string) error {
	r.log.WithContext(ctx).Infof("Updating document metadata: %s", documentID)

	// TODO: 调用 docstore 服务更新文档元数据
	// 这里是模拟实现，实际应该通过 gRPC 调用 docstore 服务

	r.log.WithContext(ctx).Infof("Document metadata updated successfully: %s", documentID)
	return nil
}

// GetDocumentStats retrieves document statistics
func (r *documentRepo) GetDocumentStats(ctx context.Context, documentID string) (*biz.DocumentStats, error) {
	r.log.WithContext(ctx).Infof("Getting document stats: %s", documentID)

	// TODO: 调用 docstore 服务获取文档统计
	// 这里是模拟实现，实际应该通过 gRPC 调用 docstore 服务

	stats := &biz.DocumentStats{
		TotalChunks:      15,
		TotalTokens:      5000,
		AverageChunkSize: 333,
		QueryCount:       42,
		LastAccessed:     time.Now().Add(-time.Hour * 2),
	}

	r.log.WithContext(ctx).Infof("Document stats retrieved: %s", documentID)
	return stats, nil
}

// CallDocstoreService calls docstore service
func (r *documentRepo) CallDocstoreService(ctx context.Context, operation string, data interface{}) (interface{}, error) {
	r.log.WithContext(ctx).Infof("Calling docstore service: %s", operation)

	// TODO: 实现对 docstore 服务的 gRPC 调用
	// 这里是模拟实现

	switch operation {
	case "save":
		return "document_saved", nil
	case "get":
		return "document_data", nil
	case "delete":
		return "document_deleted", nil
	default:
		return nil, fmt.Errorf("unknown operation: %s", operation)
	}
}

// CallPreprocessorService calls preprocessor service
func (r *documentRepo) CallPreprocessorService(ctx context.Context, content []byte, options *v1.DocumentProcessingOptions) (*biz.DocumentProcessingProgress, error) {
	r.log.WithContext(ctx).Info("Calling preprocessor service")

	// TODO: 实现对 preprocessor 服务的 gRPC 调用
	// 这里是模拟实现

	progress := &biz.DocumentProcessingProgress{
		ProgressPercentage:  25.0,
		CurrentStage:        "parsing",
		StatusMessage:       "文档解析中...",
		StartedAt:           time.Now(),
		EstimatedCompletion: time.Now().Add(time.Minute * 5),
	}

	r.log.WithContext(ctx).Info("Preprocessor service called successfully")
	return progress, nil
}

// CallEmbeddingService calls embedding service
func (r *documentRepo) CallEmbeddingService(ctx context.Context, documentID string, chunks []string) error {
	r.log.WithContext(ctx).Infof("Calling embedding service for document: %s", documentID)

	// TODO: 实现对 embedding 服务的 gRPC 调用
	// 这里是模拟实现

	r.log.WithContext(ctx).Infof("Embedding service called successfully for %d chunks", len(chunks))
	return nil
}
