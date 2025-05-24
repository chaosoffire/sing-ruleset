package file

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Adguard_Blocklists []BlocklistEntry `json:"Adguard_Blocklists"`
	IP_Lists           []BlocklistEntry `json:"IP_Lists"`
}

type BlocklistEntry struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func ReadJson(filePath string) (Config, error) {
	// Open the JSON file
	file, err := os.Open(filePath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()
	// Decode the JSON data into the Config struct
	decoder := json.NewDecoder(file)
	var config Config
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, err
	}
	// Return the populated Config struct
	return config, nil
}

func GetConfigs(config Config) (map[string]string, map[string]string) {
	// Create a map to store the blocklist names and URLs
	Adguard_Blocklists_Map := make(map[string]string)
	IP_Lists_Map := make(map[string]string)
	// Iterate over the Adguard Blocklists
	for _, blocklist := range config.Adguard_Blocklists {
		Adguard_Blocklists_Map[blocklist.Name] = blocklist.URL
	}
	// Iterate over the IP Lists
	for _, blocklist := range config.IP_Lists {
		IP_Lists_Map[blocklist.Name] = blocklist.URL
	}
	return Adguard_Blocklists_Map, IP_Lists_Map
}

func PrintConfig(config Config) {
	// Print the Adguard Blocklists
	for _, blocklist := range config.Adguard_Blocklists {
		log.Println("Adguard Blocklist Name:", blocklist.Name)
		log.Println("Adguard Blocklist URL:", blocklist.URL)
	}
	// Print the IP Lists
	for _, blocklist := range config.IP_Lists {
		log.Println("IP List Name:", blocklist.Name)
		log.Println("IP List URL:", blocklist.URL)
	}
}
