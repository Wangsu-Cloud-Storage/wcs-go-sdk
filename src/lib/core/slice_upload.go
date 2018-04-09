package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"../utility"
)

// 分片上传
// https://wcs.chinanetcenter.com/document/API/FileUpload/SliceUpload

type SliceUpload struct {
	auth        *utility.Auth
	config      *Config
	httpManager *utility.HttpManager
	uploadBatch string
}

const (
	blockBits = 22
	BlockSize = 1 << blockBits
	blockMask = (1 << blockBits) - 1
)

// BlockCount 用来计算文件的分块数量
func BlockCount(fsize int64) int64 {
	return (fsize + blockMask) >> blockBits
}

func NewSliceUpload(auth *utility.Auth, config *Config, client *http.Client) (su *SliceUpload) {
	if nil == auth {
		panic("auth is nil")
	}
	if nil == config {
		config = NewDefaultConfig()
	}
	return &SliceUpload{auth, config, utility.NewHttpManager(client), utility.GetUuid()}
}

func (this *SliceUpload) MakeBlock(block_size int64, block_order int64, chunk []byte, upload_token string, key string) (response *http.Response, err error) {
	if block_size < 0 || BlockSize < block_size {
		err = errors.New("block_size is invalid")
		return
	}
	if 0 == len(upload_token) {
		err = errors.New("upload_token is empty")
		return
	}

	request, err := http.NewRequest("POST", this.config.GetUploadUrlPrefix()+"/mkblk/"+strconv.FormatInt(block_size, 10)+"/"+strconv.FormatInt(block_order, 10),
		bytes.NewReader(chunk))
	if nil != err {
		return
	}

	utility.AddMime(request, "application/octet-stream")
	request.Header.Set("UploadBatch", this.uploadBatch)
	if len(key) > 0 {
		request.Header.Set("Key", utility.UrlSafeEncodeString(key))
	}
	return this.httpManager.DoWithToken(request, upload_token)
}

func (this *SliceUpload) Bput(context string, offset int64, chunk []byte, upload_token string, key string) (response *http.Response, err error) {
	if 0 == len(context) {
		err = errors.New("context is empty")
		return
	}
	if 0 == len(upload_token) {
		err = errors.New("upload_token is empty")
		return
	}

	request, err := http.NewRequest("POST", this.config.GetUploadUrlPrefix()+"/bput/"+context+"/"+strconv.FormatInt(offset, 10),
		bytes.NewReader(chunk))
	if nil != err {
		return
	}

	utility.AddMime(request, "application/octet-stream")
	request.Header.Set("UploadBatch", this.uploadBatch)
	if len(key) > 0 {
		request.Header.Set("Key", utility.UrlSafeEncodeString(key))
	}
	return this.httpManager.DoWithToken(request, upload_token)
}

func (this *SliceUpload) MakeFile(size int64, key string, contexts []string, upload_token string, put_extra *PutExtra) (response *http.Response, err error) {
	if size < 0 {
		err = errors.New("size is invalid")
		return
	}
	if nil == contexts {
		err = errors.New("contexts is empty")
		return
	}
	if 0 == len(upload_token) {
		err = errors.New("upload_token is empty")
		return
	}

	url := this.config.GetUploadUrlPrefix() + "/mkfile/" + strconv.FormatInt(size, 10)
	if nil != put_extra && nil != put_extra.Params {
		for k, v := range put_extra.Params {
			if strings.HasPrefix(k, "x:") && len(v) > 0 {
				url += "/" + k + "/" + utility.UrlSafeEncodeString(v)
			}
		}
	}

	ctx := ""
	for _, c := range contexts {
		ctx += "," + c
	}
	request, err := http.NewRequest("POST", url, strings.NewReader(ctx[1:]))
	if nil != err {
		return
	}

	utility.AddMime(request, "text/plain;charset=UTF-8")
	request.Header.Set("UploadBatch", this.uploadBatch)
	if len(key) > 0 {
		request.Header.Set("Key", utility.UrlSafeEncodeString(key))
	}
	if nil != put_extra {
		if len(put_extra.MimeType) > 0 {
			request.Header.Set("MimeType", put_extra.MimeType)
		}
		if -1 != put_extra.Deadline {
			request.Header.Set("Deadline", strconv.Itoa(put_extra.Deadline))
		}
	}
	return this.httpManager.DoWithToken(request, upload_token)
}

