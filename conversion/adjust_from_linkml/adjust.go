package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

func main() {
	filePath := "../fromlinkml/structs.go"
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading the file:", err)
		return
	}
	fileContent := string(content)

	if !strings.Contains(fileContent, `"LS_reader/conversion/basetypes"`) {
		packageEndIndex := strings.Index(fileContent, "\n") + 1

		fileContent = fileContent[:packageEndIndex] + `import "LS_reader/conversion/basetypes"` + "\n\n" + fileContent[packageEndIndex:]
	}

	fileContent = strings.ReplaceAll(fileContent, "int", "basetypes.Int")
	fileContent = strings.ReplaceAll(fileContent, "float64", "basetypes.Float64")
	fileContent = strings.ReplaceAll(fileContent, "bool", "basetypes.Bool")
	fileContent = strings.ReplaceAll(fileContent, "string", "basetypes.String")
	fileContent = strings.ReplaceAll(fileContent, "time.Date", "basetypes.String")
	fileContent = strings.ReplaceAll(fileContent, "QuantityValue", "basetypes.Float64")
	fileContent = strings.ReplaceAll(fileContent, "Any", "basetypes.String")
	fileContent = strings.ReplaceAll(fileContent, "type basetypes.Float64", "type QuantityValue")
	fileContent = strings.ReplaceAll(fileContent, "type basetypes.String", "type Any")

	// making sure any imported OSCEM schema works - provided it is compatible with the conversions table.
	pattern := "oscem-.*"
	replacement := "oscem"
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return
	}
	lines := strings.Split(fileContent, "\n")
	for i, line := range lines {
		if re.MatchString(line) {
			lines[i] = re.ReplaceAllString(line, replacement)
		}
	}
	output := strings.Join(lines, "\n")

	err = ioutil.WriteFile(filePath, []byte(output), 0644)
	if err != nil {
		fmt.Println("Error writing the file:", err)
		return
	}

	fmt.Println("File updated successfully!")
}
