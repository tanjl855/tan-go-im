package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// Md5Encode 小写
func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	tmpByte := h.Sum(nil)
	return hex.EncodeToString(tmpByte)
}

// MD5ToUpper 大写md5
func MD5ToUpper(data string) string {
	return strings.ToUpper(Md5Encode(data))
}

// MakePassword 加密
func MakePassword(plainpwd, salt string) string {
	return Md5Encode(plainpwd + salt)
}
