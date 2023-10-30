package configs

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type DBConfig struct {
	HOST    string `toml:"HOST"`
	PORT    int    `toml:"PORT"`
	USER    string `toml:"USER"`
	PWD     string `toml:"PWD"`
	DB_NAME string `toml:"DB_NAME"`
}

type Token struct {
	NAME   string `toml:"NAME"`
	PERMIT string `toml:"PERMIT"`
}

type Allowedlist struct {
	TOKENS []Token `toml:"TOKENS"`
}

type ServerConfig struct {
	ALLOWED_IPs                 []string           `toml:"ALLOWED_IPs"`
	CLIENT_AUTH_TOKEN			[]string           `toml:"CLIENT_AUTH_TOKEN"`
	PORT                        string   	       `toml:"PORT"`
	KEEP_ALIVE_TIMEOUT          time.Duration      `toml:"KEEP_ALIVE_TIMEOUT"`
	KEEP_ALIVE_TIMEOUT_UNIT     string             `toml:"KEEP_ALIVE_TIMEOUT_UNIT"`
	USE_TLS                     bool               `toml:"USE_TLS"`
	TLS_CERT_PATH               string             `toml:"TLS_CERT_PATH"`
	TLS_KEY_PATH                string             `toml:"TLS_KEY_PATH"`
	TLS_PORT					string             `toml:"TLS_PORT"`
	SERVE_BOTH                  bool               `toml:"SERVE_BOTH"`
	TEMPORARY_PERMIT_TIME       int                `toml:"TEMPORARY_PERMIT_TIME"`
	TEMPORARY_PERMIT_TIME_UNIT  string 	           `toml:"TEMPORARY_PERMIT_TIME_UNIT"`
	HASHING_METHOD			    string             `toml:"HASHING_METHOD"`
	LOG_TIME_UNIT               string             `toml:"LOG_TIME_UNIT"`
	LOG_MAX_AGE                 int                `toml:"LOG_MAX_AGE"`
	LOG_ROTATION_TIME           int                `toml:"LOG_ROTATION_TIME"`
	LOG_FORMATTER			    string             `toml:"LOG_FORMATTER"`
}

var SERVER_CONFIG ServerConfig
var DB_CONFIG DBConfig
var ALLOWEDLIST Allowedlist

func init() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

    if _, err := toml.DecodeFile("./configs/server.toml", &SERVER_CONFIG); err != nil {
		panic(err)
	}

	if _, err := toml.DecodeFile("./configs/database.toml", &DB_CONFIG); err != nil {
		panic(err)
	}

	if _, err := toml.DecodeFile("./configs/allowlist.toml", &ALLOWEDLIST); err != nil {
		panic(err)
	}

	checkValid()
}

func checkValid() {
	if SERVER_CONFIG.KEEP_ALIVE_TIMEOUT < 0 {
		panic("KEEP_ALIVE_TIMEOUT should be bigger or equal to 0.")
	}

	switch strings.ToUpper(SERVER_CONFIG.KEEP_ALIVE_TIMEOUT_UNIT) {
		case "HOUR", "MINUTE", "SECOND", "MILLISECOND":
		default:
			panic("KEEP_ALIVE_TIMEOUT_UNIT is not valid (Require: hour, minute, second).")
	}

	if SERVER_CONFIG.TEMPORARY_PERMIT_TIME <= 0 {
		panic("TEMPORARY_PERMIT_TIME should be bigger than 0.")
	}

	switch strings.ToUpper(SERVER_CONFIG.TEMPORARY_PERMIT_TIME_UNIT) {
		case "DAY", "HOUR", "MINUTE":
		default:
			panic("TEMPORARY_PERMIT_TIME_UNIT is not valid (Require: day, hour, minute).")
	}

	if SERVER_CONFIG.LOG_MAX_AGE <= 0 {
		panic("LOG_MAX_AGE should be bigger than 0.")
	}

	if SERVER_CONFIG.LOG_ROTATION_TIME <= 0 {
		panic("LOG_ROTATION_TIME should be bigger than 0.")
	}

	switch strings.ToUpper(SERVER_CONFIG.LOG_TIME_UNIT) {
		case "DAY", "HOUR", "MINUTE", "SECOND":
		default:
			panic("TEMPORARY_PERMIT_TIME_UNIT is not valid (Require: day, hour, minute, second).")
	}
}
