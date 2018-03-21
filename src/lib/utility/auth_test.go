package utility

import (
	"net/http"
	"strings"
	"testing"
)

// 测试中事先准备好的计算结果来自 WCS C# SDK

func TestAuthSign(t *testing.T) {
	auth := NewAuth("UMU", "key")
	sign := auth.Sign([]byte("The quick brown fox jumps over the lazy dog"))
	if "UMU:ZGU3YzliODViOGI3OGFhNmJjOGE3YTM2ZjcwYTkwNzAxYzlkYjRkOQ==" != sign {
		t.Fatal(sign, "Should be: `UMU:ZGU3YzliODViOGI3OGFhNmJjOGE3YTM2ZjcwYTkwNzAxYzlkYjRkOQ=='")
	}
	t.Log("Auth.Sign() =", sign)
}

func TestAuthSignRequest(t *testing.T) {
	auth := NewAuth("UMU", "Blog")
	req, err := http.NewRequest("GET", "http://blog.umu618.com/2017/08/09/reconstructionism-poet-keep-secret/", nil)
	if nil != err {
		t.Log(err)
	}
	sign, err := auth.SignRequest(req)
	if nil != err {
		t.Log(err)
	}
	if "UMU:MDE2ZWM0MzY0M2RjN2YwZDA5YjFkMmY3OGVkNjdhYjIyMWZkODdjMQ==" != sign {
		t.Fatal(sign, "Should be: `UMU:MDE2ZWM0MzY0M2RjN2YwZDA5YjFkMmY3OGVkNjdhYjIyMWZkODdjMQ=='")
	}
	t.Log("Auth.SignRequest() =", sign)

	req, err = http.NewRequest("GET", "http://blog.umu618.com/2018/02/15/be-pretentious-unwittingly/", strings.NewReader("UMUTech"))
	if nil != err {
		t.Log(err)
	}
	sign, err = auth.SignRequest(req)
	if nil != err {
		t.Log(err)
	}
	if "UMU:ODI3ZDY4NjI2ZGZjN2VlYWZiZTJiODA3YzFmNTZlMmEyYTYzNDA3Nw==" != sign {
		t.Fatal(sign, "Should be: `UMU:ODI3ZDY4NjI2ZGZjN2VlYWZiZTJiODA3YzFmNTZlMmEyYTYzNDA3Nw=='")
	}
	t.Log("Auth.SignRequest(body) =", sign)
}
