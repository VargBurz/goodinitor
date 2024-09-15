package main

import (
    "encoding/json"
    "fmt"
    "os"
	"io/ioutil"
)


// Load and parse the config.json file into a slice of Config structs
func loadConfig(filename string) ([]Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %v", err)
	}
	defer file.Close()

	configData, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var configs []Config
	err = json.Unmarshal(configData, &configs)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	return configs, nil
}

// Write results to a file in the specified format
func writeResultsToFile(results []Result, filename string) error {
    fmt.Printf("[writeResultsToFile] Writing results to %s\n", filename)
    data, err := json.MarshalIndent(results, "", "  ")
    if err != nil {
        return fmt.Errorf("error marshalling results: %v", err)
    }

    err = ioutil.WriteFile(filename, data, 0644)
    if err != nil {
        return fmt.Errorf("error writing results to file: %v", err)
    }

    return nil
}

// Load existing results from store.json
func getResultsMap() (map[string]Result, error) {
    existingResults, err := loadExistingResults()
    if err != nil {
        return nil, fmt.Errorf("error loading existing results: %v", err)
    }

    existingMap := make(map[string]Result)
    for _, res := range existingResults {
        existingMap[res.Name] = res
    }

    return existingMap, nil
}

// Load existing results from store.json
func loadExistingResults() ([]Result, error) {
    file, err := os.Open("store.json")
    if err != nil {
        return nil, fmt.Errorf("error opening store file: %v", err)
    }
    defer file.Close()

    var existingResults []Result
    data, err := ioutil.ReadAll(file)
    if err != nil {
        return nil, fmt.Errorf("error reading store file: %v", err)
    }

    err = json.Unmarshal(data, &existingResults)
    if err != nil {
        return nil, fmt.Errorf("error parsing store file: %v", err)
    }

    return existingResults, nil
}