package application

import (
	"context"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

func buildConfig() {
	ctx := context.Background()
	path, err := os.Getwd()
	if err != nil {
		log.WithContext(ctx).
			Error(err.Error())
		panic(err)
	}
	config, err := godotenv.Read(path + "/.env")
	if err != nil {
		log.WithContext(ctx).
			Error(err.Error())
		panic(err)
	}
	configs = config
}
