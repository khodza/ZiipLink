package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// type Config struct {
// 	Env              string `yaml:"env" env:"ENV" env-default:"local"`
// 	StoragePath      string `yaml:"storage_path" env:"STORAGE" env-required:"true"`
// 	HTTPServer       `yaml:"http_server"`
// 	MongoDBUrl       string `yaml:"mongodb_url"`
// 	TelegramBotToken string `yaml:"telegram_bot_token"`
// 	DataBaseName     string `yaml:"database_name"`
// 	Domain           string `env:"DOMAIN" env-required:"true"`
// }

// type HTTPServer struct {
// 	Address     string        `yaml:"address" env:"ADDRESS" env-default:"localhost:8080"`
// 	TimeOut     time.Duration `yaml:"timeout" env:"TIMEOUT" env-default:"4s"`
// 	IdleTimeOut time.Duration `yaml:"idle_timeout" env:"IDLETIMEOUT" env-default:"4s"`
// 	User        string        `yaml:"user" env-required:"true"`
// 	Password    string        `yaml:"password" env:"HTTP_SERVER_PASSWORD" env-required:"true"`
// }

type Config struct {
	Env              string        `yaml:"env" env:"ENV"`
	StoragePath      string        `yaml:"storage_path" env:"STORAGE"`
	MongoDBUrl       string        `yaml:"mongodb_url" env:"MONGODB_URL"`
	TelegramBotToken string        `yaml:"telegram_bot_token" env:"Telegram_Bot_Token"`
	DataBaseName     string        `yaml:"database_name" env:"DB_NAME"`
	Domain           string        `yaml:"domain" env:"DOMAIN"`
	Address          string        `yaml:"address" env:"ADDRESS"`
	TimeOut          time.Duration `yaml:"timeout" env:"TIMEOUT"`
	IdleTimeOut      time.Duration `yaml:"idle_timeout" env:"IDLETIMEOUT"`
	User             string        `yaml:"user" env-required:"true" env:"USER"`
	Password         string        `yaml:"password" env:"PASSWORD"`
}

func MustLoad() *Config {
	// configPath := os.Getenv("./config/local.yaml")
	configPath := "./.env"
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	//check if file exists

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist : %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("config file is not valid : %s", err)
	}
	fmt.Println("cfg", cfg)
	return &cfg

}
