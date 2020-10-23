package core

import (
	"errors"
	"github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib/utility"
	"net/http"
	"net/url"
	"strconv"
)

type FileManager struct {
	auth        *utility.Auth
	config      *Config
	httpManager *utility.HttpManager
}

type FetchInfo struct {
	Fetch_url     string
	Bucket        string
	Key           string
	Prefix        string
	Md5           string
	Decompression string
}

type CopyInfo struct {
	Resource string
	Bucket   string
	Key      string
	Prefix   string
}

type MoveInfo struct {
	Resource string
	Bucket   string
	Key      string
	Prefix   string
}

type DeleteInfo struct {
	Bucket string
	Key    string
}

type DeletePrefixInfo struct {
	Bucket string
	Prefix string
	Output string
}

type DeleteM3u8Info struct {
	Bucket   string
	Key      string
	Deletets int
}

type SetDeadlineInfo struct {
	Bucket   string
	Prefix   string
	Deadline int
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
// Copy resource to bucket:key
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
		fops += ";fetchURL/" + utility.UrlSafeEncodeString(v.Fetch_url) + "/bucket/" + utility.UrlSafeEncodeString(v.Bucket)
		if len(v.Key) > 0 {
			fops += "/key/" + utility.UrlSafeEncodeString(v.Key)
		}
		if len(v.Prefix) > 0 {
			fops += "/prefix/" + utility.UrlSafeEncodeString(v.Prefix)
		}
		if len(v.Md5) > 0 {
			fops += "/md5/" + v.Md5
		}
		if len(v.Decompression) > 0 {
			fops += "/decompression/" + v.Decompression
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
		if 0 == len(v.Resource) {
			err = errors.New("resource is empty")
			return
		}
		if 0 == len(v.Bucket) {
			err = errors.New("bucket is empty")
			return
		}
		fops += ";resource/" + utility.UrlSafeEncodeString(v.Resource) + "/bucket/" + utility.UrlSafeEncodeString(v.Bucket)
		if len(v.Key) > 0 {
			fops += "/key/" + utility.UrlSafeEncodeString(v.Key)
		}
		if len(v.Prefix) > 0 {
			fops += "/prefix/" + utility.UrlSafeEncodeString(v.Prefix)
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

// 移动资源
// https://wcs.chinanetcenter.com/document/API/Fmgr/move
// Move resource to bucket:key
func (this *FileManager) Move(resource string, bucket string, key string, prefix string, notify_url string, separate int) (response *http.Response, err error) {
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

	return this.move(fops, notify_url, separate)
}

func (this *FileManager) MoveMultiple(move_info []MoveInfo, notify_url string, separate int) (response *http.Response, err error) {
	var fops string
	for _, v := range move_info {
		if 0 == len(v.Resource) {
			err = errors.New("resource is empty")
			return
		}
		if 0 == len(v.Bucket) {
			err = errors.New("bucket is empty")
			return
		}
		fops += ";resource/" + utility.UrlSafeEncodeString(v.Resource) + "/bucket/" + utility.UrlSafeEncodeString(v.Bucket)
		if len(v.Key) > 0 {
			fops += "/key/" + utility.UrlSafeEncodeString(v.Key)
		}
		if len(v.Prefix) > 0 {
			fops += "/prefix/" + utility.UrlSafeEncodeString(v.Prefix)
		}
	}
	return this.move(fops[1:], notify_url, separate)
}

func (this *FileManager) move(fops string, notify_url string, separate int) (response *http.Response, err error) {
	query := make(url.Values)
	query.Add("fops", fops)
	if len(notify_url) > 0 {
		query.Add("notifyURL", utility.UrlSafeEncodeString(notify_url))
	}
	if 0 == separate || 1 == separate {
		query.Add("separate", strconv.Itoa(separate))
	}
	return InnerFOps(this.auth, this.httpManager.GetClient(), this.config.GetManageUrlPrefix()+"/fmgr/move", utility.MakeQuery(query))
}

// 删除资源
// https://wcs.chinanetcenter.com/document/API/Fmgr/delete
// Delete bucket:key
func (this *FileManager) Delete(bucket string, key string, notify_url string, separate int) (response *http.Response, err error) {
	if 0 == len(bucket) {
		err = errors.New("bucket is empty")
		return
	}
	if 0 == len(key) {
		err = errors.New("key is empty")
		return
	}

	fops := "bucket/" + utility.UrlSafeEncodeString(bucket) + "/key/" + utility.UrlSafeEncodeString(key)

	return this.delete(fops, notify_url, separate)
}

func (this *FileManager) DeleteMultiple(delete_info []DeleteInfo, notify_url string, separate int) (response *http.Response, err error) {
	var fops string
	for _, v := range delete_info {
		if 0 == len(v.Bucket) {
			err = errors.New("bucket is empty")
			return
		}
		if 0 == len(v.Key) {
			err = errors.New("key is empty")
			return
		}
		fops += ";" + "bucket/" + utility.UrlSafeEncodeString(v.Bucket) + "/key/" + utility.UrlSafeEncodeString(v.Key)
	}
	return this.delete(fops[1:], notify_url, separate)
}

func (this *FileManager) delete(fops string, notify_url string, separate int) (response *http.Response, err error) {
	query := make(url.Values)
	query.Add("fops", fops)
	if len(notify_url) > 0 {
		query.Add("notifyURL", utility.UrlSafeEncodeString(notify_url))
	}
	if 0 == separate || 1 == separate {
		query.Add("separate", strconv.Itoa(separate))
	}
	return InnerFOps(this.auth, this.httpManager.GetClient(), this.config.GetManageUrlPrefix()+"/fmgr/delete", utility.MakeQuery(query))
}

// 按前缀删除资源
// https://wcs.chinanetcenter.com/document/API/Fmgr/deletePrefix
func (this *FileManager) DeletePrefix(bucket string, prefix string, output string, notify_url string, separate int) (response *http.Response, err error) {
	if 0 == len(bucket) {
		err = errors.New("bucket is empty")
		return
	}
	if 0 == len(prefix) {
		err = errors.New("prefix is empty")
		return
	}

	fops := "bucket/" + utility.UrlSafeEncodeString(bucket) +
		"/prefix/" + utility.UrlSafeEncodeString(prefix)
	if len(output) > 0 {
		fops += "/output/" + utility.UrlSafeEncodeString(output)
	}

	return this.deletePrefix(fops, notify_url, separate)
}

func (this *FileManager) DeletePrefixMultiple(delete_prefix_info []DeletePrefixInfo, notify_url string, separate int) (response *http.Response, err error) {
	var fops string
	for _, v := range delete_prefix_info {
		if 0 == len(v.Bucket) {
			err = errors.New("bucket is empty")
			return
		}
		if 0 == len(v.Prefix) {
			err = errors.New("prefix is empty")
			return
		}
		fops += ";bucket/" + utility.UrlSafeEncodeString(v.Bucket) + "/prefix/" + utility.UrlSafeEncodeString(v.Prefix)
		if len(v.Output) > 0 {
			fops += "/output/" + utility.UrlSafeEncodeString(v.Output)
		}
	}
	return this.deletePrefix(fops[1:], notify_url, separate)
}

func (this *FileManager) deletePrefix(fops string, notify_url string, separate int) (response *http.Response, err error) {
	query := make(url.Values)
	query.Add("fops", fops)
	if len(notify_url) > 0 {
		query.Add("notifyURL", utility.UrlSafeEncodeString(notify_url))
	}
	if 0 == separate || 1 == separate {
		query.Add("separate", strconv.Itoa(separate))
	}
	return InnerFOps(this.auth, this.httpManager.GetClient(), this.config.GetManageUrlPrefix()+"/fmgr/deletePrefix", utility.MakeQuery(query))
}

// 删除m3u8文件
// https://wcs.chinanetcenter.com/document/API/Fmgr/deletem3u8
func (this *FileManager) DeleteM3u8(bucket string, key string, deletets int, notify_url string, separate int) (response *http.Response, err error) {
	if 0 == len(bucket) {
		err = errors.New("bucket is empty")
		return
	}
	if 0 == len(key) {
		err = errors.New("key is empty")
		return
	}

	fops := "bucket/" + utility.UrlSafeEncodeString(bucket) +
		"/key/" + utility.UrlSafeEncodeString(key)
	if deletets == 0 || deletets == 1 {
		fops += "/deletets/" + strconv.Itoa(deletets)
	}

	return this.deleteM3u8(fops, notify_url, separate)
}

func (this *FileManager) DeleteM3u8Multiple(delete_m3u8_info []DeleteM3u8Info, notify_url string, separate int) (response *http.Response, err error) {
	var fops string
	for _, v := range delete_m3u8_info {
		if 0 == len(v.Bucket) {
			err = errors.New("bucket is empty")
			return
		}
		if 0 == len(v.Key) {
			err = errors.New("key is empty")
			return
		}
		fops += ";bucket/" + utility.UrlSafeEncodeString(v.Bucket) + "/key/" + utility.UrlSafeEncodeString(v.Key)
		if v.Deletets == 0 || v.Deletets == 1 {
			fops += "/deletets/" + strconv.Itoa(v.Deletets)
		}
	}
	return this.deleteM3u8(fops[1:], notify_url, separate)
}

func (this *FileManager) deleteM3u8(fops string, notify_url string, separate int) (response *http.Response, err error) {
	query := make(url.Values)
	query.Add("fops", fops)
	if len(notify_url) > 0 {
		query.Add("notifyURL", utility.UrlSafeEncodeString(notify_url))
	}
	if 0 == separate || 1 == separate {
		query.Add("separate", strconv.Itoa(separate))
	}
	return InnerFOps(this.auth, this.httpManager.GetClient(), this.config.GetManageUrlPrefix()+"/fmgr/deletem3u8", utility.MakeQuery(query))
}

// 批量修改文件保存期限
// https://wcs.chinanetcenter.com/document/API/Fmgr/setdeadline
func (this *FileManager) SetDeadline(bucket string, prefix string, deadline int, notify_url string, separate int) (response *http.Response, err error) {
	if 0 == len(bucket) {
		err = errors.New("bucket is empty")
		return
	}

	if deadline < -1 {
		err = errors.New("deadline is invalid")
		return
	}

	fops := "bucket/" + utility.UrlSafeEncodeString(bucket)
	if len(prefix) > 0 {
		fops += "/prefix/" + utility.UrlSafeEncodeString(prefix)
	}
	fops += "/deadline/" + strconv.Itoa(deadline)
	return this.setDeadline(fops, notify_url, separate)
}

func (this *FileManager) SetDeadlineMultiple(set_deadline_info []SetDeadlineInfo, notify_url string, separate int) (response *http.Response, err error) {
	var fops string
	for _, v := range set_deadline_info {
		if 0 == len(v.Bucket) {
			err = errors.New("bucket is empty")
			return
		}
		if v.Deadline < -1 {
			err = errors.New("deadline is invalid")
			return
		}
		fops += ";bucket/" + utility.UrlSafeEncodeString(v.Bucket)
		if len(v.Prefix) > 0 {
			fops += "/prefix/" + utility.UrlSafeEncodeString(v.Prefix)
		}
		fops += "/deadline/" + strconv.Itoa(v.Deadline)
	}
	return this.setDeadline(fops[1:], notify_url, separate)
}

func (this *FileManager) setDeadline(fops string, notify_url string, separate int) (response *http.Response, err error) {
	query := make(url.Values)
	query.Add("fops", fops)
	if len(notify_url) > 0 {
		query.Add("notifyURL", utility.UrlSafeEncodeString(notify_url))
	}
	if 0 == separate || 1 == separate {
		query.Add("separate", strconv.Itoa(separate))
	}
	return InnerFOps(this.auth, this.httpManager.GetClient(), this.config.GetManageUrlPrefix()+"/fmgr/setdeadline", utility.MakeQuery(query))
}
