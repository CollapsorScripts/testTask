package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

const (
	LOCAL = "local"
	PROD  = "prod"
)

type Admin struct {
	Username string `yaml:"username"`
	ID       int64  `yaml:"id"`
}

type Panel struct {
	JwtSecret string        `yaml:"jwtSecret"`
	AESKey    string        `yaml:"aesKey"`
	Expires   time.Duration `yaml:"expires"`
}

type WalletSettings struct {
	Host           string `yaml:"host"`
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	WalletPassword string `yaml:"walletPassword"`
	Certificate    string `yaml:"certificate"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type Developer struct {
	ID       int64  `yaml:"id"`
	Username string `yaml:"username"`
	Wallet   string `yaml:"wallet"`
}

type Alerts struct {
	Enabled   bool  `yaml:"enabled" env-upd`
	ChannelID int64 `yaml:"channelId" env-upd`
}

type Reserv struct {
	Username string `yaml:"username" env-upd`
	Token    string `yaml:"token" env-upd`
}

type BotSettings struct {
	Token          string  `yaml:"token" env-upd`
	CryptoKey      string  `yaml:"cryptoKey"`
	Alerts         *Alerts `yaml:"alerts"`
	Commission     float64 `yaml:"commission"`
	TransactionUrl string  `yaml:"transactionUrl"`
	Shortner       bool    `json:"shortner"`
	Reserv         *Reserv `yaml:"reserv" env-upd`
}

type Paths struct {
	Files   string `yaml:"files"`
	LogDir  string `yaml:"logDir"`
	LogName string `yaml:"logName"`
}

type ServerConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type Config struct {
	Env          string        `yaml:"env" env-default:"local"`
	Database     *Database     `yaml:"database"`
	JwtSecret    string        `yaml:"jwtSecret"`
	Paths        *Paths        `yaml:"paths"`
	ServerConfig *ServerConfig `yaml:"server"`
}

var GlobalConfig *Config

var GlobalUpdatesBotErr error = nil

// MustLoad - загружает конфигурацию
func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic(any("Файл конфигурации по указанному пути отсутствует"))
	}

	return MustLoadByPath(path)
}

func (c *Config) Update(configPath string) error {
	var file *os.File
	var err error

	file, err = os.Create(configPath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	encoder := yaml.NewEncoder(file)
	defer func(encoder *yaml.Encoder) {
		_ = encoder.Close()
	}(encoder)

	if err := encoder.Encode(&c); err != nil {
		fmt.Println("Error encoding config to file:", err)
		return err
	}

	return nil
}

// MustLoadByPath - загружает конфигурацию по указанному пути
func MustLoadByPath(configPath string) *Config {

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic(any(fmt.Sprintf("Файл конфигурации не найден: %s", configPath)))
	}

	cfg := new(Config)

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		panic(any(fmt.Sprintf("Ошибка чтения файла конфигурации: %v", err)))
	}

	GlobalConfig = cfg

	return cfg
}

// fetchConfigPath - парсинг пути к конфигурации из флага или переменной окружения
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		path := "local.yaml"
		if _, err := os.Stat(path); os.IsNotExist(err) {
			path = "./config/local.yaml"
		}
		res = path
	}

	return res
}
