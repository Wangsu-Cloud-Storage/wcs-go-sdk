package core

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"../utility"
)

type BucketManager struct {
	auth        *utility.Auth
	config      *Config
	httpManager *utility.HttpManager
}

func NewBucketManager(auth *utility.Auth, config *Config, client *http.Client) (bm *BucketManager) {
	if nil == auth {
		panic("auth is nil")
	}
	if nil == config {
		config = NewDefaultConfig()
	}
	return &BucketManager{auth, config, utility.NewHttpManager(client)}
}

// 获取影视频元数据（avinfo）
// https://wcs.chinanetcenter.com/document/API/ResourceManage/avinfo
func (this *BucketManager) AvInfo(av_url string) (response *http.Response, err error) {
	if 0 == len(av_url) {
		err = errors.New("av_url is empty")
		return
	}
	url := av_url + "?op=avinfo"
	request, err := utility.CreateGetRequest(url)
	if nil != err {
		return
	}
	// 不需要 Token
	response, err = this.httpManager.DoWithToken(request, "")
	return
}

// 获取音视频简单元数据（avinfo2）
// https://wcs.chinanetcenter.com/document/API/ResourceManage/avinfo2
func (this *BucketManager) AvInfo2(av_url string) (response *http.Response, err error) {
	if 0 == len(av_url) {
		err = errors.New("av_url is empty")
		return
	}
	url := av_url + "?op=avinfo2"
	request, err := utility.CreateGetRequest(url)
	if nil != err {
		return
	}
	// 不需要 Token
	response, err = this.httpManager.DoWithToken(request, "")
	return
}

// 获取图片基本信息（imageInfo）
// https://wcs.chinanetcenter.com/document/API/ResourceManage/imageInfo
func (this *BucketManager) ImageInfo(image_url string) (response *http.Response, err error) {
	if 0 == len(image_url) {
		err = errors.New("image_url is empty")
		return
	}
	url := image_url + "?op=imageInfo"
	request, err := utility.CreateGetRequest(url)
	if nil != err {
		return
	}
	// 不需要 Token
	response, err = this.httpManager.DoWithToken(request, "")
	return
}

// 获取图片 EXIF 信息
// https://wcs.chinanetcenter.com/document/API/ResourceManage/exif
func (this *BucketManager) Exif(image_url string) (response *http.Response, err error) {
	if 0 == len(image_url) {
		err = errors.New("image_url is empty")
		return
	}
	url := image_url + "?op=exif"
	request, err := utility.CreateGetRequest(url)
	if nil != err {
		return
	}
	// 不需要 Token
	response, err = this.httpManager.DoWithToken(request, "")
	return
}

// 删除文件（delete）
// https://wcs.chinanetcenter.com/document/API/ResourceManage/delete
func (this *BucketManager) Delete(bucket string, key string) (response *http.Response, err error) {
	url := this.config.GetManageUrlPrefix() + "/delete/" + utility.UrlSafeEncodePair(bucket, key)
	request, err := utility.CreatePostRequest(url)
	if nil != err {
		return
	}
	response, err = this.httpManager.DoWithAuth(request, this.auth)
	return
}

// 获取文件信息（stat）
// https://wcs.chinanetcenter.com/document/API/ResourceManage/stat
func (this *BucketManager) Stat(bucket string, key string) (response *http.Response, err error) {
	url := this.config.GetManageUrlPrefix() + "/stat/" + utility.UrlSafeEncodePair(bucket, key)
	request, err := utility.CreateGetRequest(url)
	if nil != err {
		return
	}
	response, err = this.httpManager.DoWithAuth(request, this.auth)
	return
}

// 查询持久化处理状态（status）
// https://wcs.chinanetcenter.com/document/API/ResourceManage/PersistentStatus
func (this *BucketManager) PersistentStatus(persistent_id string) (response *http.Response, err error) {
	if 0 == len(persistent_id) {
		err = errors.New("persistent_id is empty")
		return
	}
	url := this.config.GetManageUrlPrefix() + "/status/get/prefop?persistentId=" + persistent_id
	request, err := utility.CreateGetRequest(url)
	if nil != err {
		return
	}
	// 不需要 Token
	response, err = this.httpManager.DoWithToken(request, "")
	return
}

// 文件解压缩
// https://wcs.chinanetcenter.com/document/API/ResourceManage/decompression
func (this *BucketManager) Decompression(bucket string, key string, format string, directory string, save_list string, notify_url string, force int, separate int) (response *http.Response, err error) {
	if 0 == len(bucket) {
		err = errors.New("bucket is empty")
		return
	}
	if 0 == len(key) {
		err = errors.New("key is empty")
		return
	}

	query := make(url.Values)
	query.Add("bucket", utility.UrlSafeEncodeString(bucket))
	query.Add("key", utility.UrlSafeEncodeString(key))
	fops := "decompression/" + format
	if len(directory) > 0 {
		fops += "/dir/" + utility.UrlSafeEncodeString(directory)
	}
	if len(save_list) > 0 {
		fops += "|saveas/" + utility.UrlSafeEncodeString(save_list)
	}
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
	return FOps(this.auth, this.config, this.httpManager.GetClient(), query.Encode())
}
