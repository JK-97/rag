package data

import (
	"rag/app/gateway/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"google.golang.org/grpc"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewAuthRepo,
	NewQueryRepo,
	NewDocumentRepo,
)

// Data .
type Data struct {
	// gRPC 客户端连接
	docstoreConn     *grpc.ClientConn
	preprocessorConn *grpc.ClientConn
	embeddingConn    *grpc.ClientConn
	orchestratorConn *grpc.ClientConn
	rerankerConn     *grpc.ClientConn
	assemblerConn    *grpc.ClientConn

	// 其他依赖
	logger log.Logger
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	// TODO: 初始化各个微服务的 gRPC 连接
	// 这里先用空连接，实际应该根据配置文件建立连接

	data := &Data{
		logger: logger,
	}

	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
		// TODO: 关闭所有 gRPC 连接
		if data.docstoreConn != nil {
			data.docstoreConn.Close()
		}
		if data.preprocessorConn != nil {
			data.preprocessorConn.Close()
		}
		if data.embeddingConn != nil {
			data.embeddingConn.Close()
		}
		if data.orchestratorConn != nil {
			data.orchestratorConn.Close()
		}
		if data.rerankerConn != nil {
			data.rerankerConn.Close()
		}
		if data.assemblerConn != nil {
			data.assemblerConn.Close()
		}
	}

	return data, cleanup, nil
}
