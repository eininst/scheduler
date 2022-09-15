package util

import (
	"crypto/md5"
	"encoding/hex"
	"time"
)

func Md5(value string) string {
	m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
}
func FormatTime() string {
	return time.Now().Format("2006.01.02 15:04:05")
}
