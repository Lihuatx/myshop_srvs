package chat_tools

import (
	"regexp"
	"strings"
)

// 精准匹配手机号（11位数字，以1开头）
var phoneRegex = regexp.MustCompile(`\b1[3-9]\d{9}\b`)

// 精准匹配地址（避免误伤“订单信息”等词汇）
var addressRegex = regexp.MustCompile(`(\p{Han}{2})[\p{Han}]+(街道|路|号|小区)`)

// 白名单：不处理的词汇（如“订单信息”）
var whitelist = []string{"订单信息", "订单号", "物流单号"}

func SanitizeInput(input string) string {
	// 白名单保护
	for _, word := range whitelist {
		input = strings.ReplaceAll(input, word, "{"+word+"}") // 临时占位
	}

	// 脱敏手机号（138****1234）
	masked := phoneRegex.ReplaceAllString(input, "$1****$2")

	// 脱敏地址（北京市***朝阳区）
	masked = addressRegex.ReplaceAllString(masked, "$1***$2")

	// 恢复白名单词汇
	for _, word := range whitelist {
		masked = strings.ReplaceAll(masked, "{"+word+"}", word)
	}

	return masked
}

// 匹配订单号（假设格式为字母+数字，6-10位，如 "ABC123"）
var orderSNRegex = regexp.MustCompile(`([A-Za-z0-9]{1,10})`)

func ExtractOrderSN(question string) string {
	matches := orderSNRegex.FindStringSubmatch(question)
	if len(matches) > 0 {
		return matches[0]
	}
	return ""
}
