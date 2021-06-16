package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"

	h0neytr4p "github.com/pbssubhash/h0neytr4p"
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
                                         /$$  | $$                              | $$      
       Built by a Red team, with <3     |  $$$$$$/                              | $$      
             h0neytr4p v0.1             \______/                               |__/      
        Built by zer0p1k4chu & g0dsky
    https://github.com/pbssubhash/h0neyt4p
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
	output := flag.String("output", "Default", "Output file - It's a string.")
	log := flag.String("log", "Default", "Log file - It's a string.")
	Verbose := flag.String("verbose", "true", "Use -verbose=false for disabling streaming output; by default it's true")
	help := flag.String("help", "Print Help", "Print Help")
	flag.Parse()
	if *trapsFolder == "Default" || isFlagPassed(*help) || *output == "Default" || *log == "Default" {
		fmt.Println("Wrong Arguments.. Exiting Now")
		flag.PrintDefaults()
		os.Exit(1)
	}
	trapConfig := h0neytr4p.ParseTraps(*trapsFolder)
	h0neytr4p.CreateTrapFile(*output)
	h0neytr4p.CreateLogFile(*log, *Verbose)
	var ports []string
	filteredTraps := make(map[string][]h0neytr4p.Trap)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			// sig is a ^C, handle it
			if sig == os.Interrupt {
				fmt.Println("Interrupt received. Gracefully exiting the program.")
				os.Exit(1)
			}
		}
	}()
	for _, trap := range trapConfig {
		v, found := filteredTraps[trap.Basicinfo.Port]
		if found {
			v = append(v, trap)
			filteredTraps[trap.Basicinfo.Port] = v
		} else {
			v = append(v, trap)
			filteredTraps[trap.Basicinfo.Port] = v
			ports = append(ports, trap.Basicinfo.Port)
		}
	}
	for _, port := range ports {
		wg.Add(1)
		go h0neytr4p.StartHandler(port, filteredTraps[port])
	}
	wg.Wait()
}
