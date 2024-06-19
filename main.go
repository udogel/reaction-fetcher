package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	GuildID       string `json:"guild_id"`
	EmojiID       string `json:"emoji_id"`
	Token         string `json:"token"`
	ChannelID     string `json:"channel_id"`
	MessageID     string `json:"message_id"`
	ResponseLimit string `json:"response_limit"`
	VersionAPI    string `json:"api_version,omitempty"`
	TimeoutTime   int    `json:"timeout_time,omitempty"`
}

func main() {
	config, err := readConfig("config.json")
	if err != nil {
		fmt.Printf("error loading configuration file: %v\n", err)
	}

	if err := validateConfig(config); err != nil {
		log.Fatal("Invalid configuration: ", err)
	}

	timeoutTime := config.TimeoutTime
	if timeoutTime == 0 {
		timeoutTime = 10
	}
	client := &http.Client{
		Timeout: time.Duration(timeoutTime) * time.Second,
	}

	baseURL := fmt.Sprintf("https://discord.com/api/v%s", config.VersionAPI)
	if config.VersionAPI == "" {
		baseURL = "https://discord.com/api"
	}

	usersID, err := FetchReactions(baseURL, config, client)
	if err != nil {
		log.Fatal("Error fetching reactions: ", err)
	} else {
		log.Print("Reactions fetched successfully")
	}

	fmt.Println(usersID) // Just here not to throw an unused var error
}

func readConfig(path string) (Config, error) {
	var config Config
	configFile, err := os.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("failed to read config file: %w", err)
	}

	if err = json.Unmarshal(configFile, &config); err != nil {
		return config, fmt.Errorf("failed to unmarshall config file: %w", err)
	}

	return config, nil
}

func validateConfig(config Config) error {
	if config.GuildID == "" || config.EmojiID == "" || config.Token == "" || config.ResponseLimit == "" || config.ChannelID == "" || config.MessageID == "" {
		return fmt.Errorf("all required config fields must be set")
	}
	return nil
}
