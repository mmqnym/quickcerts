package configs

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
)

type DBConfig struct {
	HOST    string   `toml:"HOST"`
	PORT    int      `toml:"PORT"`
	USER    string   `toml:"USER"`
	PWD     string   `toml:"PWD"`
	DB_NAME string   `toml:"DB_NAME"`
}

type Permission struct {
	NAME   string   `toml:"NAME"`
	TOKEN  string   `toml:"TOKEN"`
}

type Allowedlist struct {
	PERMISSIONS []Permission   `toml:"PERMISSIONS"`
}

type ServerConfig struct {
	ALLOWED_IPs                  []string        `toml:"ALLOWED_IPs"`
	USE_RUNTIME_CODE		     bool            `toml:"USE_RUNTIME_CODE"`
	RUNTIME_CODE_LENGTH		     int             `toml:"RUNTIME_CODE_LENGTH"`
	CLIENT_AUTH_TOKEN			 []string        `toml:"CLIENT_AUTH_TOKEN"`
	PORT                         string   	     `toml:"PORT"`
	KEEP_ALIVE_TIMEOUT           time.Duration   `toml:"KEEP_ALIVE_TIMEOUT"`
	KEEP_ALIVE_TIMEOUT_UNIT      string          `toml:"KEEP_ALIVE_TIMEOUT_UNIT"`
	USE_TLS                      bool            `toml:"USE_TLS"`
	TLS_CERT_PATH                string          `toml:"TLS_CERT_PATH"`
	TLS_KEY_PATH                 string          `toml:"TLS_KEY_PATH"`
	TLS_PORT					 string          `toml:"TLS_PORT"`
	TEMPORARY_PERMIT_TIME        int             `toml:"TEMPORARY_PERMIT_TIME"`
	TEMPORARY_PERMIT_TIME_UNIT   string 	     `toml:"TEMPORARY_PERMIT_TIME_UNIT"`
	HASHING_METHOD			     string          `toml:"HASHING_METHOD"`
	LOG_TIME_UNIT                string          `toml:"LOG_TIME_UNIT"`
	LOG_MAX_AGE                  int             `toml:"LOG_MAX_AGE"`
	LOG_ROTATION_TIME            int             `toml:"LOG_ROTATION_TIME"`
	LOG_FORMATTER			     string          `toml:"LOG_FORMATTER"`
}

var SERVER_CONFIG ServerConfig
var DB_CONFIG DBConfig
var ALLOWEDLIST Allowedlist

func init() {
	defer func() {
		if err := recover(); err != nil {
			color.Red(err.(error).Error())
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
	if SERVER_CONFIG.USE_RUNTIME_CODE && SERVER_CONFIG.RUNTIME_CODE_LENGTH < 6 {
		panic(errors.New("RUNTIME_CODE_LENGTH should be bigger or equal to 6"))
	}

	if SERVER_CONFIG.KEEP_ALIVE_TIMEOUT < 0 {
		panic(errors.New("KEEP_ALIVE_TIMEOUT should be bigger or equal to 0"))
	}

	switch strings.ToLower(SERVER_CONFIG.KEEP_ALIVE_TIMEOUT_UNIT) {
		case "hour", "minute", "second", "millisecond":
		default:
			panic(errors.New("KEEP_ALIVE_TIMEOUT_UNIT is not valid (Require: hour, minute, second)"))
	}

	if SERVER_CONFIG.TEMPORARY_PERMIT_TIME <= 0 {
		panic(errors.New("TEMPORARY_PERMIT_TIME should be bigger than 0"))
	}

	switch strings.ToLower(SERVER_CONFIG.TEMPORARY_PERMIT_TIME_UNIT) {
		case "day", "hour", "minute":
		default:
			panic(errors.New("TEMPORARY_PERMIT_TIME_UNIT is not valid (Require: day, hour, minute)"))
	}

	if SERVER_CONFIG.LOG_MAX_AGE <= 0 {
		panic(errors.New("LOG_MAX_AGE should be bigger than 0"))
	}

	if SERVER_CONFIG.LOG_ROTATION_TIME <= 0 {
		panic(errors.New("LOG_ROTATION_TIME should be bigger than 0"))
	}

	switch strings.ToLower(SERVER_CONFIG.LOG_TIME_UNIT) {
		case "day", "hour", "minute", "second":
		default:
			panic(errors.New("TEMPORARY_PERMIT_TIME_UNIT is not valid (Require: day, hour, minute, second)"))
	}
}
