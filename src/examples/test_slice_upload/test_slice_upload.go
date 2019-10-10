package main

import (
	"encoding/json"	// try https://github.com/json-iterator/go
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib/core"
	"github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib/utility"
	"github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/examples/test_common"
)

func main() {
	auth := test_common.EnvAuth()

	// 上传到这个 bucket
	bucket := "umu618-docs"
	key := "go-U6M.txt"

	bm := core.NewBucketManager(auth, nil, nil)

	fmt.Printf("Calling Delete(%s, %s)...\n", bucket, key)
	response, err := bm.Delete(bucket, key)
	if nil != err {
		fmt.Println("\tDelete() error:", err)
	} else {
		body, _ := ioutil.ReadAll(response.Body)
		if http.StatusOK == response.StatusCode {
			fmt.Println("\tDelete() OK!")
			fmt.Println("\tResponse:", string(body))
		} else {
			fmt.Println("\tDelete() failed, StatusCode =", response.StatusCode)
			fmt.Println("\tResponse:", string(body))
		}
	}

	// 在内存构造一个文件内容：2M 个 U，2M 个 M， 2M 个 U
	const dataSize int64 = 6 * 1024 * 1024
	data := make([]byte, dataSize)
	var i int64 = 0
	for ; i < dataSize/3; i++ {
		data[i] = 85
	}
	for ; i < dataSize/3*2; i++ {
		data[i] = 77
	}
	for ; i < dataSize; i++ {
		data[i] = 85
	}

	// 最后合成文件时的 hash
	etag := utility.ComputeEtag([]byte(data))
	fmt.Println("ETag =", etag)

	// 一个小时的超时，转为 UnixTime 毫秒数
	deadline := time.Now().Add(time.Second*3600).Unix() * 1000
	put_policy := fmt.Sprintf(`{"scope": "%s:%s","deadline": "%d"}`, bucket, key, deadline)
	fmt.Println("PutPolicy =", put_policy)
	upload_token := auth.CreateUploadToken(put_policy)
	fmt.Println("UploadToken =", upload_token)

	// 第一个分片不宜太大，因为可能遇到错误，上传太大是白费流量和时间！
	const blockSize = core.BlockSize //4 * 1024 * 1024
	const firstChunkSize = 1024

	su := core.NewSliceUpload(auth, nil, nil)

	fmt.Printf("Calling MakeBlock(%d, %d)...\n", blockSize, 0)
	response, err = su.MakeBlock(blockSize, 0, data[0:firstChunkSize], upload_token, key)
	if nil != err {
		fmt.Println("\tMakeBlock() error:", err)
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	if http.StatusOK != response.StatusCode {
		fmt.Println("\tMakeBlock() failed, StatusCode =", response.StatusCode)
		fmt.Println("\tResponse:", string(body))
		return
	}

	fmt.Println("\tMakeBlock() OK!")
	fmt.Println("\tResponse:", string(body))

	var result core.SliceUploadResult
	err = json.Unmarshal(body, &result)
	if nil != err {
		fmt.Println("\tUnmarshal() error:", err)
		return
	}

	blockCount := core.BlockCount(dataSize)
	contexts := make([]string, blockCount)
	contexts[0] = result.Context

	// 上传第 1 个 block 剩下的数据
	fmt.Printf("Calling Bput(%s, %d)...\n", contexts[0], firstChunkSize)
	response, err = su.Bput(contexts[0], firstChunkSize, data[firstChunkSize:blockSize], upload_token, key)
	if nil != err {
		fmt.Println("\tBput() error:", err)
		return
	}
	body, _ = ioutil.ReadAll(response.Body)
	if http.StatusOK != response.StatusCode {
		fmt.Println("\tBput() failed, StatusCode =", response.StatusCode)
		fmt.Println("\tResponse:", string(body))
		return
	}

	fmt.Println("\tBput() OK!")
	fmt.Println("\tResponse:", string(body))

	err = json.Unmarshal(body, &result)
	if nil != err {
		fmt.Println("\tUnmarshal() error:", err)
		return
	}

	contexts[0] = result.Context

	// 上传后续 block，每次都是一整块上传
	for block_index := int64(1); block_index < blockCount; block_index++ {
		pos := blockSize * block_index
		left_size := dataSize - pos
		var chunk_size int64
		if left_size > blockSize {
			chunk_size = blockSize
		} else {
			chunk_size = left_size
		}
		fmt.Printf("Calling MakeBlock(%d, %d)...\n", chunk_size, block_index)
		response, err = su.MakeBlock(chunk_size, block_index, data[pos:pos+chunk_size], upload_token, key)
		if nil != err {
			fmt.Println("\tMakeBlock() error:", err)
			return
		}

		body, _ := ioutil.ReadAll(response.Body)
		if http.StatusOK != response.StatusCode {
			fmt.Println("\tMakeBlock() failed, StatusCode =", response.StatusCode)
			fmt.Println("\tResponse:", string(body))
			return
		}

		fmt.Println("\tMakeBlock() OK!")
		fmt.Println("\tResponse:", string(body))

		err = json.Unmarshal(body, &result)
		if nil != err {
			fmt.Println("\tUnmarshal() error:", err)
			return
		}

		contexts[block_index] = result.Context
	}

	// 合成文件，注意与前面打印的 ETag 对比
	fmt.Printf("Calling MakeFile(%d, %s, %s)...\n", dataSize, key, contexts)
	response, err = su.MakeFile(dataSize, key, contexts, upload_token, nil)
	if nil != err {
		fmt.Println("\tMakeFile() error:", err)
		return
	}

	body, _ = ioutil.ReadAll(response.Body)
	if http.StatusOK != response.StatusCode {
		fmt.Println("\tMakeFile() failed, StatusCode =", response.StatusCode)
		fmt.Println("\tResponse:", string(body))
		return
	}

	fmt.Println("\tMakeFile() OK!")
	fmt.Println("\tResponse:", string(body))
}
