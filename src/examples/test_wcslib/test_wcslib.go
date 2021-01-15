package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib/core"
	"github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib/utility"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Document: https://wcs.chinanetcenter.com/document/API/")
		fmt.Println("Available operation (NOT all implemented):")
		fmt.Println("    FileUpload/Upload AccessKey SecretKey UseHttps UploadHost ManageHost Bucket Key Deadline:Day LocalFilename")
		fmt.Println("    FileUpload/SliceUpload AccessKey SecretKey UseHttps UploadHost ManageHost Bucket Key Deadline:Day LocalFilename")
		fmt.Println("    FileUpload/AppendUpload AccessKey SecretKey UseHttps UploadHost ManageHost Bucket Key Position Deadline:Day LocalFilename")
		fmt.Println("    ResourceManage/listbucket AccessKey SecretKey UseHttps UploadHost ManageHost")
		fmt.Println("    ResourceManage/bucketstat AccessKey SecretKey UseHttps UploadHost ManageHost BucketNames StartDate EndDate")
		fmt.Println("    ResourceManage/download")
		fmt.Println("    ResourceManage/delete AccessKey SecretKey UseHttps UploadHost ManageHost Bucket Key")
		fmt.Println("    ResourceManage/stat AccessKey SecretKey UseHttps UploadHost ManageHost Bucket Key")
		fmt.Println("    ResourceManage/imageInfo URL")
		fmt.Println("    ResourceManage/exif URL")
		fmt.Println("    ResourceManage/avinfo URL")
		fmt.Println("    ResourceManage/avinfo2 URL")
		fmt.Println("    ResourceManage/list AccessKey SecretKey UseHttps UploadHost ManageHost Bucket [Limit:int] [Prefix] [Mode:int] [Marker]")
		fmt.Println("    ResourceManage/prefetch AccessKey SecretKey UseHttps UploadHost ManageHost BucketFileKeys")
		fmt.Println("    ResourceManage/move AccessKey SecretKey UseHttps UploadHost ManageHost Source Destination")
		fmt.Println("    ResourceManage/copy AccessKey SecretKey UseHttps UploadHost ManageHost Source Destination")
		fmt.Println("    ResourceManage/decompression")
		fmt.Println("    ResourceManage/setdeadline AccessKey SecretKey UseHttps UploadHost ManageHost Bucket Key Deadline [Relevance]")
		fmt.Println("    ResourceManage/PersistentStatus UseHttps UploadHost ManageHost PersistentId")
		fmt.Println("    Fmgr/fetch AccessKey SecretKey UseHttps UploadHost ManageHost FetchUrl Bucket [Key] [Prefix] [Md5] [Decompression] [NotifyURL] [Force] [Separate]")
		fmt.Println("    Fmgr/copy AccessKey SecretKey UseHttps UploadHost ManageHost Resource Bucket [Key] [Prefix] [NotifyURL] [Separate]")
		fmt.Println("    Fmgr/move AccessKey SecretKey UseHttps UploadHost ManageHost Resource Bucket [Key] [Prefix] [NotifyURL] [Separate]")
		fmt.Println("    Fmgr/delete AccessKey SecretKey UseHttps UploadHost ManageHost Bucket Key [NotifyURL] [Separate]")
		fmt.Println("    Fmgr/deletePrefix AccessKey SecretKey UseHttps UploadHost ManageHost Bucket Prefix [Output] [NotifyURL] [Separate]")
		fmt.Println("    Fmgr/deletem3u8 AccessKey SecretKey UseHttps UploadHost ManageHost Bucket Key [Deletets] [NotifyURL] [Separate]")
		fmt.Println("    Fmgr/setdeadline AccessKey SecretKey UseHttps UploadHost ManageHost Bucket [Prefix] Deadline [NotifyURL] [Separate]")
		fmt.Println("    Fmgr/status UseHttps UploadHost ManageHost PersistentId")
		fmt.Println("    Video-op/fops AccessKey SecretKey UseHttps UploadHost ManageHost QueryString")
		fmt.Println("    Image-op/imageDetect AccessKey SecretKey UseHttps UploadHost ManageHost Image Type Bucket")
		fmt.Println("    Image-op/imagePersistentOp AccessKey SecretKey UseHttps UploadHost ManageHost Bucket Key FOps [NotifyURL] [Force] [Separate]")

		fmt.Println("    UrlSafeEncode string ...")
		fmt.Println("    UrlSafeEncodePair key value")

		Exit(-1, "No parameter specified!")
		return
	}
	switch os.Args[1] {
	case "FileUpload/Upload":
		SimpleUpload()
	case "FileUpload/SliceUpload":
		SliceUpload()
	case "FileUpload/AppendUpload":
		AppendUpload()

	case "ResourceManage/listbucket":
		//Exit(-2, "Not implemented")
		ResourceManageListBucket()
	case "ResourceManage/bucketstat":
		//Exit(-2, "Not implemented")
		ResourceManageBucketStat()
	case "ResourceManage/download":
		Exit(-2, "Not implemented")
	case "ResourceManage/delete":
		ResourceManageDelete()
	case "ResourceManage/stat":
		ResourceManageStat()
	case "ResourceManage/imageInfo":
		ResourceManageImageInfo()
	case "ResourceManage/exif":
		ResourceManageExif()
	case "ResourceManage/avinfo":
		ResourceManageAvinfo()
	case "ResourceManage/avinfo2":
		ResourceManageAvinfo2()
	case "ResourceManage/list":
		//Exit(-2, "Not implemented")
		ResourceManageList()
	case "ResourceManage/prefetch":
		//Exit(-2, "Not implemented")
		ResourceManagePrefetch()
	case "ResourceManage/move":
		//Exit(-2, "Not implemented")
		ResourceManageMove()
	case "ResourceManage/copy":
		//Exit(-2, "Not implemented")
		ResourceManageCopy()
	case "ResourceManage/decompression":
		//Exit(-2, "Not implemented")
		ResourceManageDecompression()
	case "ResourceManage/setdeadline":
		//Exit(-2, "Not implemented")
		ResourceManageSetDeadline()
	case "ResourceManage/PersistentStatus":
		ResourceManagePersistentStatus()

	case "Fmgr/fetch":
		FmgrFetch()
	case "Fmgr/copy":
		FmgrCopy()
	case "Fmgr/move":
		//Exit(-2, "Not implemented")
		FmgrMove()
	case "Fmgr/delete":
		//Exit(-2, "Not implemented")
		FmgrDelete()
	case "Fmgr/deletePrefix":
		//Exit(-2, "Not implemented")
		FmgrDeletePrefix()
	case "Fmgr/deletem3u8":
		//Exit(-2, "Not implemented")
		FmgrDeleteM3u8()
	case "Fmgr/setdeadline":
		FmgrSetDeadline()
	case "Fmgr/status":
		FmgrStatus()
	case "Video-op/fops":
		FOps()
	case "UrlSafeEncode":
		UrlSafeEncode()
	case "UrlSafeEncodePair":
		UrlSafeEncodePair()
	case "Image-op/imageDetect":
		ImageOpImageDetect()
	case "Image-op/imagePersistentOp":
		ImageOpImagePersistentOp()
	default:
		Exit(-1, "Unknown operation: "+os.Args[1])
	}
	return
}

