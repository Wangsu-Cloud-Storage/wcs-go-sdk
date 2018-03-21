package test_common

import (
	"os"
	"strings"

	"../../lib/utility"
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

func getEnvValues(variable string) (value []string) {
	env, exists := os.LookupEnv(variable)
	if exists {
		return strings.Split(env, ",")
	}
	return make([]string, 0)
}
