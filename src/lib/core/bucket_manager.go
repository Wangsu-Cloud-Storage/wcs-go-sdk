package core

import (
	"errors"
	"github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib/utility"
	"net/http"
	"strconv"
	"strings"
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
	if 0 == len(bucket) {
		err = errors.New("bucket is empty")
		return
	}
	if 0 == len(key) {
		err = errors.New("key is empty")
		return
	}
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
	if 0 == len(bucket) {
		err = errors.New("bucket is empty")
		return
	}
	if 0 == len(key) {
		err = errors.New("key is empty")
		return
	}

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

	var queryBuilder strings.Builder
	queryBuilder.WriteString("bucket=" + utility.UrlSafeEncodeString(bucket))
	queryBuilder.WriteString("&key=" + utility.UrlSafeEncodeString(key))
	fops := "decompression/" + format
	if len(directory) > 0 {
		fops += "/dir/" + utility.UrlSafeEncodeString(directory)
	}
	if len(save_list) > 0 {
		fops += "|saveas/" + utility.UrlSafeEncodeString(save_list)
	}
	queryBuilder.WriteString("&fops=" + utility.UrlSafeEncodeString(fops))
	if len(notify_url) > 0 {
		queryBuilder.WriteString("&notifyURL=" + utility.UrlSafeEncodeString(notify_url))
	}
	if 0 == force || 1 == force {
		queryBuilder.WriteString("&force=" + strconv.Itoa(force))
	}
	if 0 == separate || 1 == separate {
		queryBuilder.WriteString("&separate=" + strconv.Itoa(separate))
	}

	query := queryBuilder.String()
	return FOps(this.auth, this.config, this.httpManager.GetClient(), query)
}

// 列举空间(list bucket)
// https://wcs.chinanetcenter.com/document/API/ResourceManage/listbucket
func (this *BucketManager) ListBucket() (response *http.Response, err error) {
	url := this.config.GetManageUrlPrefix() + "/bucket/list"
	request, err := utility.CreateGetRequest(url)
	if nil != err {
		return
	}
	response, err = this.httpManager.DoWithAuth(request, this.auth)
	return
}

// 获取空间存储量(bucket stat)
// https://wcs.chinanetcenter.com/document/API/ResourceManage/bucketstat
//   bucket_names 格式：<bucket_name1>|<bucket_name2>|……
//   startdate 是 统计开始时间，格式为yyyy-mm-dd
//   enddate 是 统计结束时间，格式为yyyy-mm-dd 注：查询的时间跨度最长为六个月
func (this *BucketManager) BucketStat(bucket_names string, startdate string, enddate string) (response *http.Response, err error) {
	if 0 == len(bucket_names) {
		err = errors.New("bucket_names is empty")
		return
	}
	if 10 != len(startdate) {
		err = errors.New("startdate is invalid")
		return
	}
	if 10 != len(enddate) {
		err = errors.New("enddate is invalid")
		return
	}

	url := this.config.GetManageUrlPrefix() + "/bucket/stat?name=" + utility.UrlSafeEncodeString(bucket_names) +
		"&startdate=" + startdate + "&enddate=" + enddate
	request, err := utility.CreateGetRequest(url)
	if nil != err {
		return
	}
	response, err = this.httpManager.DoWithAuth(request, this.auth)
	return
}

// 列举资源(list)
// https://wcs.chinanetcenter.com/document/API/ResourceManage/list
func (this *BucketManager) List(bucket string, limit int, prefix string, mode int, marker string, startTime string, endTime string) (response *http.Response, err error) {
	if 0 == len(bucket) {
		err = errors.New("bucket is empty")
		return
	}
	var query = []string{}
	query = append(query, "bucket="+bucket)
	if limit >= 1 && limit <= 1000 {
		query = append(query, "limit="+strconv.Itoa(limit))
	}
	if len(prefix) > 0 {
		query = append(query, "prefix="+utility.UrlSafeEncodeString(prefix))
	}
	if 0 == mode || 1 == mode {
		query = append(query, "mode="+strconv.Itoa(mode))
	}
	if len(marker) > 0 {
		query = append(query, "marker="+marker)
	}
	if len(startTime) > 0 {
		query = append(query, "startTime="+startTime)
	}
	if len(endTime) > 0 {
		query = append(query, "endTime="+endTime)
	}

	url := this.config.GetManageUrlPrefix() + "/list?" + strings.Join(query, "&")
	request, err := utility.CreateGetRequest(url)
	if nil != err {
		return
	}
	response, err = this.httpManager.DoWithAuth(request, this.auth)
	return
}

// 更新镜像资源
// https://wcs.chinanetcenter.com/document/API/ResourceManage/prefetch
//   bucket_file_keys 格式：空间名+"："+文件名1|文件名2|文件名3…
func (this *BucketManager) Prefetch(bucket_file_keys string) (response *http.Response, err error) {
	if 0 == len(bucket_file_keys) {
		err = errors.New("bucket_file_keys is empty")
		return
	}

	url := this.config.GetManageUrlPrefix() + "/prefetch/" + utility.UrlSafeEncodeString(bucket_file_keys)
	request, err := utility.CreatePostRequest(url)
	if nil != err {
		return
	}
	response, err = this.httpManager.DoWithAuth(request, this.auth)
	return
}

// 移动资源(move)
// https://wcs.chinanetcenter.com/document/API/ResourceManage/move
func (this *BucketManager) Move(src string, dst string) (response *http.Response, err error) {
	if 0 == len(src) {
		err = errors.New("src is empty")
		return
	}
	if 0 == len(dst) {
		err = errors.New("dst is empty")
		return
	}

	url := this.config.GetManageUrlPrefix() + "/move/" + utility.UrlSafeEncodeString(src) +
		"/" + utility.UrlSafeEncodeString(dst)
	request, err := utility.CreatePostRequest(url)
	if nil != err {
		return
	}
	response, err = this.httpManager.DoWithAuth(request, this.auth)
	return
}

// 复制资源(copy)
// https://wcs.chinanetcenter.com/document/API/ResourceManage/copy
func (this *BucketManager) Copy(src string, dst string) (response *http.Response, err error) {
	if 0 == len(src) {
		err = errors.New("src is empty")
		return
	}
	if 0 == len(dst) {
		err = errors.New("dst is empty")
		return
	}

	url := this.config.GetManageUrlPrefix() + "/copy/" + utility.UrlSafeEncodeString(src) +
		"/" + utility.UrlSafeEncodeString(dst)
	request, err := utility.CreatePostRequest(url)
	if nil != err {
		return
	}
	response, err = this.httpManager.DoWithAuth(request, this.auth)
	return
}

// 设置文件保存期限
// https://wcs.chinanetcenter.com/document/API/ResourceManage/setdeadline
func (this *BucketManager) SetDeadline(bucket string, key string, deadline int, relevance int) (response *http.Response, err error) {
	if 0 == len(bucket) {
		err = errors.New("bucket is empty")
		return
	}
	if 0 == len(key) {
		err = errors.New("key is empty")
		return
	}

	if deadline < -1 {
		err = errors.New("deadline is invalid")
		return
	}

	url := this.config.GetManageUrlPrefix() + "/setdeadline"
	data := "bucket=" + utility.UrlSafeEncodeString(bucket) + "&key=" +
		utility.UrlSafeEncodeString(key) + "&deadline=" + strconv.Itoa(deadline)
	if 0 == relevance || 1 == relevance {
		data += "&relevance=" + strconv.Itoa(relevance)
	}
	return InnerFOps(this.auth, this.httpManager.GetClient(), url, data)
}
