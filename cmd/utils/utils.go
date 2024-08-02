package utils

import (
	"crypto/md5"
	"encoding/hex"
	"os"

	"github.com/google/uuid"
)

// CreateUUID 初始化时使用uuid作为salt
func CreateUUID() string {
	uuid := uuid.New()
	return uuid.String()
}

// GetMd5
func GetMd5(pwd string) string {
	h := md5.New()
	h.Write([]byte(pwd))
	s := h.Sum(nil)
	return hex.EncodeToString(s)
}

// CheckFileExist 检查DB文件是否存在
func CheckFileExist(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}
