package basic

import (
	"context"
	"encoding/json"
)

// ToBeautifulJson 结构化序列化
func ToBeautifulJson(ctx context.Context, data interface{}) string {
	marshal, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return ""
	}
	return string(marshal)
}
