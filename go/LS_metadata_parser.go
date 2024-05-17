package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// XML PART
// Definitions of the xml structure
type MicroscopeImage struct {
	XMLName    xml.Name   `xml:"MicroscopeImage"`
	Name       string     `xml:"name"`
	UniqueID   string     `xml:"uniqueID"`
	CustomData CustomData `xml:"CustomData"`
}

// for key-value
type CustomData struct {
	KeyValues []KeyValue `xml:"KeyValueOfstringanyType"`
}

type KeyValue struct {
	Key   string `xml:"Key"`
	Value string `xml:"Value"`
}

// For tag-value
type Element struct {
	XMLName  xml.Name
	Content  string    `xml:",chardata"`
	Children []Element `xml:",any"`
}

func parseElement(element Element, path string, leafNodes map[string]string) {
	currentPath := path
	if currentPath != "" {
		currentPath += "." + element.XMLName.Local
	} else {
		currentPath = element.XMLName.Local
	}

	trimmedContent := strings.TrimSpace(element.Content)
	if len(element.Children) == 0 && trimmedContent != "" {
		leafNodes[currentPath] = trimmedContent
	}
	for _, child := range element.Children {
		parseElement(child, currentPath, leafNodes)
	}
}

func process_xml(input string) (map[string]string, error) {
	xmlFile, err := os.Open(input)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer xmlFile.Close()

	xmlData, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err
	}
	var root Element
	err = xml.Unmarshal(xmlData, &root)
	if err != nil {
		fmt.Println("Error unmarshaling XML:", err)
		return nil, err
	}

	leafNodes := make(map[string]string)
	parseElement(root, "", leafNodes)

	var image MicroscopeImage
	err = xml.Unmarshal(xmlData, &image)
	if err != nil {
		fmt.Println("Error unmarshaling XML:", err)
		return nil, err
	}
	leafNodes["MicroscopeImage.Name"] = image.Name
	leafNodes["MicroscopeImage.UniqueID"] = image.UniqueID
	for _, kv := range image.CustomData.KeyValues {
		leafNodes[kv.Key] = kv.Value
	}
	return (leafNodes), err
}

