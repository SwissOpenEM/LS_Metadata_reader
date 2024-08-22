package main

import (
	"fmt"
	"io/ioutil"
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
	fileContent = strings.ReplaceAll(fileContent, "oscem-schemas", "oscem")
	fileContent = strings.ReplaceAll(fileContent, "time.Date", "basetypes.String")

	err = ioutil.WriteFile(filePath, []byte(fileContent), 0644)
	if err != nil {
		fmt.Println("Error writing the file:", err)
		return
	}

	fmt.Println("File updated successfully!")
}
