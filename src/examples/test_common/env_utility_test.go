package test_common

import (
	"testing"
)

func TestGetEnvValues(t *testing.T) {
	value := getEnvValues("WcsLibAkSk")
	if 0 == len(value) {
		t.Fatal("Please set EnvironmentVariable \"WcsLibAkSk\" to \"AccessKey,SecretKey\".")
	}
	t.Log("Read WcsLibAkSk passed, len =", len(value))

	value = getEnvValues("WcsLibConfig")
	if 0 == len(value) {
		t.Log("EnvironmentVariable \"WcsLibConfig\" not found, or invalid, use default Config().")
	} else {
		t.Log("Read WcsLibConfig passed, len =", len(value))
	}
}

func TestEnvAuth(t *testing.T) {
	auth := EnvAuth()
	t.Log("AccessKey =", auth.AccessKey)
}

func TestEnvConfig(t *testing.T) {
	config := EnvConfig()
	t.Log("UseHttps =", config.UseHttps)
}
