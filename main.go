package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/SwissOpenEM/LS_Metadata_reader/LS_Metadata_reader"
	"github.com/SwissOpenEM/LS_Metadata_reader/configuration"

	conversion "github.com/osc-em/Converter"
)

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
	metadataFolder := flag.String("folder_filter", "", "If the system deviates from standard EPU naming conventions, a regex for the folder name with the metadata files can be provided.")
	outF := flag.Bool("cli_out", false, "If you want the results also as a stdout")
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

	current, err := configuration.Getconfig()
	var grid map[string]string
	if err == nil && *p1Flag == "" && *p2Flag == "" && *p3Flag == "" {
		_ = json.Unmarshal(current, &grid)
		*p1Flag = grid["cs"]
		*p2Flag = grid["gainref_flip_rotate"]
		*p3Flag = grid["MPCPATH"]
	}

	data, err := LS_Metadata_reader.Reader(directory, *zFlag, *fFlag, *p3Flag, *metadataFolder)
	if err != nil {
		fmt.Fprintln(os.Stderr, "The extraction went wrong due to", err)
		os.Exit(1)
	}
	out, err1 := conversion.Convert(data, "", *p1Flag, *p2Flag, *oFlag)
	if err1 != nil {
		fmt.Fprintln(os.Stderr, "The extraction went wrong due to", err1)
		os.Exit(1)
	}
	if *outF {
		fmt.Printf("%s", string(out))
	}
}
