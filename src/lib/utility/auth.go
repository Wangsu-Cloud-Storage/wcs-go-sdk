package utility

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"io/ioutil"
	"net/http"
)

// https://wcs.chinanetcenter.com/document/API/Token/AccessToken
// https://wcs.chinanetcenter.com/document/Tools/GenerateManageToken

type Auth struct {
	AccessKey string
	SecretKey []byte
}

func NewAuth(accessKey, secretKey string) (auth *Auth) {
	return &Auth{accessKey, []byte(secretKey)}
}

/// <summary>
/// 生成上传凭证
/// https://wcs.chinanetcenter.com/document/API/Token/UploadToken
/// https://wcs.chinanetcenter.com/document/Tools/GenerateUploadToken
/// </summary>
/// <param name="putPolicy">上传策略，JSON 字符串</param>
/// <returns>上传凭证</returns>
func (this *Auth) CreateUploadToken(put_policy string) (token string) {
	return this.SignWithData([]byte(put_policy))
}

// Signature {
// public
func (this *Auth) Sign(data []byte) (token string) {
	return this.AccessKey + ":" + this.encodeSign(data)
}

func (this *Auth) SignWithData(data []byte) (token string) {
	encodedData := UrlSafeEncode(data)
	return this.AccessKey + ":" + this.encodeSign([]byte(encodedData)) + ":" + encodedData
}

// https://wcs.chinanetcenter.com/document/Tools/GenerateManageToken
func (this *Auth) SignRequest(reqest *http.Request) (token string, err error) {
	var data string
	u := reqest.URL
	if len(u.RawQuery) > 0 {
		data = u.Path + "?" + u.RawQuery + "\n"
	} else {
		data = u.Path + "\n"
	}

	var buffer []byte
	if reqest.Body != nil {
		var read_body io.ReadCloser
		read_body, err = reqest.GetBody()
		if nil == err {
			var body []byte
			body, err = ioutil.ReadAll(read_body)
			read_body.Close()
			if nil == err {
				buffer = append([]byte(data), body...)
			}
		}
	} else {
		buffer = []byte(data)
	}

	return this.Sign(buffer), nil
}

// private
func (this *Auth) encodeSign(data []byte) (sign string) {
	hm := hmac.New(sha1.New, this.SecretKey)
	hm.Write(data)
	sum := hm.Sum(nil)
	hexString := make([]byte, hex.EncodedLen(len(sum)))
	hex.Encode(hexString, sum)
	return UrlSafeEncode(hexString)
}

// Signature }
