package biz

import (
	"context"
	"time"

	commonv1 "rag/api/common/v1"
	v1 "rag/api/gateway/v1"

	"github.com/go-kratos/kratos/v2/log"
)

// QueryResult represents query processing result
type QueryResult struct {
	QueryID          string
	Answer           string
	RelatedDocuments []*RelatedDocument
	Metadata         *QueryMetadata
	Suggestions      []string
}

// RelatedDocument represents a document related to the query
type RelatedDocument struct {
	DocumentID     string
	Title          string
	Snippet        string
	RelevanceScore float32
	DocumentType   string
	Chunks         []*commonv1.ChunkInfo
}

// QueryMetadata contains metadata about query processing
type QueryMetadata struct {
	QueryTime              time.Time
	ProcessingTimeMs       int64
	TotalDocumentsSearched int32
	DocumentsReturned      int32
	ModelUsed              string
	DebugInfo              map[string]string
}

// QueryRepo defines the data access interface for query processing
type QueryRepo interface {
	// 调用检索服务进行文档检索
	RetrieveDocuments(ctx context.Context, query string, params *v1.QueryParameters) ([]*RelatedDocument, error)
	// 调用重排序服务优化结果
	RerankDocuments(ctx context.Context, query string, documents []*RelatedDocument) ([]*RelatedDocument, error)
	// 调用组装服务生成答案
	AssembleAnswer(ctx context.Context, query string, documents []*RelatedDocument) (string, error)
	// 保存查询历史
	SaveQueryHistory(ctx context.Context, userID, query, answer string, metadata *QueryMetadata) error
	// 获取查询建议
	GetQuerySuggestions(ctx context.Context, query string) ([]string, error)
}

// QueryUsecase handles query processing business logic
type QueryUsecase struct {
	repo QueryRepo
	log  *log.Helper
}

// NewQueryUsecase creates a new query usecase
func NewQueryUsecase(repo QueryRepo, logger log.Logger) *QueryUsecase {
	return &QueryUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// ProcessQuery handles intelligent query processing
func (uc *QueryUsecase) ProcessQuery(ctx context.Context, req *v1.QueryRequest) (*v1.QueryResponse, error) {
	startTime := time.Now()
	uc.log.WithContext(ctx).Infof("Processing query: %s", req.Query)

	// 检索相关文档
	documents, err := uc.repo.RetrieveDocuments(ctx, req.Query, req.Parameters)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("Failed to retrieve documents: %v", err)
		return nil, err
	}

	uc.log.WithContext(ctx).Infof("Retrieved %d documents", len(documents))

	// 如果启用重排序，调用重排序服务
	if req.Parameters != nil && req.Parameters.EnableReranking {
		rerankedDocs, err := uc.repo.RerankDocuments(ctx, req.Query, documents)
		if err != nil {
			uc.log.WithContext(ctx).Warnf("Failed to rerank documents: %v", err)
		} else {
			documents = rerankedDocs
			uc.log.WithContext(ctx).Info("Documents reranked successfully")
		}
	}

	// 如果启用上下文组装，调用组装服务生成答案
	var answer string
	if req.Parameters != nil && req.Parameters.EnableContextAssembly {
		generatedAnswer, err := uc.repo.AssembleAnswer(ctx, req.Query, documents)
		if err != nil {
			uc.log.WithContext(ctx).Warnf("Failed to assemble answer: %v", err)
			answer = "抱歉，无法生成答案，请查看相关文档。"
		} else {
			answer = generatedAnswer
		}
	} else {
		answer = "请查看相关文档获取信息。"
	}

	// 获取查询建议
	suggestions, err := uc.repo.GetQuerySuggestions(ctx, req.Query)
	if err != nil {
		uc.log.WithContext(ctx).Warnf("Failed to get suggestions: %v", err)
		suggestions = []string{}
	}

	// 构建查询元数据
	processingTime := time.Since(startTime).Milliseconds()
	metadata := &QueryMetadata{
		QueryTime:              startTime,
		ProcessingTimeMs:       processingTime,
		TotalDocumentsSearched: int32(len(documents)),
		DocumentsReturned:      int32(len(documents)),
		ModelUsed:              "default",
		DebugInfo:              make(map[string]string),
	}

	// 保存查询历史
	if err := uc.repo.SaveQueryHistory(ctx, req.UserId, req.Query, answer, metadata); err != nil {
		uc.log.WithContext(ctx).Warnf("Failed to save query history: %v", err)
	}

	// 转换为 protobuf 格式
	relatedDocs := make([]*v1.RelatedDocument, len(documents))
	for i, doc := range documents {
		relatedDocs[i] = &v1.RelatedDocument{
			DocumentId:     doc.DocumentID,
			Title:          doc.Title,
			Snippet:        doc.Snippet,
			RelevanceScore: doc.RelevanceScore,
			DocumentType:   doc.DocumentType,
			Chunks:         doc.Chunks,
		}
	}

	response := &v1.QueryResponse{
		QueryId:          generateQueryID(),
		Answer:           answer,
		RelatedDocuments: relatedDocs,
		Metadata: &v1.QueryMetadata{
			ProcessingTimeMs:       metadata.ProcessingTimeMs,
			TotalDocumentsSearched: metadata.TotalDocumentsSearched,
			DocumentsReturned:      metadata.DocumentsReturned,
			ModelUsed:              metadata.ModelUsed,
			DebugInfo:              metadata.DebugInfo,
		},
		Suggestions: suggestions,
	}

	uc.log.WithContext(ctx).Infof("Query processed successfully in %dms", processingTime)
	return response, nil
}

// generateQueryID generates a unique query ID
func generateQueryID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString generates a random string of specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
