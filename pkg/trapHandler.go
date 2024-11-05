package h0neytr4p

import (
	"fmt"
	"log"
	"os"
)

// Initialize log file with a unified header
func InitLogFile(filename string, verbose string) {
	var err error
	logFile, err = os.Create(filename)
	Verbose = verbose
	if err != nil {
		log.Fatalln("Error creating the log file:", err)
		os.Exit(1)
	}

	// Unified CSV header for both types of log entries
	fmt.Fprintln(logFile, "SourceIP,UserAgent,Timestamp,Path,Trapped,TrappedFor,RiskRating,References")
	fmt.Println("Logging is configured and ready.")
}

// Unified logging function
func LogEntry(details LogDetails) bool {
	logFileMutex.Lock()
	defer logFileMutex.Unlock()

	// Format the log entry line
	entry := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s",
		details.SourceIP,
		details.UserAgent,
		details.Timestamp,
		details.Path,
		details.Trapped,
		details.TrappedFor,
		details.RiskRating,
		details.References,
	)

	// Write to file and print to console if verbose mode is enabled
	fmt.Fprintln(logFile, entry)
	if Verbose == "true" {
		fmt.Println("[" + details.Timestamp + "] Path: " + details.Path + "; Trapped: " + details.Trapped)
	}
	logFile.Sync()
	return true
}
