package main

import (
	"encoding/json"
	"fmt"
	"github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/examples/test_common"
	"github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib/core"
	"io/ioutil"
	"net/http"
)

func main() {
	auth := test_common.EnvAuth()
	config := core.NewDefaultConfig()

	bm := core.NewBucketManager(auth, config, nil)
	{
		// 查看音视频元数据
		response, err := bm.AvInfo("keyUrl")
		if nil != err {
			fmt.Println("AvInfo() failed:", err)
			return
		}
		body, _ := ioutil.ReadAll(response.Body)
		if http.StatusOK == response.StatusCode {
			fmt.Println(string(body))
		} else {
			fmt.Println("Failed, StatusCode =", response.StatusCode)
			var result core.CommonResult
			err = json.Unmarshal(body, &result)
			if nil != err {
				fmt.Println(string(body))
				fmt.Println("Unmarshal() failed,", err)
				return
			}
			fmt.Println("code =", result.Code)
			fmt.Println("message =", result.Message)
		}
	}

	{
		// 查看音文件数据
		response, err := bm.Stat("umu618-docs", "各种录音程序.7z")
		if nil != err {
			fmt.Println("Stat() failed:", err)
			return
		}
		body, _ := ioutil.ReadAll(response.Body)
		if http.StatusOK == response.StatusCode {
			fmt.Println(string(body))
		} else {
			fmt.Println("Failed, StatusCode =", response.StatusCode)
			fmt.Println(string(body))
		}
	}

	{
		limit := 1000
		mode := 0
		startTime := ""
		endTime := ""
		// 列举资源
		response, err := bm.List("bucketName", limit, "prefix", mode, "marker", startTime, endTime)
		if nil != err {
			fmt.Println("List() failed:", err)
			return
		}
		body, _ := ioutil.ReadAll(response.Body)
		if http.StatusOK == response.StatusCode {
			fmt.Println(string(body))
		} else {
			fmt.Println("Failed, StatusCode =", response.StatusCode)
			fmt.Println(string(body))
		}
	}

	// 其它文件管理功能参考：wcs-go-sdk/src/lib/core/bucket_manager.go
}
