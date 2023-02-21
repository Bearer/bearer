package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type TypeGoogleSheets struct {
	UserToken *oauth2.Token
	AppConfig *oauth2.Config
}

type TypeGoogleDrive struct {
	ParentFolderId string
}

type TypeRedis struct {
	Address  string
	Password string
	Db       int
	Init     bool
}

type TypeRuntime struct {
	MaxAttempt     int
	Sheets         TypeGoogleSheets
	Drive          TypeGoogleDrive
	Redis          TypeRedis
	EFSLocation    string
	ExecutablePath string
}

var Runtime = TypeRuntime{}

func Load() {
	userTokenS := os.Getenv("GOOGLE_SHEETS_USER")
	appConfigS := os.Getenv("GOOGLE_APP")
	userToken := &oauth2.Token{}

	err := json.NewDecoder(bytes.NewBuffer([]byte(userTokenS))).Decode(userToken)
	if err != nil {
		log.Fatal().Err(
			fmt.Errorf("failed decoding sheets user token `%s` to token %#v", userTokenS, err),
		).Send()
	}

	appConfig, err := google.ConfigFromJSON(
		[]byte(appConfigS),
		"https://www.googleapis.com/auth/spreadsheets",
		"https://www.googleapis.com/auth/drive",
	)
	if err != nil {
		log.Fatal().Err(fmt.Errorf("unable to parse google app `%s` to config: %e", appConfigS, err)).Send()
	}

	Runtime.ExecutablePath = getEnv("BEARER_EXECUTABLE_PATH", "/app/bearer")

	maxAttempt, err := strconv.Atoi(getEnv("GOOGLE_MAX_ATTEMPT", ""))
	if err != nil {
		Runtime.MaxAttempt = -1
	} else {
		Runtime.MaxAttempt = maxAttempt
	}

	Runtime.Sheets = TypeGoogleSheets{
		UserToken: userToken,
		AppConfig: appConfig,
	}

	Runtime.EFSLocation = os.TempDir()
	if os.Getenv("USE_EFS") == "1" {
		Runtime.EFSLocation = "/app/battle-test-tmp"
	}

	Runtime.Drive = TypeGoogleDrive{
		ParentFolderId: os.Getenv("GOOGLE_DRIVE_PARENT_FOLDER_ID"),
	}

	initValue := getEnv("REDIS_INIT", "false")
	initBool, err := strconv.ParseBool(initValue)
	if err != nil {
		log.Fatal().Err(
			fmt.Errorf("unable to parse env variable `%s` for REDIS_INIT to bool: %e", initValue, err),
		).Send()
	}

	Runtime.Redis = TypeRedis{
		Address:  os.Getenv("BATTLE_TEST_REDIS_ADDRESS"),
		Password: "",
		Db:       0,
		Init:     initBool,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
