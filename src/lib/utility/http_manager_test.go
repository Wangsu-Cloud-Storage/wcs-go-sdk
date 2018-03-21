package utility

import (
	"net/http"
	"testing"
)

func TestDo(t *testing.T) {
	request, err := http.NewRequest("GET", "https://github.com/dokan-dev/dokany/", nil)
	if nil != err {
		return
	}

	t.Log("UserAgent =", request.Header.Get("User-Agent"))

	const _TOKEN = "UMU-Token"
	http_manager := NewDefaultHttpManager()
	response, err := http_manager.DoWithToken(request, _TOKEN)
	if nil != err {
		return
	}
	if _TOKEN != request.Header.Get("Authorization") {
		t.Fatalf("Authorization should be %s!", _TOKEN)
	}
	if userAgent != request.Header.Get("User-Agent") {
		t.Fatalf("UserAgent is %s, and should be %s!", request.Header.Get("User-Agent"), userAgent)
	}
	t.Log("StatusCode =", response.StatusCode)
	t.Log("UserAgent =", request.Header.Get("User-Agent"))

	auth := NewAuth("UMU", "Test")
	request.Header.Set("User-Agent", "Modified")
	response, err = http_manager.DoWithAuth(request, auth)
	if nil != err {
		return
	}
	if "UMU:MjUzZDA3OTBkMzJiMzMxODdiOGJiZDA0MDEyZDAwZjZlODc2MGYzMg==" != request.Header.Get("Authorization") {
		t.Fatalf("Authorization should be `UMU:MjUzZDA3OTBkMzJiMzMxODdiOGJiZDA0MDEyZDAwZjZlODc2MGYzMg=='!")
	}
	t.Log("StatusCode =", response.StatusCode)
	t.Log("UserAgent =", request.Header.Get("User-Agent"))
}
