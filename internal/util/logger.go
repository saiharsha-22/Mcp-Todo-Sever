package util

import (
	"log"
	"os"
)

func LogStartupInfo(todoFilePath string) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Failed to get current working directory: %v\n", err)
	} else {
		log.Printf("Current working directory: %s\n", cwd)
	}
	log.Printf("Using todo file path: %s\n", todoFilePath)
}

func JoinLines(lines []string) string {
	out := ""
	for i, line := range lines {
		if i > 0 {
			out += "\n"
		}
		out += line
	}
	return out
}
