package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"runtime"

	"github.com/google/uuid"
)

// CreateUUID 初始化时使用uuid作为salt
func CreateUUID() string {
	uuid := uuid.New()
	return uuid.String()
}

// GetMd5
func GetMd5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	s := h.Sum(nil)
	return hex.EncodeToString(s)
}

// SessionFile login session file
func SessionPath() string {
	var homeDir string
	var err error
	switch runtime.GOOS {
	case "windows":
		homeDir = os.Getenv("USERPROFILE")
	case "darwin":
		homeDir = os.Getenv("HOME")
	default:
		homeDir = os.Getenv("HOME")
	}

	if len(homeDir) <= 0 {
		homeDir, err = os.UserCacheDir()
		if err != nil {
			homeDir, err = os.Getwd()
			if err != nil {
				fmt.Println("Error getting current working directory:", err)
				os.Exit(1)
			}
		}
	}
	return homeDir
}

// CheckFileExist 检查DB文件是否存在
func CheckFileExist(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}
