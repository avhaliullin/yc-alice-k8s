package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	K8sHost      string
	K8sCA        string
	DockerImages map[string]string
	SAKey        string
}

func LoadFromEnv() *Config {
	return &Config{
		K8sHost:      requireString("K8S_HOST"),
		K8sCA:        requireString("K8S_CA"),
		DockerImages: requireStringMap("DOCKER_IMAGES"),
		SAKey:        os.Getenv("SA_KEY"), //non-required
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

func requireStringMap(name string) map[string]string {
	list := requireStringList(name)
	result := make(map[string]string)
	for _, value := range list {
		parts := strings.SplitN(value, ":", 2)
		if len(parts) != 2 {
			panic(fmt.Sprintf("expected list of tuples in %s, found: %s", name, value))
		}
		result[parts[0]] = parts[1]
	}
	return result
}

func requireStringList(name string) []string {
	str := requireString(name)
	return strings.Split(str, ",")
}
