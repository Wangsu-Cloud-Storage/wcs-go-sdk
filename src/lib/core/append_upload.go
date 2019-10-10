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

// 追加上传
// https://www.wangsucloud.com/#/help/details/31/4727
// 每次 append object 的长度不能超过2G

type AppendUpload struct {
	auth        *utility.Auth
	config      *Config
	httpManager *utility.HttpManager
}

const (
	MAX_APPEND_SIZE = 2 * 1024 * 1024 * 1024
)

func NewAppendUpload(auth *utility.Auth, config *Config, client *http.Client) (su *AppendUpload) {
	if nil == auth {
		panic("auth is nil")
	}
	if nil == config {
		config = NewDefaultConfig()
	}
	return &AppendUpload{auth, config, utility.NewHttpManager(client)}
}

func (this *AppendUpload) AppendData(data []byte, position int, put_policy string, key string, put_extra *PutExtra) (response *http.Response, err error) {
	if len(data) > MAX_APPEND_SIZE {
		err = errors.New("exceeded maximum append size")
		return
	}
	if 0 == len(put_policy) {
		err = errors.New("put_policy is empty")
		return
	}
	filename := key
	if 0 == len(filename) {
		filename = "goappend.tmp"
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

	request, err := http.NewRequest("POST", this.config.GetUploadUrlPrefix()+"/append/"+strconv.Itoa(position), strings.NewReader(buffer.String()))
	if nil != err {
		return
	}

	utility.AddMime(request, "multipart/form-data; boundary="+writer.Boundary())
	return this.httpManager.Do(request)
}

func (this *AppendUpload) AppendFile(local_filename string, position int, put_policy string, key string, put_extra *PutExtra) (response *http.Response, err error) {
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
		filename = "goappend.tmp"
	}

	f, err := os.Open(local_filename)
	if err != nil {
		return
	}
	defer f.Close()
	fi, err := f.Stat()
	if nil != err {
		return
	}
	if fi.Size() > MAX_APPEND_SIZE {
		err = errors.New("exceeded maximum file size")
		return
	}

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

	request, err := http.NewRequest("POST", this.config.GetUploadUrlPrefix()+"/append/"+strconv.Itoa(position), strings.NewReader(buffer.String()))
	if nil != err {
		return
	}

	utility.AddMime(request, "multipart/form-data; boundary="+writer.Boundary())
	return this.httpManager.Do(request)
}
