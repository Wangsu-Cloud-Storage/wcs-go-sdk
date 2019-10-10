package core

import (
	"errors"
	"github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib/utility"
	"net/http"
	"strings"
)

// https://wcs.chinanetcenter.com/document/API/Video-op
// fops 由客户自己组装

func FOps(auth *utility.Auth, config *Config, client *http.Client, query string) (response *http.Response, err error) {
	if nil == auth {
		err = errors.New("No Auth")
		return
	}

	if nil == config {
		config = NewDefaultConfig()
	}

	return InnerFOps(auth, client, config.GetManageUrlPrefix()+"/fops", query)
}

func InnerFOps(auth *utility.Auth, client *http.Client, url string, query string) (response *http.Response, err error) {
	request, err := http.NewRequest("POST", url, strings.NewReader(query))
	if nil != err {
		return
	}

	var http_manager *utility.HttpManager
	if nil == client {
		http_manager = utility.NewDefaultHttpManager()
	} else {
		http_manager = utility.NewHttpManager(client)
	}
	response, err = http_manager.DoWithAuth(request, auth)
	return
}
