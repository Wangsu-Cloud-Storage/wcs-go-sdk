package core

import (
	"github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib/utility"
)

type Config struct {
	UseHttps   bool
	UploadHost string
	ManageHost string
}

func NewConfig(useHttp bool, upload_host string, manage_host string) (config *Config) {
	if 0 == len(upload_host) {
		upload_host = "apitestuser.up0.v1.wcsapi.com"
	}
	if 0 == len(manage_host) {
		manage_host = "apitestuser.mgr0.v1.wcsapi.com"
	}
	return &Config{useHttp, upload_host, manage_host}
}

func NewDefaultConfig() (config *Config) {
	return NewConfig(false, "", "")
}

func (config *Config) GetManageUrlPrefix() (url_prefix string) {
	if config.UseHttps {
		url_prefix = "https://" + config.ManageHost
	} else {
		url_prefix = "http://" + config.ManageHost
	}
	return
}

func (config *Config) GetUploadUrlPrefix() (url_prefix string) {
	if config.UseHttps {
		url_prefix = "https://" + config.UploadHost
	} else {
		url_prefix = "http://" + config.UploadHost
	}
	return
}

const (
	VERSION        = "1.0.0.3"
	BLOCK_SIZE int = 4 * 1024 * 1024
)

func init() {
	utility.SetUserAgent("WCS-GO-SDK-" + VERSION)
}
