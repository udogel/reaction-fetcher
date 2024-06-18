package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Emoji struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID string `json:"id"`
}

func FetchReactions(baseURL string, config Config, client *http.Client) ([]string, error) {
	var usersID []string
	emojiResp, err := getEmoji(baseURL, config, client)
	if err != nil {
		return usersID, fmt.Errorf("failed to fetch emoji object: %w", err)
	}

	usersResp, err := getUsers(baseURL, config, client, emojiResp)
	if err != nil {
		return usersID, fmt.Errorf("failed to fetch user objects: %w", err)
	}
	usersID = extractID(usersResp)
	return usersID, nil
}

func getEmoji(baseURL string, config Config, client *http.Client) (Emoji, error) {
	var responseEmoji Emoji
	requestURL := fmt.Sprintf("%s/guilds/%s/emojis/%s", baseURL, config.GuildID, config.EmojiID)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return responseEmoji, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Authorization", "Bot "+config.Token)

	resp, err := client.Do(req)
	if err != nil {
		return responseEmoji, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return responseEmoji, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode == http.StatusOK {
		if err := json.Unmarshal(body, &responseEmoji); err != nil {
			return responseEmoji, fmt.Errorf("failed to unmarshall response: %w", err)
		}
	} else {
		return responseEmoji, fmt.Errorf("unexpected status code: %d, response body: %s", resp.StatusCode, string(body))
	}

	return responseEmoji, nil
}

func getUsers(baseURL string, config Config, client *http.Client, emoji Emoji) ([]User, error) {
	var usersResponse []User
	emojiObj := fmt.Sprintf("%s:%s", emoji.Name, emoji.ID)
	requestURL := fmt.Sprintf("%s/channels/%s/messages/%s/reactions/%s?limit=%s", baseURL, config.ChannelID, config.MessageID, emojiObj, config.ResponseLimit)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return usersResponse, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Authorization", "Bot "+config.Token)

	resp, err := client.Do(req)
	if err != nil {
		return usersResponse, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return usersResponse, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode == http.StatusOK {
		if err := json.Unmarshal(body, &usersResponse); err != nil {
			return usersResponse, fmt.Errorf("failed to unmarshall response: %w", err)
		}
	} else {
		return usersResponse, fmt.Errorf("unexpected status code: %d, response body: %s", resp.StatusCode, string(body))
	}

	return usersResponse, nil
}

func extractID(usersResponse []User) []string {
	var idArray []string
	for _, user := range usersResponse {
		idArray = append(idArray, user.ID)
	}
	return idArray
}
