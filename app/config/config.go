package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	K8sHost      string
	K8sCA        []byte
	DockerImages []string
}

func LoadFromEnv() *Config {
	return &Config{
		K8sHost:      requireString("K8S_HOST"),
		K8sCA:        requireBytes("K8S_CA"),
		DockerImages: requireStringList("DOCKER_IMAGES"),
	}
}

func requireBytes(name string) []byte {
	str := requireString(name)
	res, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		panic(fmt.Errorf("failed to decode var %v with base64: %w", err))
	}
	return res
}

func requireString(name string) string {
	res, ok := os.LookupEnv(name)
	if !ok {
		panic(fmt.Sprintf("required env var %s not found", name))
	}
	return res
}

func requireStringList(name string) []string {
	str := requireString(name)
	return strings.Split(str, ",")
}
