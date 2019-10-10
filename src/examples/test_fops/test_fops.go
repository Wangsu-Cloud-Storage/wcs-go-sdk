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
	response, err := core.FOps(auth, config, nil, "bucket=aW1hZ2Vz&key=bGVodS5tcDQ==&fops=YXZ0aHVtYi9mbHYvcy80ODB4Mzg0fHNhdmVhcy9hVzFoWjJWek9tZHFhQzVtYkhZPQ==&force=1&separate=1")
	if nil != err {
		fmt.Println("FOps() failed:", err)
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
