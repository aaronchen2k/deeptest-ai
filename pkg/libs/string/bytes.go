package _str

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5Byte(p []byte) (string, error) {
	h := md5.New()
	_, err := h.Write(p)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
