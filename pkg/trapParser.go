package h0neytr4p

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Taken from https://stackoverflow.com/questions/55300117/how-do-i-find-all-files-that-have-a-certain-extension-in-go-regardless-of-depth
func GetAllTraps(root, ext string) []string {
	var a []string
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			a = append(a, s)
		}
		return nil
	})
	return a
}

func ParseTraps(traps string) []Trap {
	var Traps []Trap
	for _, trap := range GetAllTraps(traps, ".json") {
		trapData := Trap{}
		trapFile, ok := ioutil.ReadFile(trap)
		if ok != nil {
			log.Fatalln("Error reading file: " + trap)
			log.Fatalln("Exiting Now. Check the file and re-run the program. Bye!")
			os.Exit(1)
		}
		ok = json.Unmarshal([]byte(trapFile), &trapData)
		if ok != nil {
			log.Fatalln("Looks like the file: " + trap + " isn't a valid trap.")
			log.Fatalln("Exiting Now. Chec the file and re-run the program. Bye!")
			os.Exit(1)
		}
		Traps = append(Traps, trapData)
	}
	return Traps
}
