package main

import (
	"fmt"
	"github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/examples/test_common"
	"github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib/core"
	"io/ioutil"
	"net/http"
)

func main() {
	auth := test_common.EnvAuth()
	config := core.NewDefaultConfig()

	fm := core.NewFileManager(auth, config, nil)

	{
		response, err := fm.Fetch("https://www.baidu.com/img/bd_logo1.png",
			"umu618-docs", "baidu.png", "", "", "", "", 0, 0)
		if nil != err {
			fmt.Println("Fetch() failed:", err)
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
		response, err := fm.Status("1000fe132b109c1d4db4952df6c7b1ab9ceb")
		if nil != err {
			fmt.Println("Status() failed:", err)
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