// MDOC Part
func process_mdoc(input string) (map[string]string, error) {
	var count float64 = 0.00
	re := regexp.MustCompile(`(.+?)\s*=\s*(.+)`)
	mdocFile, err := os.Open(input)
	if err != nil {
		fmt.Printf("Welp your file didnt open")
		return nil, err
	}
	defer mdocFile.Close()
	scanner := bufio.NewScanner(mdocFile)
	mdoc_results := make(map[string]string)

	for scanner.Scan() {
		// Look for special case
		//TiltAxis Angle
		tiltaxis := strings.Contains(scanner.Text(), "TiltAxisAngle")    // Tomo 5
		tiltaxis2 := strings.Contains(scanner.Text(), "Tilt axis angle") // SerialEM
		if tiltaxis {
			blabb_split := strings.Split(re.FindStringSubmatch(scanner.Text())[2], "=")[1]
			mdoc_results["TiltAxisAngle"] = (strings.TrimSpace(strings.Split(blabb_split, "  ")[0])) // this is bound to fail at some point if they dont keep their double space seperation logic
		}
		if tiltaxis2 {
			blab_split := strings.Split(re.FindStringSubmatch(scanner.Text())[2], ",")[0]
			mdoc_results["TiltAxisAngle"] = (strings.TrimSpace(strings.Split(blab_split, "=")[1]))
		}
		// general search and update for min/max values
		match := re.FindStringSubmatch(scanner.Text())
		// Quick check incase the image dimesions are only present in the header
		image := strings.Contains(scanner.Text(), "ImageSize")
		if image {
			mdoc_results["ImageDimensions_X"] = strings.Split(match[2], " ")[0]
			mdoc_results["ImageDimensions_Y"] = strings.Split(match[2], " ")[1]
		}
		if match != nil {
			if strings.TrimSpace(match[1]) == "[ZValue" {
				count++
			}
			value, exists := mdoc_results[match[1]]
			if !exists {
				mdoc_results[match[1]] = match[2]
			} else if value == match[2] {
				// Grab some Tuples
				energy := strings.Contains(scanner.Text(), "FilterSlitAndLoss")
				if energy {
					energytest, _ := strconv.ParseFloat(strings.TrimSpace(strings.Split(match[2], " ")[0]), 64)
					if energytest > float64(0.00) {
						mdoc_results["EnergyFilterUsed"] = "true"
						mdoc_results["EnergyFilterSlitWidth"] = strings.Split(match[2], " ")[0]
					}
				}
				continue
			} else if value != match[2] {
				test, err := strconv.ParseFloat(strings.TrimSpace(mdoc_results[match[1]]), 64)
				if err != nil {
					// Grab the remaining Tuples
					beamshift := strings.Contains(scanner.Text(), "Beamshift") // check for correct syntax only present in newer versions of SerialEM
					imageShift := strings.Contains(scanner.Text(), "ImageShift")
					stagepos := strings.Contains(scanner.Text(), "StagePosition")
					if stagepos {
						xcheck, xexist := mdoc_results["Stage_x_max"]
						ycheck, yexist := mdoc_results["Stage_y_max"]
						if !xexist && !yexist {
							mdoc_results["Stage_x_max"] = strings.Split(match[2], " ")[0]
							mdoc_results["Stage_y_max"] = strings.Split(match[2], " ")[1]
							mdoc_results["Stage_x_min"] = strings.Split(match[2], " ")[0]
							mdoc_results["Stage_y_min"] = strings.Split(match[2], " ")[1]
						} else {
							xtest, _ := strconv.ParseFloat(strings.TrimSpace(xcheck), 64)
							ytest, _ := strconv.ParseFloat(strings.TrimSpace(ycheck), 64)
							x_new, _ := strconv.ParseFloat(strings.TrimSpace(strings.Split(match[2], " ")[0]), 64)
							y_new, _ := strconv.ParseFloat(strings.TrimSpace(strings.Split(match[2], " ")[1]), 64)
							mdoc_results["Stage_x_max"] = strconv.FormatFloat(max(xtest, x_new), 'f', 2, 64)
							mdoc_results["Stage_y_max"] = strconv.FormatFloat(max(ytest, y_new), 'f', 2, 64)
							mdoc_results["Stage_x_min"] = strconv.FormatFloat(min(xtest, x_new), 'f', 2, 64)
							mdoc_results["Stage_y_min"] = strconv.FormatFloat(min(ytest, y_new), 'f', 2, 64)
						}
					}
					if beamshift {
						xcheck, xexist := mdoc_results["Beamshift_x_max"]
						ycheck, yexist := mdoc_results["Beamshift_y_max"]
						if !xexist && !yexist {
							mdoc_results["Beamshift_x_max"] = strings.Split(match[2], " ")[0]
							mdoc_results["Beamshift_y_max"] = strings.Split(match[2], " ")[1]
							mdoc_results["Beamshift_x_min"] = strings.Split(match[2], " ")[0]
							mdoc_results["Beamshift_y_min"] = strings.Split(match[2], " ")[1]
						} else {
							xtest, _ := strconv.ParseFloat(strings.TrimSpace(xcheck), 64)
							ytest, _ := strconv.ParseFloat(strings.TrimSpace(ycheck), 64)
							x_new, _ := strconv.ParseFloat(strings.TrimSpace(strings.Split(match[2], " ")[0]), 64)
							y_new, _ := strconv.ParseFloat(strings.TrimSpace(strings.Split(match[2], " ")[1]), 64)
							mdoc_results["Beamshift_x_max"] = strconv.FormatFloat(max(xtest, x_new), 'f', 2, 64)
							mdoc_results["Beamshift_y_max"] = strconv.FormatFloat(max(ytest, y_new), 'f', 2, 64)
							mdoc_results["Beamshift_x_min"] = strconv.FormatFloat(min(xtest, x_new), 'f', 2, 64)
							mdoc_results["Beamshift_y_min"] = strconv.FormatFloat(min(ytest, y_new), 'f', 2, 64)
						}
					}
					if imageShift {
						xcheck, xexist := mdoc_results["ImageShift_x_max"]
						ycheck, yexist := mdoc_results["ImageShift_y_max"]
						if !xexist && !yexist {
							mdoc_results["ImageShift_x_max"] = strings.Split(match[2], " ")[0]
							mdoc_results["ImageShift_y_max"] = strings.Split(match[2], " ")[1]
							mdoc_results["ImageShift_x_min"] = strings.Split(match[2], " ")[0]
							mdoc_results["ImageShift_y_min"] = strings.Split(match[2], " ")[1]
						} else {
							xtest, _ := strconv.ParseFloat(strings.TrimSpace(xcheck), 64)
							ytest, _ := strconv.ParseFloat(strings.TrimSpace(ycheck), 64)
							x_new, _ := strconv.ParseFloat(strings.TrimSpace(strings.Split(match[2], " ")[0]), 64)
							y_new, _ := strconv.ParseFloat(strings.TrimSpace(strings.Split(match[2], " ")[1]), 64)
							mdoc_results["ImageShift_x_max"] = strconv.FormatFloat(max(xtest, x_new), 'f', 2, 64)
							mdoc_results["ImageShift_y_max"] = strconv.FormatFloat(max(ytest, y_new), 'f', 2, 64)
							mdoc_results["ImageShift_x_min"] = strconv.FormatFloat(min(xtest, x_new), 'f', 2, 64)
							mdoc_results["ImageShift_y_min"] = strconv.FormatFloat(min(ytest, y_new), 'f', 2, 64)
						}
					}
					continue
				} else {
					new, _ := strconv.ParseFloat(strings.TrimSpace(match[2]), 64)
					keymin, existmin := mdoc_results[match[1]+"_min"]
					keymax, existmax := mdoc_results[match[1]+"_max"]
					if !existmin {
						mdoc_results[match[1]+"_min"] = strconv.FormatFloat(min(test, new), 'f', 2, 64)
					} else {
						oldmin, _ := strconv.ParseFloat(strings.TrimSpace(keymin), 64)
						mdoc_results[match[1]+"_min"] = strconv.FormatFloat(min(new, oldmin), 'f', 2, 64)
					}
					if !existmax {
						mdoc_results[match[1]+"_max"] = strconv.FormatFloat(max(test, new), 'f', 2, 64)
					} else {
						oldmax, _ := strconv.ParseFloat(strings.TrimSpace(keymax), 64)
						mdoc_results[match[1]+"_max"] = strconv.FormatFloat(max(new, oldmax), 'f', 2, 64)
					}
				}
			}
		}

	}
	// Numberoftilts
	mdoc_results["NumberOfTilts"] = strconv.FormatFloat(count, 'f', 2, 64)

	// get tiltangle at the end if applicable
	_, existtilt := mdoc_results["TiltAngle"]
	if existtilt {
		tiltmax, err := strconv.ParseFloat(strings.TrimSpace(mdoc_results["TiltAngle_max"]), 64)
		if err != nil {
			fmt.Println("Tilt angle increment calculation failed")
		}
		tiltmin, err := strconv.ParseFloat(strings.TrimSpace(mdoc_results["TiltAngle_min"]), 64)
		if err != nil {
			fmt.Println("Tilt angle increment calculation failed")
		}
		mdoc_results["Tilt_increment"] = strconv.FormatFloat(math.Abs(tiltmax-tiltmin)/count, 'f', 2, 64)
	}
	// Software used
	T, T_exist := mdoc_results["[T"]
	if T_exist {
		if strings.Contains(T, "TOMOGRAPHY") || strings.Contains(T, "Tomography") {
			mdoc_results["Software"] = "EPU-Tomo5"
		} else if strings.Contains(T, "SerialEM") {
			mdoc_results["Software"] = "SerialEM"
		}
	} // generalized before, if SerialEM additions/scripts were used:
	vers, vers_exist := mdoc_results["Version"]
	if vers_exist {
		mdoc_results["Software"] = vers
	}
	// Inference based things come here
	dark, darkexist := mdoc_results["DarkField"]
	if darkexist {
		te, _ := strconv.Atoi(strings.TrimSpace(dark))
		if te == 1 {
			mdoc_results["Imaging"] = "Darkfield"
		}
	}
	mag, magexist := mdoc_results["MagIndex"]
	if magexist {
		te, _ := strconv.Atoi(strings.TrimSpace(mag))
		te2, _ := strconv.Atoi(strings.TrimSpace(dark))
		if te > 0 && (te2 == 0 || !darkexist) {
			mdoc_results["Imaging"] = "Brightfield"
		}
	}
	// Currently missing Illumination modes (EMDB allowed: "Flood Beam", "Spot Scan", "Other") --
	// Problem how to differentiate Spot Scan ; most cryoEM cases definitely Flood Beam
	// Could do "Flood Beam" as baseline and add a catch later; dont know if anyone uses serialEM for spotscan anyways
	EMMode, modeexist := mdoc_results["EMmode"]
	if modeexist {
		te, _ := strconv.Atoi(strings.TrimSpace(EMMode))
		if te == 0 {
			mdoc_results["EMMode"] = "TEM"
		} else if te == 1 {
			mdoc_results["EMMode"] = "EFTEM"
		} else if te == 2 {
			mdoc_results["EMMode"] = "STEM"
		} else if te == 3 {
			mdoc_results["Imaging"] = "Diffraction"
		}
	}
	// Cleanup before return
	for key := range mdoc_results {
		_, upexist := mdoc_results[key+"_max"]
		_, dwnexist := mdoc_results[key+"_min"]
		if upexist || dwnexist {
			delete(mdoc_results, key)
		}
	}
	//return
	return mdoc_results, err
}

