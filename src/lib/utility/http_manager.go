package utility

import (
	"net/http"
)

// 比 C# SDK 的类少了一个 allowAutoRedirect 参数，这个可以在传入的 http.Client 上自己设置 CheckRedirect
// 比 C# SDK 的类少了一个 userAgent 参数

type HttpManager struct {
	client *http.Client
}

// 这个值会在 Config 里设定
var userAgent = "WCS-GO-SDK-0.0.0.0"

func SetUserAgent(ua string) {
	userAgent = ua
	return
}

func CreateGetRequest(url string) (request *http.Request, err error) {
	request, err = http.NewRequest("GET", url, nil)
	return
}

func CreatePostRequest(url string) (request *http.Request, err error) {
	request, err = http.NewRequest("POST", url, nil)
	return
}

func AddMime(reqest *http.Request, mime string) {
	if len(mime) > 0 {
		reqest.Header.Set("Content-Type", mime)
	}
	return
}

func NewHttpManager(client *http.Client) (http_manager *HttpManager) {
	if nil == client {
		return NewDefaultHttpManager()
	}
	return &HttpManager{client}
}

func NewDefaultHttpManager() (http_manager *HttpManager) {
	return &HttpManager{http.DefaultClient}
}

func (http_manager *HttpManager) GetClient() (client *http.Client) {
	return http_manager.client
}

func (http_manager *HttpManager) Do(reqest *http.Request) (response *http.Response, err error) {
	if _, ok := reqest.Header["User-Agent"]; !ok {
		reqest.Header.Set("User-Agent", userAgent)
	}
	return http_manager.client.Do(reqest)
}

func (this *HttpManager) DoWithAuth(reqest *http.Request, auth *Auth) (response *http.Response, err error) {
	if nil != auth {
		var token string
		token, err = auth.SignRequest(reqest)
		if nil != err {
			return
		}
		reqest.Header.Set("Authorization", token)
	}
	return this.Do(reqest)
}

func (this *HttpManager) DoWithToken(reqest *http.Request, token string) (response *http.Response, err error) {
	if len(token) > 0 {
		reqest.Header.Set("Authorization", token)
	}
	return this.Do(reqest)
}
