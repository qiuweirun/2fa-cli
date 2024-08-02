package utils

import (
	"fmt"
	"strings"
	"testing"
)

func TestAesGCM(t *testing.T) {
	key := strings.Repeat("a", 16)
	data := "hello word!"
	// 加密
	gcm := AesEncryptGCM(data, key)
	fmt.Printf("密钥key: %s \n", key)
	fmt.Printf("加密数据: %s \n", data)
	fmt.Printf("加密结果: %s \n", gcm)
	// 解密
	str := AesDecryptGCM(gcm, key)
	fmt.Printf("解密结果: %s \n", str)
}
