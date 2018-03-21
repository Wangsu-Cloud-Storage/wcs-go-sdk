package utility

import (
	"testing"
)

func TestComputeEtag(t *testing.T) {
	etag := ComputeEtag([]byte("UMUTech@qq.com"))
	if "FpAEzHGOWdh-QAYhattxl1vqHLMw" != etag {
		t.Fatal(etag, "Should be: FpAEzHGOWdh-QAYhattxl1vqHLMw")
	}

	const blockSize int64 = 4 * 1024 * 1024
	if 2 != blockCount(blockSize+1) {
		t.Fatal(etag, "Should be: 2")
	}

	blockData := make([]byte, blockSize+1)
	var i int64 = 0
	for ; i < blockSize/2; i++ {
		blockData[i] = 85
	}
	for ; i < blockSize; i++ {
		blockData[i] = 77
	}
	blockData[i] = '0'
	etag = ComputeEtag(blockData)
	if "ljQoZqmvLKtOnK_5OHnMRw1djggC" != etag {
		t.Fatal(etag, "Should be: ljQoZqmvLKtOnK_5OHnMRw1djggC")
	}
}

func TestComputeFileEtag(t *testing.T) {
	etag, err := ComputeFileEtag("C:\\bootmgr")
	t.Log(etag, err)
}