// MERGE and datetimechecks
func merge_to_dataset_level(listofcontents []map[string]string) map[string]string {
	overallmap := make(map[string]string)
	timeformats := []string{
		"02-Jan-06  15:04:05",
		"02-Jan-2006  15:04:05",
		"2006-Jan-02  15:04:05",
		time.RFC3339Nano,
	}
	for item := range listofcontents {
		for key := range listofcontents[item] {
			value, exists := overallmap[key]
			valuenew := (listofcontents[item])[key]
			if !exists {
				overallmap[key] = valuenew
			} else if value == valuenew {
				continue
			} else if value != valuenew {
				if strings.Contains(key, "DateTime") {
					for _, datetime := range timeformats {
						timecheck, err1 := time.Parse(datetime, value)
						timechecknew, err := time.Parse(datetime, valuenew)
						if err == nil && err1 == nil {
							_, existstart := overallmap[key+"_start"]
							_, existend := overallmap[key+"_end"]
							if !existstart {
								if timecheck.After(timechecknew) {
									overallmap[key+"_start"] = timechecknew.Format(time.RFC3339)
								} else {
									overallmap[key+"_start"] = timecheck.Format(time.RFC3339)
								}
							} else {
								timecheckold, _ := time.Parse(time.RFC3339, overallmap[key+"_start"])
								if timecheckold.After(timechecknew) {
									overallmap[key+"_start"] = timechecknew.Format(time.RFC3339)
								}
							}
							if !existend {
								if timecheck.Before(timechecknew) {
									overallmap[key+"_end"] = timechecknew.Format(time.RFC3339)
								} else {
									overallmap[key+"_end"] = timecheck.Format(time.RFC3339)
								}
							} else {
								timecheckold, _ := time.Parse(time.RFC3339, overallmap[key+"_end"])
								if timecheckold.Before(timechecknew) {
									overallmap[key+"_end"] = timechecknew.Format(time.RFC3339)
								}
							}
						}
					}
				}
				test, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
				if err != nil {
					// Get the beam and stage tuples from mdocs across multiple images
					beamshift := strings.Contains(key, "Beamshift") // check for correct syntax only present in newer versions of SerialEM
					imageShift := strings.Contains(key, "ImageShift")
					stagepos := strings.Contains(key, "StagePosition")
					if stagepos {
						xcheck, xexist := overallmap["Stage_x_max"]
						ycheck, yexist := overallmap["Stage_y_max"]
						if !xexist && !yexist {
							overallmap["Stage_x_max"] = strings.Split(valuenew, " ")[0]
							overallmap["Stage_y_max"] = strings.Split(valuenew, " ")[1]
							overallmap["Stage_x_min"] = strings.Split(valuenew, " ")[0]
							overallmap["Stage_y_min"] = strings.Split(valuenew, " ")[1]
						} else {
							xtest, _ := strconv.ParseFloat(strings.TrimSpace(xcheck), 64)
							ytest, _ := strconv.ParseFloat(strings.TrimSpace(ycheck), 64)
							x_new, _ := strconv.ParseFloat(strings.TrimSpace(strings.Split(valuenew, " ")[0]), 64)
							y_new, _ := strconv.ParseFloat(strings.TrimSpace(strings.Split(valuenew, " ")[1]), 64)
							overallmap["Stage_x_max"] = strconv.FormatFloat(max(xtest, x_new), 'f', 2, 64)
							overallmap["Stage_y_max"] = strconv.FormatFloat(max(ytest, y_new), 'f', 2, 64)
							overallmap["Stage_x_min"] = strconv.FormatFloat(min(xtest, x_new), 'f', 2, 64)
							overallmap["Stage_y_min"] = strconv.FormatFloat(min(ytest, y_new), 'f', 2, 64)
						}
					}
					if beamshift {
						xcheck, xexist := overallmap["Beamshift_x_max"]
						ycheck, yexist := overallmap["Beamshift_y_max"]
						if !xexist && !yexist {
							overallmap["Beamshift_x_max"] = strings.Split(valuenew, " ")[0]
							overallmap["Beamshift_y_max"] = strings.Split(valuenew, " ")[1]
							overallmap["Beamshift_x_min"] = strings.Split(valuenew, " ")[0]
							overallmap["Beamshift_y_min"] = strings.Split(valuenew, " ")[1]
						} else {
							xtest, _ := strconv.ParseFloat(strings.TrimSpace(xcheck), 64)
							ytest, _ := strconv.ParseFloat(strings.TrimSpace(ycheck), 64)
							x_new, _ := strconv.ParseFloat(strings.TrimSpace(strings.Split(valuenew, " ")[0]), 64)
							y_new, _ := strconv.ParseFloat(strings.TrimSpace(strings.Split(valuenew, " ")[1]), 64)
							overallmap["Beamshift_x_max"] = strconv.FormatFloat(max(xtest, x_new), 'f', 2, 64)
							overallmap["Beamshift_y_max"] = strconv.FormatFloat(max(ytest, y_new), 'f', 2, 64)
							overallmap["Beamshift_x_min"] = strconv.FormatFloat(min(xtest, x_new), 'f', 2, 64)
							overallmap["Beamshift_y_min"] = strconv.FormatFloat(min(ytest, y_new), 'f', 2, 64)
						}
					}
					if imageShift {
						xcheck, xexist := overallmap["ImageShift_x_max"]
						ycheck, yexist := overallmap["ImageShift_y_max"]
						if !xexist && !yexist {
							overallmap["ImageShift_x_max"] = strings.Split(valuenew, " ")[0]
							overallmap["ImageShift_y_max"] = strings.Split(valuenew, " ")[1]
							overallmap["ImageShift_x_min"] = strings.Split(valuenew, " ")[0]
							overallmap["ImageShift_y_min"] = strings.Split(valuenew, " ")[1]
						} else {
							xtest, _ := strconv.ParseFloat(strings.TrimSpace(xcheck), 64)
							ytest, _ := strconv.ParseFloat(strings.TrimSpace(ycheck), 64)
							x_new, _ := strconv.ParseFloat(strings.TrimSpace(strings.Split(valuenew, " ")[0]), 64)
							y_new, _ := strconv.ParseFloat(strings.TrimSpace(strings.Split(valuenew, " ")[1]), 64)
							overallmap["ImageShift_x_max"] = strconv.FormatFloat(max(xtest, x_new), 'f', 2, 64)
							overallmap["ImageShift_y_max"] = strconv.FormatFloat(max(ytest, y_new), 'f', 2, 64)
							overallmap["ImageShift_x_min"] = strconv.FormatFloat(min(xtest, x_new), 'f', 2, 64)
							overallmap["ImageShift_y_min"] = strconv.FormatFloat(min(ytest, y_new), 'f', 2, 64)
						}
					}
					continue
				} else {
					new, _ := strconv.ParseFloat(strings.TrimSpace(valuenew), 64)
					keymin, existmin := overallmap[key+"_min"]
					keymax, existmax := overallmap[key+"_max"]
					if !existmin {
						overallmap[key+"_min"] = strconv.FormatFloat(min(test, new), 'f', 2, 64)
					} else {
						oldmin, _ := strconv.ParseFloat(strings.TrimSpace(keymin), 64)
						overallmap[key+"_min"] = strconv.FormatFloat(min(new, oldmin), 'f', 2, 64)
					}
					if !existmax {
						overallmap[key+"_max"] = strconv.FormatFloat(max(test, new), 'f', 2, 64)
					} else {
						oldmax, _ := strconv.ParseFloat(strings.TrimSpace(keymax), 64)
						overallmap[key+"_max"] = strconv.FormatFloat(max(new, oldmax), 'f', 2, 64)
					}
				}
			}
		}
	}
	for key := range overallmap {
		_, upexist := overallmap[key+"_max"]
		_, dwnexist := overallmap[key+"_min"]
		_, startexist := overallmap[key+"_start"]
		_, endexist := overallmap[key+"_end"]
		if upexist || dwnexist || startexist || endexist {
			delete(overallmap, key)
		}
	}
	overallmap["NumberOfMovies"] = strconv.Itoa(len(listofcontents))
	return overallmap
}

