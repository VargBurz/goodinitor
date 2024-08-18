package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"github.com/joho/godotenv"
)

var botToken string
var chatID string

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
	fmt.Println("Sending message to Telegram chat:")
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
