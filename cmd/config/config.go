package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func VerifyIsDockerRun() (check bool) {
	isDocker := os.Getenv("DOCKER")

	return isDocker == "true"
}

func LoadEnv(path string) (err error) {
	if err = godotenv.Load(path); err != nil {
		return fmt.Errorf("error to laod env: %w", err)
	}

	return nil
}
