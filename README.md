# RAG微服务架构实现

基于设计文档实现的完整RAG（检索增强生成）微服务架构，采用Go语言和Kratos框架构建。

## 架构概览

本系统由7个核心微服务组成：

- **Gateway服务**: 统一入口、用户认证、请求路由和文档管理API
- **Orchestrator服务**: 工作流编排和服务协调
- **Preprocessor服务**: 查询预处理、文本清洗和改写
- **Embedding服务**: 文本向量化和相似度计算
- **DocStore服务**: 文档管理、分片存储和向量检索
- **Reranker服务**: 文档重排序和相关性评分
- **Assembler服务**: 上下文构建和Token管理

## 项目结构

```
.
├── api/                    # Protocol Buffers API定义
│   ├── common/v1/          # 通用消息类型和错误处理
│   ├── gateway/v1/         # Gateway服务API
│   ├── preprocessor/v1/    # Preprocessor服务API
│   ├── embedding/v1/       # Embedding服务API
│   ├── docstore/v1/        # DocStore服务API
│   ├── reranker/v1/        # Reranker服务API
│   ├── assembler/v1/       # Assembler服务API
│   └── orchestrator/v1/    # Orchestrator服务API
├── app/                    # 微服务应用
│   ├── gateway/            # Gateway服务实现
│   ├── preprocessor/       # Preprocessor服务实现
│   ├── embedding/          # Embedding服务实现
│   ├── docstore/           # DocStore服务实现
│   ├── reranker/           # Reranker服务实现
│   ├── assembler/          # Assembler服务实现
│   └── orchestrator/       # Orchestrator服务实现
├── third_party/            # 第三方Proto文件
├── monitoring/             # 监控配置
├── scripts/                # 数据库初始化脚本
├── docker-compose.yml      # Docker编排配置
└── Makefile               # 构建脚本
```

## 快速开始

### 前置要求

- Go 1.21+
- Docker & Docker Compose
- Protocol Buffers compiler
- Make

### 环境初始化

```bash
# 初始化开发环境
make init

# 生成所有API代码
make api
```

### 使用Docker Compose启动

```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f [service_name]
```

### 服务端口

| 服务 | HTTP端口 | gRPC端口 | 描述 |
|------|----------|----------|------|
| Gateway | 8080 | 9000 | 统一入口 |
| Orchestrator | 8081 | 9001 | 工作流编排 |
| Preprocessor | 8082 | 9002 | 查询预处理 |
| Embedding | 8083 | 9003 | 文本向量化 |
| DocStore | 8084 | 9004 | 文档存储 |
| Reranker | 8085 | 9005 | 文档重排序 |
| Assembler | 8086 | 9006 | 上下文构建 |

### 基础设施端口

| 服务 | 端口 | 描述 |
|------|------|------|
| PostgreSQL | 5432 | 元数据存储 |
| Milvus | 19530 | 向量数据库 |
| Redis | 6379 | 缓存 |
| Prometheus | 9090 | 监控 |
| Grafana | 3000 | 可视化 |
| Jaeger | 16686 | 链路追踪 |

## API文档

每个服务都会生成对应的Swagger文档：

- Gateway: `api/gateway/v1/gateway_svc.swagger.json`
- Preprocessor: `api/preprocessor/v1/preprocessor_svc.swagger.json`
- Embedding: `api/embedding/v1/embedding_svc.swagger.json`
- DocStore: `api/docstore/v1/docstore_svc.swagger.json`
- Reranker: `api/reranker/v1/reranker_svc.swagger.json`
- Assembler: `api/assembler/v1/assembler_svc.swagger.json`
- Orchestrator: `api/orchestrator/v1/orchestrator_svc.swagger.json`

## 开发

### 构建单个服务

```bash
cd app/[service_name]
make build
```

### 运行单个服务

```bash
cd app/[service_name]
make run
```

### 重新生成API代码

```bash
# 生成指定服务的API代码
cd app/[service_name]
make api

# 或者生成所有服务的API代码
make api
```

## 核心特性

### 完整的API定义
- 基于Protocol Buffers的类型安全API
- 支持gRPC和HTTP双协议
- 完整的错误处理和验证
- 自动生成的Swagger文档

### 微服务架构
- 服务职责清晰分离
- 支持独立部署和扩展
- 统一的配置和监控

### 文档处理流水线
- 多格式文档支持
- 智能分片策略
- 向量化和索引
- 混合检索能力

### 查询处理流程
- 查询预处理和改写
- 语义向量检索
- 文档重排序
- 智能上下文构建

### 工作流编排
- 灵活的工作流定义
- 支持条件执行和并行处理
- 完整的执行跟踪
- 错误处理和重试机制

## 监控和运维

### 健康检查
所有服务都提供健康检查端点：`/v1/health`

### 指标监控
- Prometheus指标收集
- Grafana可视化面板
- 服务性能监控

### 链路追踪
- Jaeger分布式追踪
- 请求链路分析
- 性能瓶颈定位

## 贡献指南

1. Fork项目
2. 创建功能分支
3. 提交更改
4. 创建Pull Request

## 许可证

本项目采用MIT许可证 - 查看[LICENSE](LICENSE)文件了解详情

