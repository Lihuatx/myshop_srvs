package handler

//
//import (
//	"context"
//	"encoding/json"
//	"fmt"
//	"github.com/milvus-io/milvus-sdk-go/v2/client"
//	"github.com/milvus-io/milvus-sdk-go/v2/entity"
//	"github.com/sashabaranov/go-openai"
//	"strings"
//)
//
//// 知识库条目结构
//type KnowledgeEntry struct {
//	ID        int64     `json:"id"`
//	Content   string    `json:"content"`
//	Type      string    `json:"type"`
//	Source    string    `json:"source"`
//	Embedding []float32 `json:"-"`
//}
//
//// 商品知识构建器
//type ProductKnowledgeBuilder struct {
//	goods    *Goods
//	category *Category
//	brand    *Brands
//}
//
//func NewProductKnowledgeBuilder(goods *Goods, category *Category, brand *Brands) *ProductKnowledgeBuilder {
//	return &ProductKnowledgeBuilder{
//		goods:    goods,
//		category: category,
//		brand:    brand,
//	}
//}
//
//// 构建商品知识条目
//func (b *ProductKnowledgeBuilder) Build() []KnowledgeEntry {
//	entries := make([]KnowledgeEntry, 0)
//
//	// 基础信息
//	basicInfo := fmt.Sprintf("商品名称：%s，商品编号：%s，品类：%s，品牌：%s",
//		b.goods.Name, b.goods.GoodsSn, b.category.Name, b.brand.Name)
//	entries = append(entries, KnowledgeEntry{
//		ID:      int64(b.goods.ID),
//		Content: basicInfo,
//		Type:    "basic_info",
//		Source:  "product",
//	})
//
//	// 价格信息
//	priceInfo := fmt.Sprintf("商品 %s 的市场价格为 %.2f 元，商城价格为 %.2f 元",
//		b.goods.Name, b.goods.MarketPrice, b.goods.ShopPrice)
//	entries = append(entries, KnowledgeEntry{
//		ID:      int64(b.goods.ID),
//		Content: priceInfo,
//		Type:    "price_info",
//		Source:  "product",
//	})
//
//	// 商品描述
//	if b.goods.GoodsBrief != "" {
//		entries = append(entries, KnowledgeEntry{
//			ID:      int64(b.goods.ID),
//			Content: b.goods.GoodsBrief,
//			Type:    "description",
//			Source:  "product",
//		})
//	}
//
//	return entries
//}
//
//// 向量数据库管理器
//type VectorStore struct {
//	milvusClient client.Client
//	collection   string
//	dim          int64
//}
//
//func NewVectorStore(addr string, collection string, dim int64) (*VectorStore, error) {
//	ctx := context.Background()
//	c, err := client.NewClient(ctx, client.Config{
//		Address: addr,
//	})
//	if err != nil {
//		return nil, fmt.Errorf("连接Milvus失败: %v", err)
//	}
//
//	return &VectorStore{
//		milvusClient: c,
//		collection:   collection,
//		dim:          dim,
//	}, nil
//}
//
//// 创建集合
//func (vs *VectorStore) CreateCollection() error {
//	ctx := context.Background()
//	schema := &entity.Schema{
//		CollectionName: vs.collection,
//		Fields: []*entity.Field{
//			{
//				Name:       "id",
//				DataType:   entity.FieldTypeInt64,
//				PrimaryKey: true,
//			},
//			{
//				Name:      "content",
//				DataType:  entity.FieldTypeVarChar,
//				MaxLength: 2048,
//			},
//			{
//				Name:      "type",
//				DataType:  entity.FieldTypeVarChar,
//				MaxLength: 64,
//			},
//			{
//				Name:      "source",
//				DataType:  entity.FieldTypeVarChar,
//				MaxLength: 64,
//			},
//			{
//				Name:     "embedding",
//				DataType: entity.FieldTypeFloatVector,
//				TypeParams: map[string]string{
//					"dim": fmt.Sprintf("%d", vs.dim),
//				},
//			},
//		},
//	}
//
//	err := vs.milvusClient.CreateCollection(ctx, schema, entity.DefaultShardNumber)
//	if err != nil {
//		return fmt.Errorf("创建集合失败: %v", err)
//	}
//
//	// 创建索引
//	idx, err := entity.NewIndexIvfFlat(entity.L2)
//	if err != nil {
//		return fmt.Errorf("创建索引失败: %v", err)
//	}
//
//	err = vs.milvusClient.CreateIndex(ctx, vs.collection, "embedding", idx, false)
//	if err != nil {
//		return fmt.Errorf("创建索引失败: %v", err)
//	}
//
//	return nil
//}
//
//// RAG问答系统
//type RAGChat struct {
//	vectorStore  *VectorStore
//	openaiClient *openai.Client
//	embedDim     int
//}
//
//func NewRAGChat(vectorStore *VectorStore, openaiKey string) *RAGChat {
//	return &RAGChat{
//		vectorStore:  vectorStore,
//		openaiClient: openai.NewClient(openaiKey),
//		embedDim:     1536, // OpenAI embedding 维度
//	}
//}
//
//// 生成文本嵌入向量
//func (rc *RAGChat) generateEmbedding(text string) ([]float32, error) {
//	resp, err := rc.openaiClient.CreateEmbeddings(
//		context.Background(),
//		openai.EmbeddingRequest{
//			Input: []string{text},
//			Model: openai.AdaEmbeddingV2,
//		},
//	)
//	if err != nil {
//		return nil, err
//	}
//
//	// 转换为float32
//	embedding := make([]float32, len(resp.Data[0].Embedding))
//	for i, v := range resp.Data[0].Embedding {
//		embedding[i] = float32(v)
//	}
//
//	return embedding, nil
//}
//
//// 检索相关文档
//func (rc *RAGChat) searchSimilarDocs(embedding []float32, limit int) ([]KnowledgeEntry, error) {
//	ctx := context.Background()
//
//	searchParams, err := entity.NewIndexIvfFlatSearchParam(10)
//	if err != nil {
//		return nil, err
//	}
//
//	sp, err := rc.vectorStore.milvusClient.Search(
//		ctx,
//		rc.vectorStore.collection,
//		[]string{},
//		"", // 无表达式
//		[]string{"id", "content", "type", "source"},
//		[]entity.Vector{embedding},
//		"embedding",
//		entity.L2,
//		limit,
//		searchParams,
//	)
//	if err != nil {
//		return nil, err
//	}
//
//	results := make([]KnowledgeEntry, 0)
//	for _, result := range sp {
//		entry := KnowledgeEntry{
//			ID:      result.IDs.(*entity.ColumnInt64).Data()[0],
//			Content: result.Fields["content"].(*entity.ColumnVarChar).Data()[0],
//			Type:    result.Fields["type"].(*entity.ColumnVarChar).Data()[0],
//			Source:  result.Fields["source"].(*entity.ColumnVarChar).Data()[0],
//		}
//		results = append(results, entry)
//	}
//
//	return results, nil
//}
//
//// 处理用户问题
//func (rc *RAGChat) ProcessQuestion(ctx context.Context, question string) (string, error) {
//	// 1. 生成问题的嵌入向量
//	queryEmbedding, err := rc.generateEmbedding(question)
//	if err != nil {
//		return "", fmt.Errorf("生成问题向量失败: %v", err)
//	}
//
//	// 2. 检索相关文档
//	docs, err := rc.searchSimilarDocs(queryEmbedding, 3)
//	if err != nil {
//		return "", fmt.Errorf("检索相关文档失败: %v", err)
//	}
//
//	// 3. 构建 Prompt
//	var contextBuilder strings.Builder
//	contextBuilder.WriteString("以下是相关的商品信息：\n")
//	for _, doc := range docs {
//		contextBuilder.WriteString(fmt.Sprintf("- %s\n", doc.Content))
//	}
//
//	prompt := fmt.Sprintf(`基于以下信息回答用户问题:
//
//背景信息:
//%s
//
//用户问题:
//%s
//
//要求:
//1. 只使用提供的信息来回答
//2. 如果信息不足，说明无法回答
//3. 使用客服的语气，简洁友好
//
//回答：`, contextBuilder.String(), question)
//
//	// 4. 调用大模型生成回答
//	resp, err := rc.openaiClient.CreateChatCompletion(
//		ctx,
//		openai.ChatCompletionRequest{
//			Model: openai.GPT3Dot5Turbo,
//			Messages: []openai.ChatCompletionMessage{
//				{
//					Role:    openai.ChatMessageRoleUser,
//					Content: prompt,
//				},
//			},
//			Temperature: 0.7,
//		},
//	)
//	if err != nil {
//		return "", fmt.Errorf("生成回答失败: %v", err)
//	}
//
//	return resp.Choices[0].Message.Content, nil
//}

//how to use

//func main() {
//// 1. 初始化向量存储
//vectorStore, err := NewVectorStore("localhost:19530", "product_knowledge", 1536)
//if err != nil {
//log.Fatal(err)
//}
//
//// 2. 创建RAG聊天系统
//ragChat := NewRAGChat(vectorStore, "your-openai-key")
//
//// 3. 处理用户问题
//answer, err := ragChat.ProcessQuestion(context.Background(), "这个商品的价格是多少？")
//if err != nil {
//log.Fatal(err)
//}
//fmt.Println(answer)
//}
