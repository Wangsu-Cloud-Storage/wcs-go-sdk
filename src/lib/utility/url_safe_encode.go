package utility

import (
	"bytes"
	"encoding/base64"
	"net/url"
	"sort"
)

func UrlSafeEncode(data []byte) string {
	return base64.URLEncoding.EncodeToString(data)
}

func UrlSafeEncodeString(str string) string {
	return base64.URLEncoding.EncodeToString([]byte(str))
}

func StdEncodeString(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func UrlSafeEncodePair(bucket string, key string) string {
	return base64.URLEncoding.EncodeToString([]byte(bucket + ":" + key))
}

func MakeQuery(v url.Values) string {
	if v == nil {
		return ""
	}
	var buf bytes.Buffer
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		prefix := k + "="
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(prefix)
			buf.WriteString(v)
		}
	}
	return buf.String()
}
