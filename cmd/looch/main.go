package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

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
func main() {

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	interval := getInterval()
	logger.Info().Msgf("starting with interval: %v", interval)
	for {
		logger.Info().Msgf("sleeping for interval: %v", interval)
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
