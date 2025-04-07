package utils

import (
	"fmt"
	"strconv"
)

const (
	VALUE_TYPE_INT16  = 1
	VALUE_TYPE_UINT16 = 2
	VALUE_TYPE_INT32  = 3
	VALUE_TYPE_UINT32 = 4
	VALUE_TYPE_INT64  = 5
	VALUE_TYPE_UINT64 = 6
	VALUE_TYPE_FLOAT  = 7
	VALUE_TYPE_DOUBLE = 8
	VALUE_TYPE_BIT    = 9
	VALUE_TYPE_STRING = 10
	VALUE_TYPE_BYTES  = 11
)

// BitCount 函数根据 ValueType 返回所需的字节数
func BitCount(val uint8) int {
	switch val {
	case VALUE_TYPE_INT16, VALUE_TYPE_UINT16, VALUE_TYPE_BIT:
		return 1
	case VALUE_TYPE_INT32, VALUE_TYPE_UINT32, VALUE_TYPE_FLOAT:
		return 2
	case VALUE_TYPE_INT64, VALUE_TYPE_UINT64, VALUE_TYPE_DOUBLE:
		return 4
	case VALUE_TYPE_STRING:
		return 1
	default:
		return 0
	}
}

// ApplyByteOrder 根据符号重新排列字节
func ApplyByteOrder(data []uint8, symbol string) []uint8 {
	switch symbol {
	case "B":
		// 2,1 for int16/uint16
		if len(data) >= 2 {
			return []uint8{data[1], data[0]}
		}
	case "L":
		// 1,2 for int16/uint16 (default big endian)
		return data
	case "LL":
		// 1,2,3,4 for int32/uint32/float32 (default big endian)
		return data
	case "LB":
		// 2,1,4,3 for int32/uint32/float32
		if len(data) >= 4 {
			return []uint8{data[1], data[0], data[3], data[2]}
		}
	case "BL":
		// 3,4,1,2 for int32/uint32/float32
		if len(data) >= 4 {
			return []uint8{data[2], data[3], data[0], data[1]}
		}
	case "BB":
		// 4,3,2,1 for int32/uint32/float32
		if len(data) >= 4 {
			return []uint8{data[3], data[2], data[1], data[0]}
		}
	default:
		return data // 默认不改变
	}
	return data
}

func StringToUint16(s string) (uint16, error) {
	// 将字符串解析为 int
	parsedInt, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	// 检查是否在 uint16 范围内
	if parsedInt < 0 || parsedInt > 65535 {
		return 0, fmt.Errorf("value out of range for uint16")
	}

	// 转换为 uint16 类型
	return uint16(parsedInt), nil
}
func StringToUint8(s string) (uint8, error) {
	// 将字符串解析为 int
	parsedInt, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	// 检查是否在 uint16 范围内
	if parsedInt < 0 || parsedInt > 254 {
		return 0, fmt.Errorf("value out of range for uint16")
	}

	// 转换为 uint16 类型
	return uint8(parsedInt), nil
}
