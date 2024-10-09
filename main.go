package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strings"
    "time"
)


const woltParseTimeout = 30 // minutes
const telegramCheckTimeout = 15 // seconds

type WebImage struct {
    URL string `json:"url"`
}

// Define the Item struct
type Item struct {
    Name                string   `json:"name"`
    PurchasableBalance  *int     `json:"purchasable_balance"`
    Images              []WebImage  `json:"images"`
}

// Define the AssortmentResponse struct
type AssortmentResponse struct {
    Items []Item `json:"items"`
}

// Define the Result struct to store search results
type Result struct {
    Venue   string `json:"venue"`
    Name    string `json:"name"`
    Founded bool   `json:"founded"`
    Time    string `json:"time"` // using string to store the timestamp in a readable format
    Image   string `json:"image"`
}

// Define the Config struct to match the structure in config.json
type Config struct {
    Endpoint string   `json:"endpoint"`
    Names    []string `json:"names"`
    Venue    string   `json:"venue"`
}

// send test message to telegram from bot
func httpHandler(w http.ResponseWriter, r *http.Request) {
    message := "Hello from Telegram Bot!"
    err := sendTelegramMessage(message)
    if err != nil {
        fmt.Fprintf(w, "Failed to send message: %v", err)
        return
    }
	fmt.Fprintf(w, "Message sent to Telegram!")
}

func containsStatus(text string) bool {
	// Convert text to lowercase
	lowerText := strings.ToLower(text)

	// Check if "status" is present
	return strings.Contains(lowerText, "status")
}

func formatResultsMessage(results []Result) string {
    var messageBuilder strings.Builder

    for _, result := range results {
        status := "out of stock"
        if result.Founded {
            status = "available"
        }
        // Append the formatted string to the message
        messageBuilder.WriteString(fmt.Sprintf("%s | *%s* is %s\n", result.Venue, result.Name, status))
    }

    // Convert the message builder to a string
    return messageBuilder.String()
}

// Perform the GET request and parse the JSON response
func fetchItems(endpoint string) (*AssortmentResponse, error) {
    response, err := http.Get(endpoint)
    if err != nil {
        return nil, fmt.Errorf("error making GET request: %v", err)
    }
    defer response.Body.Close()

    body, err := io.ReadAll(response.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading response body: %v", err)
    }

    var assortment AssortmentResponse
    err = json.Unmarshal(body, &assortment)
    if err != nil {
        return nil, fmt.Errorf("error parsing JSON response: %v", err)
    }

    return &assortment, nil
}

// Search for items that match any of the provided names
func searchItems(assortment *AssortmentResponse, searchNames []string) map[string][]Item {
    // assortment - fetched items from api
    // searchNames - names from config
    matchingItems := make(map[string][]Item)

    for _, searchName := range searchNames {
        fmt.Printf("[searchItems] Searching for: %s\n", searchName)
        for _, item := range assortment.Items {
            if strings.Contains(strings.ToLower(item.Name), strings.ToLower(searchName)) {
                fmt.Printf("[searchItems] Found matching item: %s\n, PurchasableBalance: %s\n", item.Name, item.PurchasableBalance)
                if item.PurchasableBalance != nil && *item.PurchasableBalance == 0 {
                    continue
                }
                matchingItems[searchName] = append(matchingItems[searchName], item)
            }
        }
    }

    return matchingItems
}

// Compare fetched results with the existing results and send differences to Telegram
func compareResults(newResults []Result, existingResults []Result) {
    for _, newResult := range newResults {
        if existingResult, exists := existingResults[newResult.Name]; exists {
            fmt.Printf("[compareResults] Product: %s\n, old status: %b\n, new status: %b\n", newResult.Name, existingResult.Founded, newResult.Founded)
            if existingResult.Founded != newResult.Founded {
                fmt.Printf("[compareResults] Status changed for product: %s\n", newResult.Name)
                var statusMessage string
                if newResult.Founded {
                    statusMessage = "The product is available!"
                } else {
                    statusMessage = "The product is out of stock..."
                }
                message := fmt.Sprintf("**Product**: %s\n**Venue**: %s\n**Status**: %s\n![image](%s)", newResult.Name, newResult.Venue, statusMessage, newResult.Image)
                err := sendTelegramMessage(message)
                if err != nil {
                    fmt.Printf("Failed to send message: %v\n", err)
                }
            }
        }
    }
}