// Readin - evaluation
func readin(directory string) ([]map[string]string, []map[string]string, error) {
	var xmlContents []map[string]string
	var mdocContents []map[string]string

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, nil, err
	}
	for _, file := range files {
		if !file.IsDir() && !isHidden(file.Name()) {
			filePath := filepath.Join(directory, file.Name())
			switch filepath.Ext(file.Name()) {
			case ".xml":
				xmlContent, err := process_xml(filePath)
				if err == nil {
					xmlContents = append(xmlContents, xmlContent)
				} else {
					fmt.Println("Import of ", filePath, " failed")
				}
			case ".mdoc":
				mdocContent, err := process_mdoc(filePath)
				if err == nil {
					mdocContents = append(mdocContents, mdocContent)
				} else {
					fmt.Println("Import of ", filePath, " failed")
				}
			}
		}
	}

	return mdocContents, xmlContents, err
}

// minicheck against hidden files
func isHidden(name string) bool {
	return len(name) > 0 && name[0] == '.'
}

func main() {
	// Check that the number of arguments fits expectation
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run Metadata_extractor.go <directory>")
		return
	}

	// Get the directory from the command-line argument
	directory := os.Args[1]

	// Check if the provided directory exists
	fileInfo, err := os.Stat(directory)
	if os.IsNotExist(err) {
		fmt.Printf("Error: Directory '%s' does not exist.\n", directory)
		return
	}

	// Check if the provided path is a directory
	if !fileInfo.IsDir() {
		fmt.Printf("Error: '%s' is not a directory.\n", directory)
		return
	}

	mdoc_files, xml_files, err := readin(os.Args[1])
	if err != nil {
		fmt.Println("Are you sure this was the correct directory?", err)
		return
	}
	var out map[string]string
	if mdoc_files != nil && xml_files == nil {
		out = merge_to_dataset_level(mdoc_files)
	} else if xml_files != nil && mdoc_files == nil {
		out = merge_to_dataset_level(xml_files)
	} else if xml_files != nil && mdoc_files != nil {
		out := merge_to_dataset_level(mdoc_files)
		b := merge_to_dataset_level(xml_files)
		for x, y := range b {
			out[x] = y
		}
	} else {
		fmt.Println("Something went wrong, nothing was read out")
		os.Exit(1)
	}

	jsonData, err := json.MarshalIndent(out, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return
	}
	nameout1, _ := filepath.Abs(directory)
	counter := strings.Split(nameout1, "/")
	var nameout string
	if len(counter) > 0 {
		nameout = counter[len(counter)-1] + ".json"
	} else {
		fmt.Println("Name generation failed, returning to default")
		nameout = "Dataset_out.json"
	}
	err = os.WriteFile(nameout, jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file:", err)
		return
	}

	fmt.Println("Extracted data has been written to ", nameout)
}
