package main

import (
    "net/http"
	"fmt"
	"io"
	"encoding/json"
	"strings"
)

type CityConfig map[string]struct {
	Name     string
	Location [2]float64
}

var CITIES = CityConfig{
	"tbilisi": {
		Name:     "Tbilisi",
		Location: [2]float64{44.7965812683105, 41.7024604154103}, // longitude, latitude
	},
	"batumi": {
		Name:     "Batumi",
		Location: [2]float64{41.629900932312, 41.6442796894184},
	},
}

type AllStoreResponse struct {
    Sections []Section `json:"sections"`
}

type Section struct {
    Items []StoreItem `json:"items,omitempty"` // 'omitempty' to handle cases where 'items' might be absent
}

type StoreItem struct {
    Image   Image  `json:"image"`
    Title   string `json:"title"`
    TrackID string `json:"track_id"`
	Venue   Venue  `json:"venue"`
}

type Venue struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
	Address	string `json:"address"`
}

type Image struct {
    URL string `json:"url"`
}

// https://consumer-api.wolt.com/v1/pages/retail?lat=41.7024604154103&lon=44.7965812683105

type WoltConnector struct {
	city string
}

func SetWoltConnector(city string) *WoltConnector {
	return &WoltConnector{city: city}
}

func (wc *WoltConnector) getAllStoresInCity() ([]StoreItem, error) {
	// Implementation will go here
	allStoresBaseURL := "https://consumer-api.wolt.com/v1/pages/retail"
	lon := CITIES[wc.city].Location[0]
	lat := CITIES[wc.city].Location[1]
	endpoint := fmt.Sprintf("%s?lat=%f&lon=%f", allStoresBaseURL, lat, lon)
	response, err := http.Get(endpoint)
    if err != nil {
		return nil, fmt.Errorf("error making all stores GET request: %v", err)
    }
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading all stores response body: %v", err)
	}

	var allStoresResponse AllStoreResponse
	err = json.Unmarshal(body, &allStoresResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing all stores JSON response: %v", err)
	}
	fmt.Printf("Alnill stores length: %f\n", len(allStoresResponse.Sections[1].Items))
	if(len(allStoresResponse.Sections) < 2) {
		return nil, fmt.Errorf("no stores for this city: %v", err)
	}
	return allStoresResponse.Sections[1].Items, nil
}

func (wc *WoltConnector) getStoreByName(storeName string) ([]StoreItem, error) {
	allStores, err := wc.getAllStoresInCity()
	if err != nil {
		return nil, fmt.Errorf("error getting all stores: %v", err)
	}
	var filteredStores []StoreItem

    for idx, store := range allStores {
        if strings.Contains(strings.ToLower(store.Title), strings.ToLower(storeName)) {
            filteredStores = append(filteredStores, store)
			fmt.Printf("Store %d: %s\n", idx, store.Title)
        }
    }

    if len(filteredStores) == 0 {
        return nil, fmt.Errorf("no stores found with the name: %s", storeName)
    }

    return filteredStores, nil
}

func (wc *WoltConnector) getCategoryByProduct(storeSlug string, productName string) {
	// Implementation will go here
	categoriesBaseURL := fmt.Sprintf("https://consumer-api.wolt.com/consumer-api/consumer-assortment/v1/venues/slug/%s/assortment", storeSlug)
	response, err := http.Get(categoriesBaseURL)
	if err != nil {
		fmt.Printf("Error fetching categories: %v\n", err)
		return
	}
	defer response.Body.Close()

}

func main() {
	wc := SetWoltConnector("tbilisi")
	stores, err := wc.getStoreByName("gastronome")
	if err != nil {
		fmt.Printf("Error getting store: %v\n", err)
		return
	}
	fmt.Printf("Stores: %v\n", stores)
}
