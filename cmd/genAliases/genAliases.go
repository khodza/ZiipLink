package main

import (
	"fmt"
	"os"
	"zipinit/internal/config"
	mongodb "zipinit/internal/storage/mongoDB"
)

func main() {
	cfg := config.MustLoad()
	dataBaseName := cfg.DataBaseName
	storage, err := mongodb.NewStorage(cfg.MongoDBUrl, dataBaseName)
	if err != nil {
		fmt.Println("failed to init storage", err)
		os.Exit(1)
	}
	storage.SaveAndGenerateRandomStrings(10)
}
