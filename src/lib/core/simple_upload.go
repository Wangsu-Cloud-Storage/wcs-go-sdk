package core

import (
	"bytes"
	"errors"
	"github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib/utility"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// 普通上传
// https://wcs.chinanetcenter.com/document/API/FileUpload/Upload
// 若文件大小超过500M，必须使用分片上传。

type SimpleUpload struct {
	auth        *utility.Auth
	config      *Config
	httpManager *utility.HttpManager
}

func NewSimpleUpload(auth *utility.Auth, config *Config, client *http.Client) (su *SimpleUpload) {
	if nil == auth {
		panic("auth is nil")
	}
	if nil == config {
		config = NewDefaultConfig()
	}
	return &SimpleUpload{auth, config, utility.NewHttpManager(client)}
}

// <summary>
// 上传数据
// putPolicy 中的 JSON 字符串是严格模式，不允许最后一个元素后面有逗号
// https://stackoverflow.com/questions/201782/can-you-use-a-trailing-comma-in-a-json-object
// 请确保 JSON 字符串的正确性和紧凑性，最好用 JSON 库生成，而不要自己用字符串拼接。
// </summary>
// <param name="data">待上传的数据</param>
// <param name="put_policy">上传策略数据，JSON 字符串</param>
// <param name="key">可选，要保存的key</param>
// <param name="put_extra">可选，上传可选设置</param>
// <returns>上传数据后的返回结果</returns>
func (this *SimpleUpload) UploadData(data []byte, put_policy string, key string, put_extra *PutExtra) (response *http.Response, err error) {
	if 0 == len(put_policy) {
		err = errors.New("put_policy is empty")
		return
	}
	filename := key
	if 0 == len(filename) {
		filename = "goupload.tmp"
	}

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	writeMultipart(writer, this.auth.CreateUploadToken(put_policy), key, put_extra)

	var w io.Writer
	if w, err = writer.CreateFormFile("file", filename); nil != err {
		return
	}
	w.Write(data)
	writer.Close()

	request, err := http.NewRequest("POST", this.config.GetUploadUrlPrefix()+"/file/upload", strings.NewReader(buffer.String()))
	if nil != err {
		return
	}

	utility.AddMime(request, "multipart/form-data; boundary="+writer.Boundary())
	return this.httpManager.Do(request)
}

func (this *SimpleUpload) UploadFile(local_filename string, put_policy string, key string, put_extra *PutExtra) (response *http.Response, err error) {
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

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	err = writeMultipart(writer, this.auth.CreateUploadToken(put_policy), key, put_extra)
	if err != nil {
		return
	}

	var w io.Writer
	if w, err = writer.CreateFormFile("file", filename); nil != err {
		return
	}
	io.Copy(w, f)
	writer.Close()

	request, err := http.NewRequest("POST", this.config.GetUploadUrlPrefix()+"/file/upload", strings.NewReader(buffer.String()))
	if nil != err {
		return
	}

	utility.AddMime(request, "multipart/form-data; boundary="+writer.Boundary())
	return this.httpManager.Do(request)
}

func writeMultipart(writer *multipart.Writer, upload_token string, key string, put_extra *PutExtra) (err error) {
	if err = writer.WriteField("token", upload_token); nil != err {
		return
	}

	if nil != put_extra && nil != put_extra.Params {
		for k, v := range put_extra.Params {
			if strings.HasPrefix(k, "x:") && len(v) > 0 {
				err = writer.WriteField(k, v)
				if nil != err {
					return
				}
			}
		}
	}

	if len(key) > 0 {
		if err = writer.WriteField("key", key); nil != err {
			return
		}
	}

	if nil != put_extra && len(put_extra.MimeType) > 0 {
		if err = writer.WriteField("mimeType", put_extra.MimeType); nil != err {
			return
		}
	}

	if nil != put_extra && -1 != put_extra.Deadline {
		if err = writer.WriteField("deadline", strconv.Itoa(put_extra.Deadline)); nil != err {
			return
		}
	}
	return
}
