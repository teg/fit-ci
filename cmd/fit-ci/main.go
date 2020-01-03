package main

import (
	"os"

	"github.com/rs/zerolog"

	"github.com/teg/fit-ci/internal/api"
)

func main() {
	config, err := ReadConfig("config.yml")
	if err != nil {
		panic(err)
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	api := api.New(config.Server, config.Github, logger)

	err = api.Start()
	if err != nil {
		panic(err)
	}
}
