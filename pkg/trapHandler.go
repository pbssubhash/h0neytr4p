package h0neytr4p

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Initialize log file
func InitLogFile(filename string, verbose string) {
	var err error
	logFile, err = os.Create(filename)
	Verbose = verbose
	if err != nil {
		log.Fatalln("Error creating the log file:", err)
		os.Exit(1)
	}
	fmt.Println("Logging is configured and ready.")
}

// Initialize payload folder
func InitPayloadFolder(folder string, verbose string) {
	payloadFolder = folder
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		log.Fatalln("Error creating payload folder: ", err)
		os.Exit(1)
	}
	fmt.Println("Payload folder is configured and ready.")
}

func LogEntry(details map[string]string) bool {
	logFileMutex.Lock()
	defer logFileMutex.Unlock()

	// Convert log details to JSON
	entryJSON, err := json.Marshal(details)
	if err != nil {
		log.Println("Error marshalling log entry to JSON:", err)
		return false
	}

	// Write JSON log entry to file
	fmt.Fprintln(logFile, string(entryJSON))

	// Print to console if verbose mode is enabled
	if Verbose == "true" {
		fmt.Printf("[%s] [Path: %s] [Trapped: %s]\n", details["timestamp"], details["request_uri"], details["trapped"])
	}

	logFile.Sync()
	return true
}
