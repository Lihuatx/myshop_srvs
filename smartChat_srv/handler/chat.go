package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"time"

	"myshop_srvs/smartChat_srv/global"
	"myshop_srvs/smartChat_srv/proto"
	"myshop_srvs/smartChat_srv/utils/chat_tools"
)

type ChatServer struct {
	proto.UnimplementedChatServiceServer
}

func (s *ChatServer) ProcessMessage(ctx context.Context, req *proto.ChatRequest) (*proto.ChatResponse, error) {
	// 1. 敏感信息过滤
	sanitizedQuestion := chat_tools.SanitizeInput(req.Question)
	sanitizedQuestion = req.Question
	// 2. 意图分类（调用API）
	intent, err := ClassifyIntent(sanitizedQuestion)
	//intent := ClassifyByKeywords(sanitizedQuestion)

	if err != nil {
		zap.S().Errorf("意图分类失败: %v", err)
		intent = ClassifyByKeywords(sanitizedQuestion)
	}

	// 3. 根据意图调用其他服务（如订单服务）
	var answer string
	switch intent {
	case "order_query":
		// 优先使用前端传递的订单号，否则从问题中提取
		orderID := req.OrderId
		if orderID == 0 {
			tmp, _ := strconv.Atoi(chat_tools.ExtractOrderSN(sanitizedQuestion))
			orderID = int32(tmp)
		}

		if orderID != 0 {
			detailResp, err := global.OrderSrvClient.OrderDetail(ctx, &proto.OrderRequest{
				UserId: 1,
				Id:     orderID,
			})
			if err != nil || detailResp.OrderInfo == nil {
				answer = "订单查询失败，请稍后再试"
			} else {
				// 使用大模型生成回答
				answer, err = generateOrderAnswer(sanitizedQuestion, detailResp.OrderInfo)
				if err != nil {
					answer = fmt.Sprintf("订单 %s 状态：%s",
						detailResp.OrderInfo.OrderSn,
						detailResp.OrderInfo.Status)
				}
			}
		} else {
			// 未提供订单号，查询订单列表并返回最近一条
			listResp, err := global.OrderSrvClient.OrderList(ctx, &proto.OrderFilterRequest{
				//UserId: req.UserId,
				UserId: 1,
			})
			if err != nil || len(listResp.Data) == 0 {
				answer = "您当前没有订单"
			} else {
				// 取最近一条订单
				latestOrder := listResp.Data[0]
				answer = fmt.Sprintf("您最近的订单 %s 状态：%s（如需详情，请提供订单号）",
					latestOrder.OrderSn,
					latestOrder.Status,
				)
			}
		}
	default:
		// 调用大模型生成回答（DeepSeek API）
		generatedAnswer, err := GenerateAnswer(sanitizedQuestion)
		if err != nil {
			return nil, fmt.Errorf("生成回答失败: %v", err)
		}
		answer = generatedAnswer
	}

	// 4. 记录到数据库
	//log := model.ChatLog{
	//	UserID:   req.UserId,
	//	Question: req.Question,
	//	Answer:   answer,
	//	Intent:   intent,
	//}
	//global.DB.Create(&log)

	return &proto.ChatResponse{Answer: answer, Intent: intent}, nil
}

// 定义支持的分类标签
var allowedIntents = []string{"order_query", "refund", "logistics", "complaint", "other"}

// ClassifyIntent 使用 DeepSeek API 进行意图分类
func ClassifyIntent(question string) (string, error) {
	// 从配置获取 DeepSeek API 密钥
	apiKey := global.ServerConfig.DeepSeepKey.Key
	if apiKey == "" {
		return "other", fmt.Errorf("DeepSeek API 密钥未配置")
	}

	// 设计分类 Prompt
	prompt := fmt.Sprintf(`请对以下用户问题进行意图分类，可选标签：%s。直接返回标签，不要解释。
用户问题：%s`, strings.Join(allowedIntents, "、"), question)

	// 构建 API 请求
	url := "https://api.deepseek.com/v1/chat/completions"
	requestBody := map[string]interface{}{
		"model": "deepseek-chat",
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"temperature": 0.3, // 降低随机性，提高稳定性
		"max_tokens":  10,  // 限制返回长度
	}
	jsonData, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "other", fmt.Errorf("API 请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "other", fmt.Errorf("解析响应失败: %v", err)
	}

	// 提取分类标签
	if len(result.Choices) == 0 {
		return "other", fmt.Errorf("API 未返回有效结果")
	}
	answer := strings.TrimSpace(result.Choices[0].Message.Content)

	// 验证标签是否合法
	for _, intent := range allowedIntents {
		if strings.EqualFold(answer, intent) {
			return intent, nil
		}
	}
	return "other", fmt.Errorf("无效的分类标签: %s", answer)
}

// 降级策略：关键词匹配
func ClassifyByKeywords(question string) string {
	keywords := map[string]string{
		"订单": "order_query",
		"退货": "refund",
		"物流": "logistics",
		"投诉": "complaint",
	}
	for kw, intent := range keywords {
		if strings.Contains(question, kw) {
			fmt.Println("命中关键词：", kw)
			return intent
		}
	}
	return "other"
}

func generateOrderAnswer(question string, order *proto.OrderInfoResponse) (string, error) {

	// 构建Prompt模板
	prompt := fmt.Sprintf(`根据用户问题和订单信息生成回答：
<用户问题>
%s
</用户问题>

<订单数据>
- 订单号：%s
- 状态：%s
- 下单时间：%s
</订单数据>

要求：
1. 直接回答问题，不要复述数据。
2. 用口语化中文，不超过50字。
3. 如果用户未明确问物流单号，不要主动透露。

示例：
用户问题："我的订单到哪了？"
回答："您的订单 %s 已发货，物流预计明天送达，请耐心等待！"`,
		question,
		order.OrderSn,
		order.Status,
		order.AddTime,
		order.OrderSn,
	)

	// 调用大模型API
	answer, err := GenerateAnswer(prompt)
	if err != nil {
		// 降级为默认模板
		return fmt.Sprintf("订单 %s 状态：%s", order.OrderSn, order.Status), nil
	}
	return answer, nil
}
