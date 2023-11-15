package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/natemarks/looch/db"
	"github.com/natemarks/secret-hoard/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	// MasterSecretEnvVar is the name of the environment variable that contains master secret value
	MasterSecretEnvVar = "MASTER_SECRET"
)

func getEnvVar(key string) (value string, err error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return value, fmt.Errorf("Environment variable '%s' not set", key)
	}
	return value, nil
}

func getInterval() int {
	intervalStr, exists := os.LookupEnv("INTERVAL")
	if !exists {
		panic("Environment variable 'INTERVAL' not set")
	}

	interval, err := strconv.Atoi(intervalStr)
	if err != nil {
		panic(fmt.Sprintf("Failed to convert 'interval' to an integer: %v", err))
	}

	return interval
}

func getSecretFromEnvVar(envVarKey string, log *zerolog.Logger) (secret types.RDSSecretData, err error) {
	secretStr, err := getEnvVar(envVarKey)
	if err != nil {
		log.Error().Err(err).Msgf("error getting env var %s", envVarKey)
		return secret, err
	}

	err = json.Unmarshal([]byte(secretStr), &secret)
	if err != nil {
		log.Error().Err(err).Msgf("error unmarshalling secret from env var %s", envVarKey)
		return secret, err
	}
	log.Info().Msgf("secret from env var %s: %v", envVarKey, secret)
	return secret, err
}
func main() {

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	interval := getInterval()
	logger.Info().Msgf("starting with interval: %v", interval)
	for {
		// Get master secret
		masterSecret, _ := getSecretFromEnvVar(MasterSecretEnvVar, &logger)
		log.Info().Msgf("master secret host: %v", masterSecret.Host)
		err := db.TCPConn(masterSecret)
		if err != nil {
			logger.Error().Err(err).Msgf("error opening TCP connection to master secret host: %v", masterSecret.Host)
		} else {
			logger.Info().Msgf("Opened TCP connection to master secret host: %v", masterSecret.Host)
		}
		db.ConnectAndPingDB(masterSecret, &logger)
		logger.Info().Msgf("sleeping for interval: %v", interval)
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
