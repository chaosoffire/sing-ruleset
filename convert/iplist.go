package convert

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

type Config struct {
	Version int    `json:"version"`
	Rules   []Rule `json:"rules"`
}

type Rule struct {
	IpCidr []string `json:"ip_cidr,omitempty"`
}

func NewConfigFromIpList(ipList []string) *Config {
	return &Config{
		Version: 1,
		Rules: []Rule{
			{
				IpCidr: ipList,
			},
		},
	}
}

// func NewIpListFromFile(sourceFile string) ([]string, error) {
// 	// Read the IP list from the source file
// 	file, err := os.Open(sourceFile)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to open source file %v: %w", sourceFile, err)
// 	}
// 	defer file.Close()

// 	var IpList []string
// 	scanner := bufio.NewScanner(file)
// 	for scanner.Scan() {
// 		if line := strings.TrimSpace(scanner.Text()); line != "" {
// 			IpList = append(IpList, line)
// 		}
// 	}
// 	if err := scanner.Err(); err != nil {
// 		return nil, fmt.Errorf("failed to read source file %v: %w", sourceFile, err)
// 	}

// 	var indexs []int
// 	for index, line := range IpList {
// 		flag := false
// 		if net.ParseIP(line) != nil {
// 			flag = true
// 		}
// 		if _, ipcidr, _ := net.ParseCIDR(line); ipcidr != nil {
// 			flag = true
// 		}
// 		if !flag {
// 			indexs = append(indexs, index)
// 		}
// 	}
// 	if len(indexs) > 0 {
// 		for _, index := range indexs {
// 			IpList = append(IpList[:index], IpList[index+1:]...)
// 		}
// 	}

// 	return IpList, nil
// }

func ConvertFromIpList(sourceFile string, targetFile string) error {
	// Read the IP list from the source file
	file, err := os.Open(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to open source file %v: %w", sourceFile, err)
	}
	defer file.Close()

	var IpList []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if line := strings.TrimSpace(scanner.Text()); line != "" {
			IpList = append(IpList, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read source file %v: %w", sourceFile, err)
	}

	var indexs []int
	for index, line := range IpList {
		flag := false
		if net.ParseIP(line) != nil {
			flag = true
		}
		if _, ipcidr, _ := net.ParseCIDR(line); ipcidr != nil {
			flag = true
		}
		if flag {
			indexs = append(indexs, index)
		}
	}
	var validIpList []string
	if len(indexs) > 0 {
		validIpList = make([]string, len(indexs))
		for i, index := range indexs {
			validIpList[i] = IpList[index]
		}
	} else {
		return fmt.Errorf("no valid IPs found in the source file %v", sourceFile)
	}

	config := NewConfigFromIpList(validIpList)

	// Write the config to the target file
	targetFileHandle, err := os.Create(targetFile)
	if err != nil {
		return fmt.Errorf("failed to create target file %v: %w", targetFile, err)
	}
	defer targetFileHandle.Close()
	writer := bufio.NewWriter(targetFileHandle)
	JSONData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config to JSON: %w", err)
	}
	_, err = writer.Write(JSONData)
	if err != nil {
		return fmt.Errorf("failed to write JSON data to target file %v: %w", targetFile, err)
	}
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}
	// Close the target file handle
	err = targetFileHandle.Close()
	if err != nil {
		return fmt.Errorf("failed to close target file %v: %w", targetFile, err)
	}
	targetFileSrs := strings.TrimSuffix(targetFile, ".json") + ".srs"
	err = ConvertWithSingBox(targetFile, targetFileSrs, "iplist")
	if err != nil {
		return fmt.Errorf("failed to convert file %v: %w", targetFile, err)
	}
	// Remove the original target file
	// err = os.Remove(targetFile)
	// if err != nil {
	// 	return fmt.Errorf("failed to remove original target file %v: %w", targetFile, err)
	// }
	return nil
}
