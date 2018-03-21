package core

import (
	"testing"
)

func TestUrl(t *testing.T) {
	config := NewConfig(true, "", "")
	url := config.GetManageUrlPrefix()
	if "https://apitestuser.mgr0.v1.wcsapi.com" != url {
		t.Fatal("GetManageUrlPrefix() =", url, "Should be: `https://apitestuser.mgr0.v1.wcsapi.com'")
	}

	url = config.GetUploadUrlPrefix()
	if "https://apitestuser.up0.v1.wcsapi.com" != url {
		t.Fatal("GetUploadUrlPrefix() =", url, "Should be: `https://apitestuser.up0.v1.wcsapi.com'")
	}

	conf := &Config{false, "upload.umutech.com", "umutech.com"}
	t.Log(conf.GetManageUrlPrefix())
	t.Log(conf.GetUploadUrlPrefix())
}
