package h0neytr4p

import (
	"fmt"
	"log"
	"os"
)

// Frequency based file sync -- for large scale? Need to work on a feasible solution for this.
// func SyncFiles() {
// 	var wg sync.WaitGroup
// 	for {
// 		go func() {
// 			wg.Add(1)
// 			time.Sleep(5 * time.Second)
// 			logFileGlobal.Sync()
// 			filenameGlobal.Sync()
// 			time.Sleep(5 * time.Second)
// 			wg.Done()
// 		}()
// 		wg.Wait()
// 	}
// }
func CreateLogFile(filename string, verbose string) {
	f, e := os.Create(filename)
	Verbose = verbose
	if e != nil {
		log.Fatalln("Error with creating the file")
		os.Exit(1)
	}
	logFileGlobal = f
	fmt.Fprintln(f, "SourceIP,UserAgent,Timestamp,Path,Trapped")
	fmt.Println("Logging is configured and ready. Enjoying Trapping..")
}
func CreateTrapFile(filename string) {
	f, e := os.Create(filename)
	if e != nil {
		log.Fatalln("Error with creating the file")
		os.Exit(1)
	}
	fmt.Fprintln(f, "SourceIP,UserAgent,TrappedFor,RiskRating,References")
	filenameGlobal = f
	fmt.Println("Output is configured and ready.")
}

func LogIt(details LogEntry) bool {
	// f, e := os.Open(logFileGlobal)
	// if e != nil {
	// 	log.Fatalln("Error with opening the log file..")
	// 	os.Exit(1)
	// }
	fmt.Fprintln(logFileGlobal, details.SourceIP+","+details.UserAgent+","+details.Timestamp+","+details.Path+","+details.Trapped)
	if Verbose == "true" {
		fmt.Println("[" + details.Timestamp + "] Path: " + details.Path + "; Trapped: " + details.Trapped)
	}
	logFileGlobal.Sync()
	return true
}
func TrapIt(details Attacker) bool {
	// // f, e := os.Open(filenameGlobal)
	// if e != nil {
	// 	log.Fatalln("Error with opening the file..")
	// 	os.Exit(1)
	// }
	fmt.Fprintln(filenameGlobal, details.SourceIP+","+details.UserAgent+","+details.TrappedFor+","+details.RiskRating+","+details.References)
	filenameGlobal.Sync()
	return true
}
