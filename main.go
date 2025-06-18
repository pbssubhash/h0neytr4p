package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"

	h0neytr4p "github.com/t3chn0m4g3/h0neytr4p"
)

func PrintBanner() {
	fmt.Println(`                                                                
 /$$        /$$$$$$                                  /$$               /$$   /$$          
| $$       /$$$_  $$                                | $$              | $$  | $$          
| $$$$$$$ | $$$$\ $$ /$$$$$$$   /$$$$$$  /$$   /$$ /$$$$$$    /$$$$$$ | $$  | $$  /$$$$$$ 
| $$__  $$| $$ $$ $$| $$__  $$ /$$__  $$| $$  | $$|_  $$_/   /$$__  $$| $$$$$$$$ /$$__  $$
| $$  \ $$| $$\ $$$$| $$  \ $$| $$$$$$$$| $$  | $$  | $$    | $$  \__/|_____  $$| $$  \ $$
| $$  | $$| $$ \ $$$| $$  | $$| $$_____/| $$  | $$  | $$ /$$| $$            | $$| $$  | $$
| $$  | $$|  $$$$$$/| $$  | $$|  $$$$$$$|  $$$$$$$  |  $$$$/| $$            | $$| $$$$$$$/
|__/  |__/ \______/ |__/  |__/ \_______/ \____  $$   \___/  |__/            |__/| $$____/ 
                                        /$$  | $$                               | $$      
                                       |  $$$$$$/                [ v0.32 ]      | $$      
                                        \______/                                |__/      
Built by a Red team, with <3
Built by zer0p1k4chu & g0dsky (https://github.com/pbssubhash/h0neytr4p)
Adjusted for T-Pot by t3chn0m4g3 (https://github.com/t3chn0m4g3/h0neytr4p)
	`)
}

// Taken from https://stackoverflow.com/questions/35809252/check-if-flag-was-provided-in-go/35809400
func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func main() {
	PrintBanner()
	var wg sync.WaitGroup
	trapsFolder := flag.String("traps", "Default", "Traps folder - It's a string.")
	log := flag.String("log", "Default", "Log file - It's a string.")
	catchall := flag.String("catchall", "true", "Catch all or only trap based payloads - It's a string.")
	payload := flag.String("payload", "Default", "Payload folder - It's a string.")
	cert := flag.String("cert", "Default", "Certificate File")
	key := flag.String("key", "Default", "Certificate File")
	Verbose := flag.String("verbose", "true", "Use -verbose=false for disabling streaming output; by default it's true")
	wildcard := flag.String("wildcard", "false", "Load all traps on ports 80 and 443 - It's a string.")
	help := flag.String("help", "Print Help", "Print Help")
	flag.Parse()

	if *trapsFolder == "Default" || isFlagPassed(*help) || *log == "Default" || *payload == "Default" {
		fmt.Println("Wrong Arguments.. Exiting Now")
		flag.PrintDefaults()
		os.Exit(1)
	}
	fmt.Printf("[ Traps folder            ] -> [ %-30s]\n", *trapsFolder)
	fmt.Printf("[ Logfile                 ] -> [ %-30s]\n", *log)
	fmt.Printf("[ Payloads folder         ] -> [ %-30s]\n", *payload)
	fmt.Printf("[ Catch all payloads      ] -> [ %-30s]\n", *catchall)
	fmt.Printf("[ Payload multipart limit ] -> [ %-30d]\n", h0neytr4p.MaxMultipartSize)
	fmt.Printf("[ Payload other limit     ] -> [ %-30d]\n", h0neytr4p.MaxJSONFormSize)
	fmt.Printf("[ Wildcard mode           ] -> [ %-30s]\n", *wildcard)
	fmt.Println()
	trapConfig := h0neytr4p.ParseTraps(*trapsFolder)
	h0neytr4p.InitLogFile(*log, *Verbose)
	h0neytr4p.InitPayloadFolder(*payload, *Verbose)

	// List of ports to load each trap for if wildcard mode is enabled
	predefinedPorts := []string{"443", "80"}
	filteredTraps := make(map[string][]h0neytr4p.Trap)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			if sig == os.Interrupt {
				fmt.Println("Interrupt received. Gracefully exiting the program.")
				os.Exit(1)
			}
		}
	}()

	// Load traps according to the --wildcard flag
	for _, trap := range trapConfig {
		if *wildcard == "true" {
			// Add each trap to all predefined ports
			for _, port := range predefinedPorts {
				filteredTraps[port] = append(filteredTraps[port], trap)
			}
		} else {
			// Add trap only to its specified port
			port := trap.Basicinfo.Port
			filteredTraps[port] = append(filteredTraps[port], trap)
		}
	}

	// Start handlers based on the selected ports
	portsToLoad := predefinedPorts
	if *wildcard == "false" {
		// Only load ports specified in traps if wildcard is false
		portsToLoad = make([]string, 0, len(filteredTraps))
		for port := range filteredTraps {
			portsToLoad = append(portsToLoad, port)
		}
	}

	for _, port := range portsToLoad {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			h0neytr4p.StartHandler(p, filteredTraps[p], *cert, *key, *catchall)
		}(port)
	}

	wg.Wait()
}
