package cms

import (
	"crypto/md5"
	"strings"
)

var cc = []byte("0123456789ABCDEF")
func base16(arr []byte) string {
	sb := strings.Builder{}
	for _, b := range arr {
		sb.WriteByte(cc[(b >> 4) & 0x0f])
		sb.WriteByte(cc[b & 0x0f])
	}
	return sb.String()
}

func md5Str(arr []byte) string {
	var bytes []byte
	for _, b := range md5.Sum(arr) {
		bytes = append(bytes, b)
	}
	return base16(bytes)
}
