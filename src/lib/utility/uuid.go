package utility

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

// 代替方案：go get -u github.com/satori/go.uuid

func GetUuid() (uuid string) {
	const uuidByteSize = 18
	b := make([]byte, uuidByteSize)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "554d5500-abcd-abcd-abcd-abcd1234abcd"
	}
	t := make([]byte, hex.EncodedLen(uuidByteSize))
	hex.Encode(t, b)
	t[8] = '-'
	t[13] = '-'
	t[18] = '-'
	t[23] = '-'
	uuid = string(t)
	return
}
