package basic

import (
	"reflect"
	"strings"
)

func GetStructName(instance any) string {
	// 获取实例的类型字符串
	typeStr := reflect.TypeOf(instance).String()

	// 处理指针类型（去掉开头的 *）
	if strings.HasPrefix(typeStr, "*") {
		typeStr = typeStr[1:]
	}

	// 处理泛型类型（去掉泛型参数部分）
	if idx := strings.Index(typeStr, "["); idx != -1 {
		typeStr = typeStr[:idx]
	}

	// 提取结构体名称
	parts := strings.Split(typeStr, ".")
	if len(parts) == 0 {
		return "" // 防止空字符串
	}
	return parts[len(parts)-1]
}

func Zero[T any]() (v T) {
	return
}
