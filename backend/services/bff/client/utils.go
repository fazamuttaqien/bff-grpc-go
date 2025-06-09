package client

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func loadLocalEnv() any {
	if _, runningInContainer := os.LookupEnv("CONTAINER"); !runningInContainer {
		err := godotenv.Load("../.env.local")
		if err != nil {
			log.Fatalln(err)
		}
	}
	return nil
}

func GetEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalln("Environment variable not found: ", key)
	}
	return value
}
