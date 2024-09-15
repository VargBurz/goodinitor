package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var botToken string
var chatID string

type Update struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID       int    `json:"id"`
			Username string `json:"username"`
		} `json:"from"`
		Text string `json:"text"`
		Chat struct {
			ID int `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}

type UpdateResponse struct {
	OK     bool     `json:"ok"`
	Result []Update `json:"result"`
}

// Initialize the variables in an init function
func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	botToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID = os.Getenv("TELEGRAM_CHAT_ID")

	fmt.Println("TELEGRAM_BOT_TOKEN:", botToken, "TELEGRAM_CHAT_ID:", chatID)

	if botToken == "" || chatID == "" {
		fmt.Println("Error: TELEGRAM_BOT_TOKEN and TELEGRAM_CHAT_ID environment variables are required")
		os.Exit(1)
	}
}

// Function to send a message to a Telegram chat
func sendTelegramMessage(message string) error {
	const telegramAPI = "https://api.telegram.org/bot"
	apiURL := fmt.Sprintf("%s%s/sendMessage", telegramAPI, botToken)

	// Prepare data for sending the message
	data := url.Values{}
	data.Set("chat_id", chatID)
	data.Set("text", message)
	data.Set("parse_mode", "Markdown")
	fmt.Println("[sendTelegramMessage] Sending message to Telegram chat:")
	// Send the message via POST request
	response, err := http.PostForm(apiURL, data)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	} else {
		fmt.Println("Message sent successfully!")
	}

	return nil
}

// Function to get updates from the Telegram bot
func getUpdates(offset int) ([]Update, error) {
	const telegramAPI = "https://api.telegram.org/bot"
	apiURL := fmt.Sprintf("%s%s/getUpdates", telegramAPI, botToken)

	// Prepare data for fetching updates
	data := url.Values{}
	data.Set("offset", strconv.Itoa(offset))
	// data.Set("limit", "5") // Adjust limit as needed

	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		return nil, fmt.Errorf("failed to get updates: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var updateResponse UpdateResponse
	err = json.NewDecoder(resp.Body).Decode(&updateResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if !updateResponse.OK {
		return nil, fmt.Errorf("failed to fetch updates")
	}

	return updateResponse.Result, nil
}