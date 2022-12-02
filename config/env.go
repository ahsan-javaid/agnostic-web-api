package config

import (
	"os"
	"strings"
	utils "agnostic-web-api/utils"
)

func LoadEnv(key string) {
	envFile, err := os.ReadFile(key)
	utils.Check(err)
	for _, line := range strings.Split(string(envFile), "\n") {
		envList := strings.Split(line, "=")
    os.Setenv(envList[0], envList[1])
	}
}