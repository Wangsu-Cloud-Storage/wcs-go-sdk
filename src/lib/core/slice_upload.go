package core

import (
	"bytes"
	"errors"
	"net/http"
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
