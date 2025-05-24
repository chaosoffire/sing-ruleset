package convert

import (
	"fmt"
	"log"
	"os/exec"
)

func ConvertWithSingBox(sourceFile string, targetFile string, aType string) error {
	var cmd *exec.Cmd
	switch aType {
	case "adguard-srs":
		cmd = exec.Command("sing-box", "rule-set", "convert", "--type", "adguard", "--output", targetFile, sourceFile)
	default:
		cmd = exec.Command("sing-box", "rule-set", "compile", "--output", targetFile, sourceFile)
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("error executing sing-box command: %v", cmd.Args)
		log.Printf("output: %s", output)
		return fmt.Errorf("error converting file %v to %v: %v", sourceFile, targetFile, err)
	}
	return nil
}
