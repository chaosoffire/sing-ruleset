package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sing-ruleset/convert"
	"sing-ruleset/file"
	"sing-ruleset/network"
	"sync"
)

func main() {
	// Set the working directory to the current directory
	workPath, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Current working directory: %s\n", workPath)
	// Check if the config file exists
	configPath := filepath.Join(workPath, "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("config file does not exist")
	}
	log.Printf("Config file path: %s\n", configPath)
	// Read the config file
	Config, err := file.ReadJson(configPath)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	fmt.Println("Config file read successfully.")
	// Print the config file
	file.PrintConfig(Config)

	Adguard_Blocklists_Map, IP_Lists_Map := file.GetConfigs(Config)

	log.Println("Start downloading files...")
	totalFiles := len(Adguard_Blocklists_Map) + len(IP_Lists_Map)
	log.Printf("Total files to download: %d\n", totalFiles)

	// Downloading Adguard Blocklists txt files
	log.Println("Start downloading Adguard Blocklists...")
	outputPath := filepath.Join(workPath, "output", "Adguard_Blocklists")
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
			log.Fatalln(err)
		}
	}
	waitDownloadABM := sync.WaitGroup{}
	for name, url := range Adguard_Blocklists_Map {
		filePath := filepath.Join(outputPath, name+".txt")
		waitDownloadABM.Add(1)
		go func() {
			defer waitDownloadABM.Done()
			log.Println("Downloading file:", name)
			if err := network.DownloadFileFromURL(url, filePath); err != nil {
				log.Println("Error downloading file:", err)
			}
		}()
	}
	// Downloading IP Lists txt files
	log.Println("Start downloading IP Lists...")
	outputPath = filepath.Join(workPath, "output", "IP_Lists")
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
			log.Fatalln(err)
		}
	}
	waitDownloadIPM := sync.WaitGroup{}
	for name, url := range IP_Lists_Map {
		filePath := filepath.Join(outputPath, name+".txt")
		waitDownloadIPM.Add(1)
		go func() {
			defer waitDownloadIPM.Done()
			log.Println("Downloading file:", name)
			if err := network.DownloadFileFromURL(url, filePath); err != nil {
				log.Println("Error downloading file:", err)
			}
		}()
	}
	wait := sync.WaitGroup{}
	wait.Add(2)
	go func() {
		waitDownloadABM.Wait()
		wait.Done()
	}()
	go func() {
		waitDownloadIPM.Wait()
		wait.Done()
	}()
	wait.Wait()
	log.Println("All downloads completed successfully.")

	// Check if the sing-box binary exists
	cmd := exec.Command("sing-box", "version")
	_, err = cmd.Output()
	if err != nil {
		log.Fatalln("sing-box binary not found. Please install it first.")
	}

	outputPath = filepath.Join(workPath + "output")
	// if _, err := os.Stat(outputPath); os.IsNotExist(err) {
	// 	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
	// 		log.Fatalln(err)
	// 	}
	// }

	// Start converting files (Adguard_Blocklists)
	log.Println("Start converting files (Adguard Blocklists)...")
	// downloadPath = filepath.Join(workPath, "downloads", "Adguard_Blocklists")
	outputPath = filepath.Join(workPath, "output", "Adguard_Blocklists")
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
			log.Fatalln(err)
		}
	}
	waitConvert := sync.WaitGroup{}
	for name := range Adguard_Blocklists_Map {
		waitConvert.Add(1)
		go func(name string) {
			defer waitConvert.Done()
			sourceFilePath := filepath.Join(outputPath, name+".txt")
			targetFilePath := filepath.Join(outputPath, name+".srs")
			err := convert.ConvertFromAdguard(sourceFilePath, targetFilePath)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("File converted successfully:", name)
		}(name)
	}

	waitConvert.Wait()
	// Start converting files (IP_Lists)
	log.Println("Start converting files (IP Lists)...")
	// downloadPath = filepath.Join(workPath, "downloads", "IP_Lists")
	outputPath = filepath.Join(workPath, "output", "IP_Lists")
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
			log.Fatalln(err)
		}
	}
	for name := range IP_Lists_Map {
		waitConvert.Add(1)
		go func(name string) {
			defer waitConvert.Done()
			sourceFilePath := filepath.Join(outputPath, name+".txt")
			targetFilePath := filepath.Join(outputPath, name+".json")
			err := convert.ConvertFromIpList(sourceFilePath, targetFilePath)
			if err != nil {
				log.Println("Error converting file:", err)
				return
			}
			log.Println("File converted successfully:", name)
		}(name)
	}
	waitConvert.Wait()
	log.Println("All files converted successfully.")
}
