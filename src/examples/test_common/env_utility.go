package test_common

import (
	"os"
	"strings"
	"github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib/core"
	"github.com/Wangsu-Cloud-Storage/wcs-go-sdk/src/lib/utility"
)

func EnvAuth() (auth *utility.Auth) {
	aksk := getEnvValues("WcsLibAkSk")
	if len(aksk) < 2 {
		panic("Please set EnvironmentVariable \"WcsLibAkSk\" to \"AccessKey,SecretKey\".")
	}
	return utility.NewAuth(aksk[0], aksk[1])
}

func EnvAuthEx(key string) (auth *utility.Auth) {
	aksk := getEnvValues(key)
	if len(aksk) < 2 {
		panic(`Please set EnvironmentVariable "` + key + `" to "AccessKey,SecretKey".`)
	}
	return utility.NewAuth(aksk[0], aksk[1])
}

func EnvConfig() (config *core.Config) {
	v := getEnvValues("WcsLibConfig")
	if len(v) < 3 {
		panic("Please set EnvironmentVariable \"WcsLibConfig\" to \"UseHttps,UploadHost,ManageHost\".")
	}
	return core.NewConfig("true" == v[0], v[1], v[2])
}

func EnvConfigEx(key string) (config *core.Config) {
	v := getEnvValues(key)
	if len(v) < 3 {
		panic(`Please set EnvironmentVariable "` + key + `" to "UseHttps,UploadHost,ManageHost".`)
	}
	return core.NewConfig("true" == v[0], v[1], v[2])
}

func getEnvValues(variable string) (value []string) {
	env, exists := os.LookupEnv(variable)
	if exists {
		return strings.Split(env, ",")
	}
	return make([]string, 0)
}