func Exit(exit_code int, message string) {
	fmt.Println(message)
	os.Exit(exit_code)
}

func GetArgvInt(index int, required bool) int {
	if index < len(os.Args) {
		if arg, err := strconv.Atoi(os.Args[index]); nil == err {
			return arg
		}
	}
	if required {
		Exit(-1, fmt.Sprintf("Argument[%d]:int is required!", index))
	}
	return 0
}

func GetArgvBool(index int, required bool) bool {
	if index < len(os.Args) {
		if "false" == os.Args[index] {
			return false
		}
		if "true" == os.Args[index] {
			return true
		}
		Exit(-1, fmt.Sprintf("Argument[%d]:bool is invalid!", index))
	}
	if required {
		Exit(-1, fmt.Sprintf("Argument[%d]:bool is required!", index))
	}
	return false
}

func GetArgv(index int, required bool) string {
	if index < len(os.Args) {
		return os.Args[index]
	}

	if required {
		Exit(-1, fmt.Sprintf("Argument[%d] is required!", index))
	}
	return ""
}

func SimpleUpload() {
	ak := GetArgv(2, true)
	sk := GetArgv(3, true)
	use_https := GetArgvBool(4, true)
	upload_host := GetArgv(5, true)
	manage_host := GetArgv(6, true)

	bucket := GetArgv(7, true)
	key := GetArgv(8, true)

	days_remain_to_delete := GetArgvInt(9, true)
	localfilename := GetArgv(10, true)

	auth := utility.NewAuth(ak, sk)
	config := core.NewConfig(use_https, upload_host, manage_host)
	su := core.NewSimpleUpload(auth, config, nil)
	put_extra := core.NewPutExtra(days_remain_to_delete)

	deadline := time.Now().Add(time.Second*3600).Unix() * 1000
	put_policy := "{\"scope\": \"" + bucket + "\",\"deadline\": \"" + strconv.FormatInt(deadline, 10) + "\"}"
	response, err := su.UploadFile(localfilename, put_policy, key, put_extra)
	if nil != err {
		Exit(-3, fmt.Sprintf("UploadFile() failed: %s", err))
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	Exit(response.StatusCode, string(body))
	return
}

func SliceUpload() {
	ak := GetArgv(2, true)
	sk := GetArgv(3, true)
	use_https := GetArgvBool(4, true)
	upload_host := GetArgv(5, true)
	manage_host := GetArgv(6, true)

	bucket := GetArgv(7, true)
	key := GetArgv(8, true)

	days_remain_to_delete := GetArgvInt(9, true)
	localfilename := GetArgv(10, true)

	pool_size_str := GetArgv(11, false)
	pool_size := int(1)
	if pool_size_str != "" {
		pool_size, _ = strconv.Atoi(pool_size_str)
	}

	block_size_str := GetArgv(12, false)
	var block_size int
	if block_size_str != "" {
		block_size, _ = strconv.Atoi(block_size_str)
	}

	auth := utility.NewAuth(ak, sk)
	config := core.NewConfig(use_https, upload_host, manage_host)
	su := core.NewSliceUpload(auth, config, nil)
	put_extra := core.NewPutExtra(days_remain_to_delete)

	deadline := time.Now().Add(time.Second*3600).Unix() * 1000
	put_policy := "{\"scope\": \"" + bucket + "\",\"deadline\": \"" + strconv.FormatInt(deadline, 10) + "\"}"

	var response *http.Response
	var err error
	if block_size_str != "" {
		response, err = su.UploadFileWithBlockSize(localfilename, put_policy, key, put_extra, block_size)
	} else {
		if pool_size == 1 {
			response, err = su.UploadFile(localfilename, put_policy, key, put_extra)
		} else {
			response, err = su.UploadFileConcurrent(localfilename, put_policy, key, put_extra, pool_size)
		}
	}

	if nil != err {
		Exit(-3, fmt.Sprintf("UploadFile() failed: %s", err))
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	Exit(response.StatusCode, string(body))
	return
}

func AppendUpload() {
	ak := GetArgv(2, true)
	sk := GetArgv(3, true)
	use_https := GetArgvBool(4, true)
	upload_host := GetArgv(5, true)
	manage_host := GetArgv(6, true)

	bucket := GetArgv(7, true)
	key := GetArgv(8, true)

	position := GetArgvInt(9, true)

	days_remain_to_delete := GetArgvInt(10, true)
	localfilename := GetArgv(11, true)

	auth := utility.NewAuth(ak, sk)
	config := core.NewConfig(use_https, upload_host, manage_host)
	au := core.NewAppendUpload(auth, config, nil)
	put_extra := core.NewPutExtra(days_remain_to_delete)

	deadline := time.Now().Add(time.Second*3600).Unix() * 1000
	put_policy := "{\"scope\": \"" + bucket + "\",\"deadline\": \"" + strconv.FormatInt(deadline, 10) + "\"}"
	response, err := au.AppendFile(localfilename, position, put_policy, key, put_extra)
	if nil != err {
		Exit(-3, fmt.Sprintf("AppendFile() failed: %s", err))
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	Exit(response.StatusCode, string(body))
	return
}

func ResourceManageDelete() {
	ak := GetArgv(2, true)
	sk := GetArgv(3, true)
	use_https := GetArgvBool(4, true)
	upload_host := GetArgv(5, true)
	manage_host := GetArgv(6, true)

	bucket := GetArgv(7, true)
	key := GetArgv(8, true)

	auth := utility.NewAuth(ak, sk)
	config := core.NewConfig(use_https, upload_host, manage_host)
	bm := core.NewBucketManager(auth, config, nil)
	response, err := bm.Delete(bucket, key)
	if nil != err {
		Exit(-3, fmt.Sprintf("Delete() failed: %s", err))
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	Exit(response.StatusCode, string(body))
	return
}

func ResourceManageStat() {
	ak := GetArgv(2, true)
	sk := GetArgv(3, true)
	use_https := GetArgvBool(4, true)
	upload_host := GetArgv(5, true)
	manage_host := GetArgv(6, true)

	bucket := GetArgv(7, true)
	key := GetArgv(8, true)

	auth := utility.NewAuth(ak, sk)
	config := core.NewConfig(use_https, upload_host, manage_host)
	bm := core.NewBucketManager(auth, config, nil)
	response, err := bm.Stat(bucket, key)
	if nil != err {
		Exit(-3, fmt.Sprintf("Stat() failed: %s", err))
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	Exit(response.StatusCode, string(body))
	return
}

func ResourceManageImageInfo() {
	url := GetArgv(2, true)

	auth := utility.NewAuth("", "")
	config := core.NewDefaultConfig()
	bm := core.NewBucketManager(auth, config, nil)
	response, err := bm.ImageInfo(url)
	if nil != err {
		Exit(-3, fmt.Sprintf("ImageInfo() failed: %s", err))
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	Exit(response.StatusCode, string(body))
	return
}

func ResourceManageExif() {
	url := GetArgv(2, true)

	auth := utility.NewAuth("", "")
	config := core.NewDefaultConfig()
	bm := core.NewBucketManager(auth, config, nil)
	response, err := bm.Exif(url)
	if nil != err {
		Exit(-3, fmt.Sprintf("Exif() failed: %s", err))
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	Exit(response.StatusCode, string(body))
	return
}

func ResourceManageAvinfo() {
	url := GetArgv(2, true)

	auth := utility.NewAuth("", "")
	config := core.NewDefaultConfig()
	bm := core.NewBucketManager(auth, config, nil)
	response, err := bm.AvInfo(url)
	if nil != err {
		Exit(-3, fmt.Sprintf("AvInfo() failed: %s", err))
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	Exit(response.StatusCode, string(body))
	return
}

func ResourceManageAvinfo2() {
	url := GetArgv(2, true)

	auth := utility.NewAuth("", "")
	config := core.NewDefaultConfig()
	bm := core.NewBucketManager(auth, config, nil)
	response, err := bm.AvInfo2(url)
	if nil != err {
		Exit(-3, fmt.Sprintf("AvInfo2() failed: %s", err))
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	Exit(response.StatusCode, string(body))
	return
}

func ResourceManagePersistentStatus() {
	use_https := GetArgvBool(2, true)
	upload_host := GetArgv(3, true)
	manage_host := GetArgv(4, true)

	persistent_id := GetArgv(5, true)

	auth := utility.NewAuth("", "")
	config := core.NewConfig(use_https, upload_host, manage_host)
	bm := core.NewBucketManager(auth, config, nil)
	response, err := bm.PersistentStatus(persistent_id)
	if nil != err {
		Exit(-3, fmt.Sprintf("PersistentStatus() failed: %s", err))
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	Exit(response.StatusCode, string(body))
}

func FmgrFetch() {
	ak := GetArgv(2, true)
	sk := GetArgv(3, true)
	use_https := GetArgvBool(4, true)
	upload_host := GetArgv(5, true)
	manage_host := GetArgv(6, true)

	fetch_url := GetArgv(7, true)
	bucket := GetArgv(8, true)
	key := GetArgv(9, false)
	prefix := GetArgv(10, false)
	md5 := GetArgv(11, false)
	decompression := GetArgv(12, false)
	notify_url := GetArgv(13, false)
	force := GetArgvInt(14, false)
	separate := GetArgvInt(15, false)

	auth := utility.NewAuth(ak, sk)
	config := core.NewConfig(use_https, upload_host, manage_host)
	fm := core.NewFileManager(auth, config, nil)
	response, err := fm.Fetch(fetch_url, bucket, key, prefix, md5, decompression, notify_url, force, separate)
	if nil != err {
		Exit(-3, fmt.Sprintf("Fetch() failed: %s", err))
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	Exit(response.StatusCode, string(body))
	return
}

func FmgrCopy() {
	ak := GetArgv(2, true)
	sk := GetArgv(3, true)
	use_https := GetArgvBool(4, true)
	upload_host := GetArgv(5, true)
	manage_host := GetArgv(6, true)

	resource := GetArgv(7, true)
	bucket := GetArgv(8, true)
	key := GetArgv(9, false)
	prefix := GetArgv(10, false)
	notify_url := GetArgv(11, false)
	separate := GetArgvInt(12, false)

	auth := utility.NewAuth(ak, sk)
	config := core.NewConfig(use_https, upload_host, manage_host)
	fm := core.NewFileManager(auth, config, nil)
	response, err := fm.Copy(resource, bucket, key, prefix, notify_url, separate)
	if nil != err {
		Exit(-3, fmt.Sprintf("Copy() failed: %s", err))
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	Exit(response.StatusCode, string(body))
	return
}

func FmgrStatus() {
	use_https := GetArgvBool(2, true)
	upload_host := GetArgv(3, true)
	manage_host := GetArgv(4, true)

	persistent_id := GetArgv(5, true)

	auth := utility.NewAuth("", "")
	config := core.NewConfig(use_https, upload_host, manage_host)
	fm := core.NewFileManager(auth, config, nil)
	response, err := fm.Status(persistent_id)
	if nil != err {
		Exit(-3, fmt.Sprintf("Status() failed: %s", err))
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	Exit(response.StatusCode, string(body))
}

func FOps() {
	ak := GetArgv(2, true)
	sk := GetArgv(3, true)
	use_https := GetArgvBool(4, true)
	upload_host := GetArgv(5, true)
	manage_host := GetArgv(6, true)

	query := GetArgv(7, true)

	auth := utility.NewAuth(ak, sk)
	config := core.NewConfig(use_https, upload_host, manage_host)
	response, err := core.FOps(auth, config, nil, query)
	if nil != err {
		Exit(-3, fmt.Sprintf("FOps() failed: %s", err))
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	Exit(response.StatusCode, string(body))
	return
}

func UrlSafeEncode() {
	for i := 2; i < len(os.Args); i++ {
		fmt.Println(os.Args[i], "=>", utility.UrlSafeEncodeString(os.Args[i]))
	}
}

func UrlSafeEncodePair() {
	if len(os.Args) < 4 {
		fmt.Println("Not enough arguments")
	} else {
		fmt.Println(os.Args[2], os.Args[3], "=>", utility.UrlSafeEncodePair(os.Args[2], os.Args[3]))
	}
}

// 二期需求
func ResourceManageListBucket() {
	ak := GetArgv(2, true)
	sk := GetArgv(3, true)
	use_https := GetArgvBool(4, true)
	upload_host := GetArgv(5, true)
	manage_host := GetArgv(6, true)

	auth := utility.NewAuth(ak, sk)
	config := core.NewConfig(use_https, upload_host, manage_host)
	bm := core.NewBucketManager(auth, config, nil)
	response, err := bm.ListBucket()
	if nil != err {
		Exit(-3, fmt.Sprintf("ListBucket() failed: %s", err))
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	Exit(response.StatusCode, string(body))
	return
}

func ResourceManageBucketStat() {
	ak := GetArgv(2, true)
	sk := GetArgv(3, true)
	use_https := GetArgvBool(4, true)
	upload_host := GetArgv(5, true)
	manage_host := GetArgv(6, true)

	bucket_names := GetArgv(7, true)
	startdate := GetArgv(8, true)
	enddate := GetArgv(9, false)

	auth := utility.NewAuth(ak, sk)
	config := core.NewConfig(use_https, upload_host, manage_host)
	bm := core.NewBucketManager(auth, config, nil)
	response, err := bm.BucketStat(bucket_names, startdate, enddate)
	if nil != err {
		Exit(-3, fmt.Sprintf("BucketStat() failed: %s", err))
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	Exit(response.StatusCode, string(body))
	return
}

func ResourceManageList() {
	ak := GetArgv(2, true)
	sk := GetArgv(3, true)
	use_https := GetArgvBool(4, true)
	upload_host := GetArgv(5, true)
	manage_host := GetArgv(6, true)

	bucket := GetArgv(7, true)
	limit := GetArgvInt(8, false)
	prefix := GetArgv(9, false)
	mode := GetArgvInt(10, false)
	marker := GetArgv(11, false)
	startTime := GetArgv(12, false)
	endTime := GetArgv(13, false)

	auth := utility.NewAuth(ak, sk)
	config := core.NewConfig(use_https, upload_host, manage_host)

	bm := core.NewBucketManager(auth, config, nil)
	response, err := bm.List(bucket, limit, prefix, mode, marker, startTime, endTime)
	if nil != err {
		Exit(-3, fmt.Sprintf("List() failed: %s", err))
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	Exit(response.StatusCode, string(body))
	return
}

func ResourceManagePrefetch() {
	ak := GetArgv(2, true)
	sk := GetArgv(3, true)
	use_https := GetArgvBool(4, true)
	upload_host := GetArgv(5, true)
	manage_host := GetArgv(6, true)

	bucket_file_keys := GetArgv(7, true)

	auth := utility.NewAuth(ak, sk)
	config := core.NewConfig(use_https, upload_host, manage_host)
	bm := core.NewBucketManager(auth, config, nil)
	response, err := bm.Prefetch(bucket_file_keys)
	if nil != err {
		Exit(-3, fmt.Sprintf("Prefetch() failed: %s", err))
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	Exit(response.StatusCode, string(body))
	return
}

func ResourceManageMove() {
	ak := GetArgv(2, true)
	sk := GetArgv(3, true)
	use_https := GetArgvBool(4, true)
	upload_host := GetArgv(5, true)
	manage_host := GetArgv(6, true)

	src := GetArgv(7, true)
	dst := GetArgv(8, true)

	auth := utility.NewAuth(ak, sk)
	config := core.NewConfig(use_https, upload_host, manage_host)
	bm := core.NewBucketManager(auth, config, nil)
	response, err := bm.Move(src, dst)
	if nil != err {
		Exit(-3, fmt.Sprintf("Move() failed: %s", err))
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	Exit(response.StatusCode, string(body))
	return
}

func ResourceManageCopy() {
	ak := GetArgv(2, true)
	sk := GetArgv(3, true)
	use_https := GetArgvBool(4, true)
	upload_host := GetArgv(5, true)
	manage_host := GetArgv(6, true)

	src := GetArgv(7, true)
	dst := GetArgv(8, true)

	auth := utility.NewAuth(ak, sk)
	config := core.NewConfig(use_https, upload_host, manage_host)
	bm := core.NewBucketManager(auth, config, nil)
	response, err := bm.Copy(src, dst)
	if nil != err {
		Exit(-3, fmt.Sprintf("Copy() failed: %s", err))
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	Exit(response.StatusCode, string(body))
	return
}

func ResourceManageDecompression() {
	ak := GetArgv(2, true)
	sk := GetArgv(3, true)
	use_https := GetArgvBool(4, true)
	upload_host := GetArgv(5, true)
	manage_host := GetArgv(6, true)

	bucket := GetArgv(7, true)
	key := GetArgv(8, true)
	format := GetArgv(9, true)
	directory := GetArgv(10, false)
	save_list := GetArgv(11, false)
	notify_url := GetArgv(12, false)
	force := GetArgvInt(13, false)
	separate := GetArgvInt(14, false)

	auth := utility.NewAuth(ak, sk)
	config := core.NewConfig(use_https, upload_host, manage_host)
	bm := core.NewBucketManager(auth, config, nil)

	response, err := bm.Decompression(bucket, key, format, directory, save_list, notify_url, force, separate)
	if nil != err {
		Exit(-3, fmt.Sprintf("Copy() failed: %s", err))
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	Exit(response.StatusCode, string(body))
	return
}

func ResourceManageSetDeadline() {
	ak := GetArgv(2, true)
	sk := GetArgv(3, true)
	use_https := GetArgvBool(4, true)
	upload_host := GetArgv(5, true)
	manage_host := GetArgv(6, true)

	bucket_names := GetArgv(7, true)
	key := GetArgv(8, true)
	deadline := GetArgvInt(9, true)
	relevance := GetArgvInt(10, false)

	auth := utility.NewAuth(ak, sk)
	config := core.NewConfig(use_https, upload_host, manage_host)
	bm := core.NewBucketManager(auth, config, nil)
	response, err := bm.SetDeadline(bucket_names, key, deadline, relevance)
	if nil != err {
		Exit(-3, fmt.Sprintf("SetDeadline() failed: %s", err))
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	Exit(response.StatusCode, string(body))
	return
}

func ImageOpImageDetect() {
	ak := GetArgv(2, true)
	sk := GetArgv(3, true)
	use_https := GetArgvBool(4, true)
	upload_host := GetArgv(5, true)
	manage_host := GetArgv(6, true)

	image := GetArgv(7, true)
	_type := GetArgv(8, true)
	bucket := GetArgv(9, true)

	auth := utility.NewAuth(ak, sk)
	config := core.NewConfig(use_https, upload_host, manage_host)
	image_op := core.NewImageOp(auth, config, nil)
	response, err := image_op.ImageDetect(image, _type, bucket)
	if nil != err {
		Exit(-3, fmt.Sprintf("ImageDetect() failed: %s", err))
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	Exit(response.StatusCode, string(body))
	return
}

func ImageOpImagePersistentOp() {
	ak := GetArgv(2, true)
	sk := GetArgv(3, true)
	use_https := GetArgvBool(4, true)
	upload_host := GetArgv(5, true)
	manage_host := GetArgv(6, true)

	bucket := GetArgv(7, true)
	key := GetArgv(8, true)
	fops := GetArgv(9, true)

	notify_url := GetArgv(10, false)
	force := GetArgvInt(11, false)
	separate := GetArgvInt(12, false)

	auth := utility.NewAuth(ak, sk)
	config := core.NewConfig(use_https, upload_host, manage_host)
	image_op := core.NewImageOp(auth, config, nil)
	response, err := image_op.ImagePersistentOp(bucket, key, fops, notify_url, force, separate)
	if nil != err {
		Exit(-3, fmt.Sprintf("ImagePersistentOp() failed: %s", err))
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	Exit(response.StatusCode, string(body))
	return
}

func FmgrMove() {
	ak := GetArgv(2, true)
	sk := GetArgv(3, true)
	use_https := GetArgvBool(4, true)
	upload_host := GetArgv(5, true)
	manage_host := GetArgv(6, true)

	resource := GetArgv(7, true)
	bucket := GetArgv(8, true)
	key := GetArgv(9, false)
	prefix := GetArgv(10, false)
	notify_url := GetArgv(11, false)
	separate := GetArgvInt(12, false)

	auth := utility.NewAuth(ak, sk)
	config := core.NewConfig(use_https, upload_host, manage_host)
	fm := core.NewFileManager(auth, config, nil)
	response, err := fm.Move(resource, bucket, key, prefix, notify_url, separate)
	if nil != err {
		Exit(-3, fmt.Sprintf("Move() failed: %s", err))
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	Exit(response.StatusCode, string(body))
	return
}

func FmgrDelete() {
	ak := GetArgv(2, true)
	sk := GetArgv(3, true)
	use_https := GetArgvBool(4, true)
	upload_host := GetArgv(5, true)
	manage_host := GetArgv(6, true)

	bucket := GetArgv(7, true)
	key := GetArgv(8, true)
	notify_url := GetArgv(9, false)
	separate := GetArgvInt(10, false)

	auth := utility.NewAuth(ak, sk)
	config := core.NewConfig(use_https, upload_host, manage_host)
	fm := core.NewFileManager(auth, config, nil)
	response, err := fm.Delete(bucket, key, notify_url, separate)
	if nil != err {
		Exit(-3, fmt.Sprintf("Delete() failed: %s", err))
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	Exit(response.StatusCode, string(body))
	return
}

func FmgrDeletePrefix() {
	ak := GetArgv(2, true)
	sk := GetArgv(3, true)
	use_https := GetArgvBool(4, true)
	upload_host := GetArgv(5, true)
	manage_host := GetArgv(6, true)

	bucket := GetArgv(7, true)
	prefix := GetArgv(8, true)
	output := GetArgv(9, false)
	notify_url := GetArgv(10, false)
	separate := GetArgvInt(11, false)

	auth := utility.NewAuth(ak, sk)
	config := core.NewConfig(use_https, upload_host, manage_host)
	fm := core.NewFileManager(auth, config, nil)
	response, err := fm.DeletePrefix(bucket, prefix, output, notify_url, separate)
	if nil != err {
		Exit(-3, fmt.Sprintf("DeletePrefix() failed: %s", err))
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	Exit(response.StatusCode, string(body))
	return
}

func FmgrDeleteM3u8() {
	ak := GetArgv(2, true)
	sk := GetArgv(3, true)
	use_https := GetArgvBool(4, true)
	upload_host := GetArgv(5, true)
	manage_host := GetArgv(6, true)

	bucket := GetArgv(7, true)
	key := GetArgv(8, true)
	deletets := GetArgvInt(9, false)
	notify_url := GetArgv(10, false)
	separate := GetArgvInt(11, false)

	auth := utility.NewAuth(ak, sk)
	config := core.NewConfig(use_https, upload_host, manage_host)
	fm := core.NewFileManager(auth, config, nil)
	response, err := fm.DeleteM3u8(bucket, key, deletets, notify_url, separate)
	if nil != err {
		Exit(-3, fmt.Sprintf("DeleteM3u8() failed: %s", err))
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	Exit(response.StatusCode, string(body))
	return
}

func FmgrSetDeadline() {
	ak := GetArgv(2, true)
	sk := GetArgv(3, true)
	use_https := GetArgvBool(4, true)
	upload_host := GetArgv(5, true)
	manage_host := GetArgv(6, true)

	bucket := GetArgv(7, true)
	prefix := GetArgv(8, false)
	deadline := GetArgvInt(9, true)
	notify_url := GetArgv(10, false)
	separate := GetArgvInt(11, false)

	auth := utility.NewAuth(ak, sk)
	config := core.NewConfig(use_https, upload_host, manage_host)
	fm := core.NewFileManager(auth, config, nil)
	response, err := fm.SetDeadline(bucket, prefix, deadline, notify_url, separate)
	if nil != err {
		Exit(-3, fmt.Sprintf("SetDeadline() failed: %s", err))
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	Exit(response.StatusCode, string(body))
	return
}
