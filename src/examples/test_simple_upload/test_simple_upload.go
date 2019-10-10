package main

import (
	"fmt"
	"github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/examples/test_common"
	"github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib/core"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func main() {
	auth := test_common.EnvAuth()
	config := core.NewDefaultConfig()

	su := core.NewSimpleUpload(auth, config, nil)

	// UnixTime 毫秒数
	deadline := time.Now().Add(time.Second*3600).Unix() * 1000
	put_policy := "{\"scope\": \"umu618-docs\",\"deadline\": \"" + strconv.FormatInt(deadline, 10) + "\"}"
	fmt.Println(put_policy)
	{
		response, err := su.UploadData([]byte("UMUTech@qq.com"), put_policy, "UMUTech-email.txt", nil)
		if nil != err {
			fmt.Println("UploadData() failed:", err)
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
		response, err := su.UploadFile(`C:\Windows\WindowsShell.Manifest`, put_policy, "WindowsShell.txt", nil)
		if nil != err {
			fmt.Println("UploadFile() failed:", err)
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