// Process each config entry by fetching items and searching for matches
func processConfig(config Config) ([]Result, error) {
    assortment, err := fetchItems(config.Endpoint)
    if err != nil {
        return nil, err
    }

    matches := searchItems(assortment, config.Names)
    var results []Result

    for _, searchName := range config.Names {
        items := matches[searchName]
        if len(items) == 0 {
            result := Result{
                Venue:   config.Venue,
                Name:    searchName,
                Founded: false,
                Time:    time.Now().Format(time.RFC3339),
                Image:   "",
            }
            results = append(results, result)
            continue
        }
        foundedItem := items[0] 
        var image string = ""
        if foundedItem.Images != nil && len(foundedItem.Images) > 0 {
            image = foundedItem.Images[0].URL
        }
        result := Result{
            Venue:   config.Venue,
            Name:    searchName,
            Founded: true,
            Time:    time.Now().Format(time.RFC3339),
            Image:   image,
        }
        results = append(results, result)
    }

    return results, nil
}

func main() {
    configs, err := loadConfig("config.json")
    if err != nil {
        fmt.Println("Error loading config:", err)
        return
    }

    // Run the 30-minute loop in a goroutine
    go func() {
        for {
            fmt.Println("[main] Starting new iteration: ", time.Now().Format("02-01 15:04:05"))
            existingResults, err := loadExistingResults()
            if err != nil {
                fmt.Println("Error loading existing results:", err)
                return
            }

            var allResults []Result

            for _, config := range configs {
                fmt.Printf("Start processing config for venue: %s\n", config.Venue)
                results, err := processConfig(config)
                if err != nil {
                    fmt.Println("Error processing config:", err)
                    continue
                }
                fmt.Printf("Config processed successfully for venue: %s\n", config.Venue)
                // Compare the new results with existing results and send to Telegram
                compareResults(results, existingResults)

                allResults = append(allResults, results...)
            }

            err = writeResultsToFile(allResults, "store.json")
            if err != nil {
                fmt.Println("Error writing results to file:", err)
            } else {
                fmt.Println("Results successfully written to store.json")
            }
            time.Sleep(woltParseTimeout * time.Minute)
        }
    }()

    offset := 0
    // Run the 15-seconds loop in another goroutine
    go func() {
		for {
            updates, err := getUpdates(offset)
            if err != nil {
                fmt.Printf("Error getting updates: %v\n", err)
                continue
            }
            for _, update := range updates {
                // Print the received message
                fmt.Printf("New message from %s: %s\n", update.Message.From.Username, update.Message.Text)
                containsStatus := containsStatus(update.Message.Text)
                if containsStatus {
                    fmt.Println("The message contains the word 'status'")
                    existingResults, err := loadExistingResults()
                    if err != nil {
                        fmt.Println("Error loading existing results:", err)
                        return
                    }
                    resultsMessage := formatResultsMessage(existingResults)
                    fmt.Println("loading existing results:", resultsMessage)
                    err = sendTelegramMessage(resultsMessage)
                    if err != nil {
                        fmt.Println("Error sending message:", err)
                    }
                }
                // Update offset to the latest update ID + 1 to avoid fetching the same message again
                offset = update.UpdateID + 1
            }
			// Your code for the 15-seconds task here
			time.Sleep(telegramCheckTimeout * time.Second)
		}
	}()

    // Start the HTTP server in a separate goroutine
    go func() {
        http.HandleFunc("/send-text", httpHandler)
        if err := http.ListenAndServe(":8088", nil); err != nil {
            fmt.Println("Error starting server:", err)
            return
        }
    }()

    // Keep the main function alive with an infinite loop
	select {} // This blocks the main function forever to keep the program running
}