func (this *SliceUpload) UploadFile(local_filename string, put_policy string, key string, put_extra *PutExtra) (response *http.Response, err error) {
	if 0 == len(local_filename) {
		err = errors.New("local_filename is empty")
		return
	}
	if 0 == len(put_policy) {
		err = errors.New("put_policy is empty")
		return
	}
	filename := key
	if 0 == len(filename) {
		filename = "goupload.tmp"
	}

	f, err := os.Open(local_filename)
	if err != nil {
		return
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return
	}

	var block_size int64
	// 第一个分片不宜太大，因为可能遇到错误，上传太大是白费流量和时间！
	var first_chunk_size int64

	if fi.Size() < 1024 {
		block_size = fi.Size()
		first_chunk_size = fi.Size()
	} else {
		if fi.Size() < BlockSize {
			block_size = fi.Size()
		} else {
			block_size = BlockSize
		}
		first_chunk_size = 1024
	}

	first_chunk := make([]byte, first_chunk_size)
	n, err := f.Read(first_chunk)
	if nil != err {
		return
	}
	if first_chunk_size != int64(n) {
		err = errors.New("Read size < request size")
		return
	}

	upload_token := this.auth.CreateUploadToken(put_policy)
	response, err = this.MakeBlock(block_size, 0, first_chunk, upload_token, key)
	if nil != err {
		return
	}
	if http.StatusOK != response.StatusCode {
		return
	}
	body, _ := ioutil.ReadAll(response.Body)
	var result SliceUploadResult
	err = json.Unmarshal(body, &result)
	if nil != err {
		return
	}

	block_count := BlockCount(fi.Size())
	contexts := make([]string, block_count)
	contexts[0] = result.Context

	// 上传第 1 个 block 剩下的数据
	if block_size > first_chunk_size {
		first_block_left_size := block_size - first_chunk_size
		left_chunk := make([]byte, first_block_left_size)
		n, err = f.Read(left_chunk)
		if nil != err {
			return
		}
		if first_block_left_size != int64(n) {
			err = errors.New("Read size < request size")
			return
		}
		response, err = this.Bput(contexts[0], first_chunk_size, left_chunk, upload_token, key)
		if nil != err {
			return
		}
		if http.StatusOK != response.StatusCode {
			return
		}
		body, _ := ioutil.ReadAll(response.Body)
		var result SliceUploadResult
		err = json.Unmarshal(body, &result)
		if nil != err {
			return
		}
		contexts[0] = result.Context

		// 上传后续 block，每次都是一整块上传
		for block_index := int64(1); block_index < block_count; block_index++ {
			pos := block_size * block_index
			left_size := fi.Size() - pos
			var chunk_size int64
			if left_size > block_size {
				chunk_size = block_size
			} else {
				chunk_size = left_size
			}
			block := make([]byte, chunk_size)
			n, err = f.Read(block)
			if nil != err {
				return
			}
			if chunk_size != int64(n) {
				err = errors.New("Read size < request size")
				return
			}
			response, err = this.MakeBlock(chunk_size, block_index, block, upload_token, key)
			if nil != err {
				return
			}
			if http.StatusOK != response.StatusCode {
				return
			}
			body, _ := ioutil.ReadAll(response.Body)
			err = json.Unmarshal(body, &result)
			if nil != err {
				return
			}

			contexts[block_index] = result.Context
		}
	}

	response, err = this.MakeFile(fi.Size(), key, contexts, upload_token, nil)
	return
}
