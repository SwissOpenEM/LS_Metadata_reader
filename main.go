package main

import (
	"LS_reader/LS_Metadata_reader"
	"LS_reader/configuration"
	"LS_reader/conversion"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
)

//go:embed conversion/conversions.csv
var content embed.FS

func main() {

	zFlag := flag.Bool("z", false, "Toggle whether to make a zip archive of all xml files - default: false")
	fFlag := flag.Bool("f", false, "Toggle whether the full metadata is also written out in addition to the OSCEM schema conform one- default: false")
	cFlag := flag.Bool("c", false, "If you want to reset your config file")
	flag.Parse()
	posArgs := flag.Args()

	// allow for reconfiguration of the config
	if *cFlag {
		current := configuration.Getconfig()
		var grid map[string]string
		err := json.Unmarshal(current, &grid)
		if err != nil {
			fmt.Println("Config exists but reading it failed", err)
		}
		fmt.Println("current config:\n", grid)
		configuration.Changeconfig()
	}
	// Check that there are arguments
	if len(posArgs) == 0 {
		fmt.Println("No arguments; correct usage: ./LS_reader --z --f <directory>")
		return
	}

	// Get the directory from the command-line argument
	directory := posArgs[0]
	data, err := LS_Metadata_reader.Reader(directory, *zFlag, *fFlag)
	if err != nil {
		fmt.Println("The extraction went wrong due to", err)
	}
	err1 := conversion.Convert(data, content)
	if err1 != nil {
		fmt.Println("The extraction went wrong due to", err1)
	}
}
