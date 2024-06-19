package main

import (
	"LS_reader/LS_Metadata_reader"
	"LS_reader/conversion"
	"embed"
	"flag"
	"fmt"
)

//go:embed conversion/conversions.csv
var content embed.FS

func main() {

	zFlag := flag.Bool("z", false, "Toggle whether to make a zip archive of all xml files - default: false")
	fFlag := flag.Bool("f", false, "Toggle whether the full metadata is also written out in addition to the OSCEM schema conform one- default: false")
	flag.Parse()
	posArgs := flag.Args()

	// Check that there are arguments
	if len(posArgs) == 0 {
		fmt.Println("Usage: ./LS_reader.go --z --f <directory>")
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
