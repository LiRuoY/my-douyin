package util

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

func Equal(password, md5Password string) bool {
	return strings.EqualFold(Encrypted(password), md5Password)
}
func Encrypted(password string) string {
	h := md5.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}
