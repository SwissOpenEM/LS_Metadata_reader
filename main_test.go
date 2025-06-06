package main_test

import (
	"LS_reader/LS_Metadata_reader"
	"encoding/json"
	"os"
	"testing"

	conversion "github.com/osc-em/Converter"

	"github.com/stretchr/testify/assert"
)

func TestReaderTableDriven(t *testing.T) {
	readJSONFile := func(filepath string) string {
		data, err := os.ReadFile(filepath)
		if err != nil {
			t.Fatalf("Failed to read expected data file %s: %v", filepath, err)
		}
		return string(data)
	}
	//reader
	targetXML := readJSONFile("./tests/xml_full.json")
	targetMdoc := readJSONFile("./tests/mdocs_full.json")
	targetCombine := readJSONFile("./tests/combine_full.json")
	targetmdocspa := readJSONFile("./tests/mdocspa_full.json")
	targetdepth := readJSONFile("./tests/depthcheck_full.json")
	//converter
	target2XML := readJSONFile("./tests/xml_correct.json")
	target2Mdoc := readJSONFile("./tests/mdocs_correct.json")
	target2Combine := readJSONFile("./tests/combine_correct.json")
	target2mdocspa := readJSONFile("./tests/mdocspa_correct.json")
	target2depth := readJSONFile("./tests/depthcheck_correct.json")

	tests := []struct {
		name       string
		directory  string
		zFlag      bool
		fFlag      bool
		wantData   string // reader only
		wantErr    bool
		wantData2  string // e2e
		p1Flag     string
		p2Flag     string
		p3Flag     string
		folderFlag string
		outF       bool
	}{
		{
			name:       "xmls",
			directory:  "./tests/xml",
			zFlag:      false,
			fFlag:      false,
			wantData:   targetXML,
			wantErr:    false,
			wantData2:  target2XML,
			p1Flag:     "2.7",
			p2Flag:     "none",
			p3Flag:     "",
			folderFlag: "",
			outF:       false,
		},
		{
			name:       "mdocs",
			directory:  "./tests/mdocs",
			zFlag:      false,
			fFlag:      false,
			wantData:   targetMdoc,
			wantErr:    false,
			wantData2:  target2Mdoc,
			p1Flag:     "2.7",
			p2Flag:     "none",
			p3Flag:     "",
			folderFlag: "",
			outF:       false,
		},
		{
			name:       "Both",
			directory:  "./tests/combine",
			zFlag:      false,
			fFlag:      false,
			wantData:   targetCombine,
			wantErr:    false,
			wantData2:  target2Combine,
			p1Flag:     "2.7",
			p2Flag:     "none",
			p3Flag:     "",
			folderFlag: "",
			outF:       false,
		},
		{
			name:       "mdocspa",
			directory:  "./tests/mdocspa",
			zFlag:      false,
			fFlag:      false,
			wantData:   targetmdocspa,
			wantErr:    false,
			wantData2:  target2mdocspa,
			p1Flag:     "2.7",
			p2Flag:     "none",
			p3Flag:     "",
			folderFlag: "",
			outF:       false,
		},
		{
			name:       "depthcheck",
			directory:  "./tests/depthcheck",
			zFlag:      false,
			fFlag:      false,
			wantData:   targetdepth,
			wantErr:    false,
			wantData2:  target2depth,
			p1Flag:     "2.7",
			p2Flag:     "none",
			p3Flag:     "",
			folderFlag: "",
			outF:       false,
		},
		{
			name:       "folderFlag",
			directory:  "./tests/empty",
			zFlag:      false,
			fFlag:      false,
			wantData:   targetXML,
			wantErr:    false,
			wantData2:  target2XML,
			p1Flag:     "2.7",
			p2Flag:     "none",
			p3Flag:     "",
			folderFlag: "myfoldername",
			outF:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := LS_Metadata_reader.Reader(tt.directory, tt.zFlag, tt.fFlag, tt.p3Flag, tt.folderFlag)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Reader() error = %v, wantErr %v", err, tt.wantErr)
			}
			// rerun the json marshalling to ensure no issues with whitespaces etc
			var jsonData map[string]string
			if err := json.Unmarshal(data, &jsonData); err != nil {
				t.Fatalf("Failed to unmarshal returned data: %v", err)
			}
			var jsonDataclean map[string]string
			if err := json.Unmarshal([]byte(tt.wantData), &jsonDataclean); err != nil {
				t.Fatalf("Failed to unmarshal returned data: %v", err)
			}
			// exclude irrelevant keys:
			excludeKeys := []string{"MicroscopeImage.UniqueID", "MicroscopeImage.uniqueID", "MicroscopeImage.microscopeData.core.Guid", "ImageFile", "MinMaxMean", "SubFramePath", "[T"}
			cleaned_json := preprocessMap(jsonData, excludeKeys)
			cleaned_target := preprocessMap(jsonDataclean, excludeKeys)
			actualDataBytes, err := json.Marshal(cleaned_json)
			if err != nil {
				t.Fatalf("Failed to re-marshal returned data: %v", err)
			}
			targetDataBytes, err := json.Marshal(cleaned_target)
			if err != nil {
				t.Fatalf("Failed to re-marshal returned data: %v", err)
			}
			assert.JSONEqf(t, string(targetDataBytes), string(actualDataBytes), "Mismatch in test case %s", tt.name)

			data2, err2 := conversion.Convert(data, "", tt.p1Flag, tt.p2Flag, tt.folderFlag)

			if (err2 != nil) != tt.wantErr {
				t.Fatalf("Reader() error = %v, wantErr %v", err2, tt.wantErr)
			}
			var jsonData2 interface{}
			if err := json.Unmarshal(data2, &jsonData2); err != nil {
				t.Fatalf("Failed to unmarshal returned data: %v", err)
			}

			actualDataBytes2, err := json.Marshal(jsonData2)
			if err != nil {
				t.Fatalf("Failed to re-marshal returned data: %v", err)
			}

			assert.JSONEqf(t, tt.wantData2, string(actualDataBytes2), "Mismatch in test case %s", tt.name)
		})
	}
}
func filterMap(input map[string]string, excludeKeys []string) map[string]string {
	result := make(map[string]string)
	excludeSet := make(map[string]struct{}, len(excludeKeys))
	for _, key := range excludeKeys {
		excludeSet[key] = struct{}{}
	}
	for k, v := range input {
		if _, found := excludeSet[k]; !found {
			result[k] = v
		}
	}
	return result
}
func preprocessMap(input map[string]string, excludeKeys []string) map[string]string {
	return filterMap(input, excludeKeys)
}
