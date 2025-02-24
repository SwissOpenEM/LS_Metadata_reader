package main

import (
	"LS_reader/LS_Metadata_reader"
	"LS_reader/configuration"
	"LS_reader/conversion"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

//go:embed conversion/conversions.csv
var content embed.FS

func main() {
	//for benchmarking
	/*f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := trace.Start(f); err != nil {
		panic(err)
	}
	defer trace.Stop()*/

	zFlag := flag.Bool("z", false, "Toggle whether to make a zip archive of all xml files - default: false")
	fFlag := flag.Bool("f", false, "Toggle whether the full metadata is also written out in addition to the OSCEM schema conform one- default: false")
	cFlag := flag.Bool("c", false, "If you want to reset your config file")
	oFlag := flag.String("o", "", "Provide target output path and name for your metadata file, leave empty to write to current working directory")
	iFlag := flag.String("i", "", "Provide target input folder - will take first positional argument if --i is missing")
	p1Flag := flag.String("cs", "", "Provide CS value here, if you dont want to use configs")
	p2Flag := flag.String("gain_flip_rotate", "", "Provide whether and how to flip the gain ref here, if you dont want to use configs")
	p3Flag := flag.String("epu", "", "Provide the path to the mirrored EPU folder containing all the xmls of the datacollections here, if you dont want to use configs")
	flag.Parse()
	posArgs := flag.Args()

	// allow for reconfiguration of the config
	if *cFlag {
		current, err := configuration.Getconfig()
		var grid map[string]string
		if err != nil {
			fmt.Fprintln(os.Stderr, " No prior config obtainable", err)
		}
		_ = json.Unmarshal(current, &grid)
		fmt.Println("current config:\n", grid)
		configuration.Changeconfig()
		return
	}
	var directory string
	// Check that there are arguments
	if len(posArgs) == 0 && *iFlag == "" {
		fmt.Println("No arguments; correct minimum arguments: ./LS_reader <directory>")
		return
	} else if *iFlag != "" {
		directory = *iFlag
	} else {
		directory = posArgs[0]
	}

	data, err := LS_Metadata_reader.Reader(directory, *zFlag, *fFlag, *p3Flag)
	if err != nil {
		fmt.Fprintln(os.Stderr, "The extraction went wrong due to", err)
	}
	output, err1 := conversion.Convert(data, content, *p1Flag, *p2Flag)
	if err1 != nil {
		fmt.Fprintln(os.Stderr, "The extraction went wrong due to", err1)
	}
	if *oFlag == "" {
		cwd, _ := os.Getwd()
		cut := strings.Split(cwd, string(os.PathSeparator))
		name := cut[len(cut)-1] + ".json"
		os.WriteFile(name, output, 0644)
		fmt.Println()
		fmt.Println("Extracted data was written to: ", name)

	} else {
		twd := *oFlag
		if !strings.Contains(twd, ".json") {
			var conc []string
			conc = append(conc, twd, "json")
			twd = strings.Join(conc, ".")
		}
		os.WriteFile(twd, output, 0644)
		fmt.Println()
		fmt.Printf("Extracted data was written to: %s", twd)
	}
}
