package utils

// ByteBuilder 拼接多个字节切片
// 使用 make + copy 方式，性能最优
func ByteBuilder(data ...[]byte) []byte {
	// 计算总长度
	totalLen := 0
	for _, b := range data {
		totalLen += len(b)
	}

	// 直接分配内存
	result := make([]byte, totalLen)

	// 复制数据
	pos := 0
	for _, b := range data {
		copy(result[pos:], b)
		pos += len(b)
	}

	return result
}
