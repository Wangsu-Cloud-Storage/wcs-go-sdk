package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"../../lib/core"
	"../../lib/utility"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Available operation (NOT all implemented):")
		fmt.Println("    FileUpload/SimpleUpload")
		fmt.Println("    FileUpload/SliceUpload")
		fmt.Println("    FileUpload/AppendUpload")
		fmt.Println("    ResourceManage/listbucket")
		fmt.Println("    ResourceManage/bucketstat")
		fmt.Println("    ResourceManage/download")
		fmt.Println("    ResourceManage/delete AccessKey SecretKey UseHttps UploadHost ManageHost Bucket Key")
		fmt.Println("    ResourceManage/stat AccessKey SecretKey UseHttps UploadHost ManageHost Bucket Key")
		fmt.Println("    ResourceManage/imageInfo URL")
		fmt.Println("    ResourceManage/exif URL")
		fmt.Println("    ResourceManage/avinfo URL")
		fmt.Println("    ResourceManage/avinfo2 URL")
		fmt.Println("    ResourceManage/list")
		fmt.Println("    ResourceManage/prefetch")
		fmt.Println("    ResourceManage/move")
		fmt.Println("    ResourceManage/copy")
		fmt.Println("    ResourceManage/decompression")
		fmt.Println("    ResourceManage/setdeadline")
		fmt.Println("    ResourceManage/PersistentStatus UseHttps UploadHost ManageHost PersistentId")
		fmt.Println("    Fmgr/fetch fetch_url bucket [key] [prefix] [md5] [decompression] [notify_url] [force] [separate]")
		fmt.Println("    Fmgr/copy resource bucket [key] [prefix] [notify_url] [separate]")
		fmt.Println("    Fmgr/move")
		fmt.Println("    Fmgr/delete")
		fmt.Println("    Fmgr/deletePrefix")
		fmt.Println("    Fmgr/deletem3u8")
		fmt.Println("    Fmgr/status PersistentId")
		fmt.Println("    Video-op/fops AccessKey SecretKey UseHttps UploadHost ManageHost QueryString")

		Exit(-1, "No parameter specified!")
		return
	}
	switch os.Args[1] {
	case "FileUpload/SimpleUpload":
		Exit(-2, "Not implemented")
	case "FileUpload/SliceUpload":
		Exit(-2, "Not implemented")
	case "FileUpload/AppendUpload":
		Exit(-2, "Not implemented")

	case "ResourceManage/listbucket":
		Exit(-2, "Not implemented")
	case "ResourceManage/bucketstat":
		Exit(-2, "Not implemented")
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
		Exit(-2, "Not implemented")
	case "ResourceManage/prefetch":
		Exit(-2, "Not implemented")
	case "ResourceManage/move":
		Exit(-2, "Not implemented")
	case "ResourceManage/copy":
		Exit(-2, "Not implemented")
	case "ResourceManage/decompression":
		Exit(-2, "Not implemented")
	case "ResourceManage/setdeadline":
		Exit(-2, "Not implemented")
	case "ResourceManage/PersistentStatus":
		ResourceManagePersistentStatus()

	case "Fmgr/fetch":
		FmgrFetch()
	case "Fmgr/copy":
		FmgrCopy()
	case "Fmgr/move":
		Exit(-2, "Not implemented")
	case "Fmgr/delete":
		Exit(-2, "Not implemented")
	case "Fmgr/deletePrefix":
		Exit(-2, "Not implemented")
	case "Fmgr/deletem3u8":
		Exit(-2, "Not implemented")
	case "Fmgr/status":
		FmgrStatus()
	case "Video-op/fops":
		FOps()
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