package rag

import (
	"context"
	"fmt"
	"github.com/qdrant/go-client/qdrant"
	"go.uber.org/zap"
	"myshop_srvs/goods_srv/model"
	"myshop_srvs/smartChat_srv/global"
)

type VectorDB struct {
	client qdrant.Client
}

func syncGoodsToVectorDB() {
	var goods []model.Goods
	if err := global.DB.Preload("Category").Preload("Brands").Find(&goods).Error; err != nil {
		return
	}

	vdb := NewVectorDB()
	if err := vdb.StoreGoodsVectors(goods); err != nil {
		zap.S().Errorf("向量数据同步失败: %v", err)
	}
}

func NewVectorDB() *VectorDB {
	conn, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334,
	})
	if err != nil {
		panic(err)
	}
	return &VectorDB{client: *conn}
}

// 创建集合（只需运行一次）
func (v *VectorDB) CreateCollection() error {
	err := v.client.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: "goods_vectors",
		VectorsConfig: &qdrant.VectorsConfig{
			Config: &qdrant.VectorsConfig_Params{
				Params: &qdrant.VectorParams{
					Size:     384, // 根据嵌入模型维度设置
					Distance: qdrant.Distance_Cosine,
				},
			},
		},
	})

	return err
}

// 商品数据向量化存储
func (v *VectorDB) StoreGoodsVectors(goods []model.Goods) error {
	// 获取嵌入模型
	embeddings := NewEmbeddingModel()

	points := make([]*qdrant.PointStruct, 0, len(goods))
	for _, g := range goods {
		// 构造知识文本块
		text := fmt.Sprintf(`商品名称：%s
简介：%s
分类：%s
品牌：%s
价格：%.2f
是否在售：%v`,
			g.Name, g.GoodsBrief,
			getCategoryName(g.CategoryID),
			getBrandName(g.BrandsID),
			g.ShopPrice, g.OnSale)

		// 生成嵌入向量
		vector, err := embeddings.GetEmbedding(text)
		if err != nil {
			continue
		}

		points = append(points, &qdrant.PointStruct{
			Id: &qdrant.PointId{
				PointIdOptions: &qdrant.PointId_Num{
					Num: uint64(g.ID),
				},
			},
			Vectors: &qdrant.Vectors{
				VectorsOptions: &qdrant.Vectors_Vector{
					Vector: &qdrant.Vector{
						Data: vector,
					},
				},
			},
			Payload: map[string]*qdrant.Value{
				"text": {Kind: &qdrant.Value_StringValue{StringValue: text}},
				"id":   {Kind: &qdrant.Value_IntegerValue{IntegerValue: int64(g.ID)}},
			},
		})
	}

	_, err := v.client.Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: "goods_vectors",
		Points:         points,
	})
	return err
}

// 辅助函数：获取分类名称（需实现）
func getCategoryName(id int32) string {
	// 查询分类服务或本地缓存
	return ""
}

func getBrandName(id int32) string {
	// 查询分类服务或本地缓存
	return ""
}
