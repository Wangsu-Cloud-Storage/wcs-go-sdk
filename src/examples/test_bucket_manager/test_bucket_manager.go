package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"../../lib/core"
	"../test_common"
)

func main() {
	auth := test_common.EnvAuth()
	config := core.NewDefaultConfig()

	bm := core.NewBucketManager(auth, config, nil)
	{
		response, err := bm.AvInfo("http://images.w.wcsapi.biz.matocloud.com/1.mp4")
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
}
