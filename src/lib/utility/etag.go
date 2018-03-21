package utility

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"io"
	"os"
)

// HTTP 协议：ETag == URL 的 Entity Tag，用于标示 URL 对象是否改变，区分不同语言和 Session 等等。

func ComputeEtag(data []byte) (etag string) {
	tag := make([]byte, 0, 1+sha1.Size)
	h := sha1.New()
	if len(data) < _BLOCK_SIZE {
		tag = append(tag, 0x16)
		h.Write(data)
		tag = h.Sum(tag)
	} else {
		tag = append(tag, 0x96)
		block_count := blockCount(int64(len(data)))
		all_blocks_sha1 := make([]byte, 0, block_count*sha1.Size)
		for i := 0; i < block_count; i++ {
			var readBytes int
			if i < block_count-1 {
				readBytes = _BLOCK_SIZE
			} else {
				readBytes = len(data) - _BLOCK_SIZE*i
			}
			h.Write(data[_BLOCK_SIZE*i : _BLOCK_SIZE*i+readBytes])
			all_blocks_sha1 = h.Sum(all_blocks_sha1)
			h.Reset()
		}
		h.Write(all_blocks_sha1)
		tag = h.Sum(tag)
	}
	return base64.URLEncoding.EncodeToString(tag)
}

func ComputeFileEtag(filename string) (etag string, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return
	}

	fsize := fi.Size()
	block_count := blockCount(fsize)
	tag := []byte{}

	if block_count <= 1 { // file size <= 4M
		tag, err = computeSha1([]byte{0x16}, f)
		if err != nil {
			return
		}
	} else { // file size > 4M
		all_blocks_sha1 := []byte{}

		for i := 0; i < block_count; i++ {
			body := io.LimitReader(f, _BLOCK_SIZE)
			all_blocks_sha1, err = computeSha1(all_blocks_sha1, body)
			if err != nil {
				return
			}
		}

		tag, _ = computeSha1([]byte{0x96}, bytes.NewReader(all_blocks_sha1))
	}

	etag = base64.URLEncoding.EncodeToString(tag)
	return
}

// private:
const (
	_BLOCK_BITS = 22               // 2 ^ 22 = 4M
	_BLOCK_SIZE = 1 << _BLOCK_BITS // 4M
)

func blockCount(size int64) int {
	return int((size + (_BLOCK_SIZE - 1)) >> _BLOCK_BITS)
}

func computeSha1(b []byte, r io.Reader) ([]byte, error) {
	h := sha1.New()
	_, err := io.Copy(h, r)
	if err != nil {
		return nil, err
	}
	return h.Sum(b), nil
}
