package main

//加密解密通用函数
import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"time"
)

//加密
func base64Encode(src []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(src))
}

//解密
func base64Decode(src []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(src))
}

func md5Encode(src []byte) []byte {
	h := md5.New()
	h.Write(src)
	return []byte(hex.EncodeToString(h.Sum(nil)))
}

//时间
func CreatedAt() string {
	return time.Now().Format(TIME_FORMAT)
}
