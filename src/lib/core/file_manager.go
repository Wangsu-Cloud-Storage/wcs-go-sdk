package core

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"../utility"
)

type FileManager struct {
	auth        *utility.Auth
	config      *Config
	httpManager *utility.HttpManager
}

type FetchInfo struct {
	fetch_url     string
	bucket        string
	key           string
	prefix        string
	md5           string
	decompression string
}

type CopyInfo struct {
	resource string
	bucket   string
	key      string
	prefix   string
}

func NewFileManager(auth *utility.Auth, config *Config, client *http.Client) (fm *FileManager) {
	if nil == auth {
		panic("auth is nil")
	}
	if nil == config {
		config = NewDefaultConfig()
	}
	return &FileManager{auth, config, utility.NewHttpManager(client)}
}

// 抓取资源（fmgr/fetch）
// https://wcs.chinanetcenter.com/document/API/Fmgr/fetch
func (this *FileManager) Fetch(fetch_url string, bucket string, key string, prefix string, md5 string, decompression string,
	notify_url string, force int, separate int) (response *http.Response, err error) {
	if 0 == len(fetch_url) {
		err = errors.New("fetch_url is empty")
		return
	}
	if 0 == len(bucket) {
		err = errors.New("bucket is empty")
		return
	}

	fops := "fetchURL/" + utility.UrlSafeEncodeString(fetch_url) + "/bucket/" + utility.UrlSafeEncodeString(bucket)
	if len(key) > 0 {
		fops += "/key/" + utility.UrlSafeEncodeString(key)
	}
	if len(prefix) > 0 {
		fops += "/prefix/" + utility.UrlSafeEncodeString(prefix)
	}
	if len(md5) > 0 {
		fops += "/md5/" + md5
	}
	if len(decompression) > 0 {
		fops += "/decompression/" + decompression
	}

	return this.fetch(fops, notify_url, force, separate)
}

func (this *FileManager) FetchMultiple(fetch_info []FetchInfo,
	notify_url string, force int, separate int) (response *http.Response, err error) {
	var fops string
	for _, v := range fetch_info {
		fops += ";" + "fetchURL/" + utility.UrlSafeEncodeString(v.fetch_url) + "/bucket/" + utility.UrlSafeEncodeString(v.bucket)
		if len(v.key) > 0 {
			fops += "/key/" + utility.UrlSafeEncodeString(v.key)
		}
		if len(v.prefix) > 0 {
			fops += "/prefix/" + utility.UrlSafeEncodeString(v.prefix)
		}
		if len(v.md5) > 0 {
			fops += "/md5/" + v.md5
		}
		if len(v.decompression) > 0 {
			fops += "/decompression/" + v.decompression
		}
	}
	return this.fetch(fops[1:], notify_url, force, separate)
}

func (this *FileManager) fetch(fops string, notify_url string, force int, separate int) (response *http.Response, err error) {
	query := make(url.Values)
	query.Add("fops", fops)
	if len(notify_url) > 0 {
		query.Add("notifyURL", utility.UrlSafeEncodeString(notify_url))
	}
	if 0 == force || 1 == force {
		query.Add("force", strconv.Itoa(force))
	}
	if 0 == separate || 1 == separate {
		query.Add("separate", strconv.Itoa(separate))
	}
	return InnerFOps(this.auth, this.httpManager.GetClient(), this.config.GetManageUrlPrefix()+"/fmgr/fetch", utility.MakeQuery(query))
}

// 复制资源（fmgr/copy）
// https://wcs.chinanetcenter.com/document/API/Fmgr/copy
func (this *FileManager) Copy(resource string, bucket string, key string, prefix string, notify_url string, separate int) (response *http.Response, err error) {
	if 0 == len(resource) {
		err = errors.New("resource is empty")
		return
	}
	if 0 == len(bucket) {
		err = errors.New("bucket is empty")
		return
	}

	fops := "resource/" + utility.UrlSafeEncodeString(resource) + "/bucket/" + utility.UrlSafeEncodeString(bucket)
	if len(key) > 0 {
		fops += "/key/" + utility.UrlSafeEncodeString(key)
	}
	if len(prefix) > 0 {
		fops += "/prefix/" + utility.UrlSafeEncodeString(prefix)
	}

	return this.copy(fops, notify_url, separate)
}

func (this *FileManager) CopyMultiple(copy_info []CopyInfo, notify_url string, separate int) (response *http.Response, err error) {
	var fops string
	for _, v := range copy_info {
		fops += ";" + "resource/" + utility.UrlSafeEncodeString(v.resource) + "/bucket/" + utility.UrlSafeEncodeString(v.bucket)
		if len(v.key) > 0 {
			fops += "/key/" + utility.UrlSafeEncodeString(v.key)
		}
		if len(v.prefix) > 0 {
			fops += "/prefix/" + utility.UrlSafeEncodeString(v.prefix)
		}
	}
	return this.copy(fops[1:], notify_url, separate)
}

func (this *FileManager) copy(fops string, notify_url string, separate int) (response *http.Response, err error) {
	query := make(url.Values)
	query.Add("fops", fops)
	if len(notify_url) > 0 {
		query.Add("notifyURL", utility.UrlSafeEncodeString(notify_url))
	}
	if 0 == separate || 1 == separate {
		query.Add("separate", strconv.Itoa(separate))
	}
	return InnerFOps(this.auth, this.httpManager.GetClient(), this.config.GetManageUrlPrefix()+"/fmgr/copy", utility.MakeQuery(query))
}

// fmgr任务查询（fmgr/status）
// https://wcs.chinanetcenter.com/document/API/Fmgr/status
func (this *FileManager) Status(persistent_id string) (response *http.Response, err error) {
	if 0 == len(persistent_id) {
		err = errors.New("persistent_id is empty")
		return
	}
	url := this.config.GetManageUrlPrefix() + "/fmgr/status?persistentId=" + persistent_id
	request, err := utility.CreateGetRequest(url)
	if nil != err {
		return
	}
	// 不需要 Token
	response, err = this.httpManager.DoWithToken(request, "")
	return
}
