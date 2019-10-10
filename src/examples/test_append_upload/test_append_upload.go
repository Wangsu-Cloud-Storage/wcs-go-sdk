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
	auth := test_common.EnvAuthEx("WcsLibAkSkAppend")
	config := test_common.EnvConfigEx("WcsLibConfigAppend")

	su := core.NewAppendUpload(auth, config, nil)

	// UnixTime 毫秒数
	deadline := time.Now().Add(time.Second*3600).Unix() * 1000
	put_policy := "{\"scope\": \"r35-wcsm2\",\"deadline\": \"" + strconv.FormatInt(deadline, 10) + "\"}"
	fmt.Println(put_policy)
	{
		response, err := su.AppendData([]byte("UMUTech@qq.com"), 0, put_policy, "UMU-append.txt", nil)
		if nil != err {
			fmt.Println("AppendData() failed:", err)
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
		response, err := su.AppendFile(`C:\Windows\WindowsShell.Manifest`, 14, put_policy, "UMU-append.txt", nil)
		if nil != err {
			fmt.Println("AppendFile() failed:", err)
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
