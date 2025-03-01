package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"myshop_srvs/smartChat_srv/global"
	"net/http"
	"time"
)

type DeepSeekResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func GenerateAnswer(question string) (string, error) {
	// 从Nacos配置获取API密钥
	apiKey := global.ServerConfig.DeepSeepKey.Key
	url := "https://api.deepseek.com/v1/chat/completions"

	prompt := fmt.Sprintf("用户问题：%s。请用中文简洁回答，不超过100字。", question)
	body := map[string]interface{}{
		"model": "deepseek-chat",
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"temperature": 0.7,
	}
	jsonData, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API请求失败: %v", err)
	}
	defer resp.Body.Close()

	var result DeepSeekResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	if len(result.Choices) > 0 {
		return result.Choices[0].Message.Content, nil
	}
	return "暂时无法回答此问题，请联系人工客服", nil
}
