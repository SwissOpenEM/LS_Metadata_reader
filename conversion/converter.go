package conversion

import (
	"LS_reader/configuration"
	"LS_reader/conversion/basetypes"
	"LS_reader/conversion/generated"
	"embed"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type CSVRecord struct {
	Field1  string
	Field2  string
	Field3  string
	Field4  string
	Field5  string
	Field6  string
	Field7  string
	Field8  string
	Field9  string
	Field10 string
	Field11 string
	Field12 string
	Field13 string
	Field14 string
	Field15 string
}

func SetField(obj interface{}, parent, name, value, unit string, prio bool) error {
	structValue := reflect.ValueOf(obj).Elem()
	return traverseAndSetField(structValue, parent, name, value, unit, prio)
}

func traverseAndSetField(structValue reflect.Value, parent, name, value, unit string, prio bool) error {
	if parent == "" {
		if err := setFieldValue(structValue, name, value, unit, prio); err == nil {
			return nil
		}
	}
	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)
		structField := structValue.Type().Field(i)

		if field.Kind() == reflect.Ptr {
			// Initialize the nested structs
			if field.IsNil() {
				field.Set(reflect.New(field.Type().Elem()))
			}
			field = field.Elem()
		}

		if field.Kind() == reflect.Struct {
			if structField.Name == parent {
				if err := setFieldValue(field, name, value, unit, prio); err == nil {
					return nil
				}
			}
			// traverse the nested structs
			if err := traverseAndSetField(field, parent, name, value, unit, prio); err == nil {
				return nil
			}
		}
	}

	return fmt.Errorf("no such field: %s in obj with parent: %s", name, parent)
}

// Set field value with custom types
func setFieldValue(structValue reflect.Value, name, value, unit string, prio bool) error {
	fieldValue := structValue.FieldByName(name)
	if !fieldValue.IsValid() {
		return fmt.Errorf("no such field: %s", name)
	}

	if !fieldValue.CanSet() {
		return fmt.Errorf("cannot set field %s", name)
	}
	test := fieldValue.FieldByName("HasSet")
	if test.Bool() {
		if !prio {
			return nil
		}
	}
	// Attempt conversions for custom types
	switch fieldValue.Interface().(type) {
	case basetypes.Int:
		intValue, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
		if err == nil {
			fieldValue.Set(reflect.ValueOf(basetypes.Int{Value: intValue, HasSet: true, Unit: unit}))
			return nil
		} else { // if the original software writes i.e. 300.00 for Acceleration Voltage
			floattest, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
			if err == nil {
				fieldValue.Set(reflect.ValueOf(basetypes.Int{Value: int64(floattest), HasSet: true, Unit: unit}))
				return nil
			}
		}
	case basetypes.Float64:
		floatValue, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
		if err == nil {
			fieldValue.Set(reflect.ValueOf(basetypes.Float64{Value: floatValue, HasSet: true, Unit: unit}))
			return nil
		}
	case basetypes.Bool:
		boolValue, err := strconv.ParseBool(strings.TrimSpace(value))
		if err == nil {
			fieldValue.Set(reflect.ValueOf(basetypes.Bool{Value: boolValue, HasSet: true}))
			return nil
		}
	case basetypes.String:
		fieldValue.Set(reflect.ValueOf(basetypes.String{Value: value, HasSet: true}))
		return nil
	default:
		return fmt.Errorf("unsupported kind %s", fieldValue.Kind())
	}

	return fmt.Errorf("provided value %s could not be converted to the appropriate type", value)
}

func Convert(jsonin []byte, content embed.FS) error {

	csvRecords, err := readCSVFile(content)
	if err != nil {
		fmt.Println("Error reading CSV file:", err)
		return err
	}

	var jsonData map[string]string
	err = json.Unmarshal(jsonin, &jsonData)
	if err != nil {
		return err
	}

	var testing generated.Instrument
	mh := make(map[string]generated.Instrument)

	for k, v := range jsonData {
		for _, test := range csvRecords {
			if test.Field14 != "" && test.Field5 == k {
				v, err = unitcrunch(v, test)
				if err != nil {
					fmt.Println("Unit crunching failed: ", err)
					continue
				}
			}
			prio := false
			if test.Field5 == k || test.Field6 == k || test.Field13 == k {
				if test.Field6 == k {
					prio = true
				}
				if test.Field2 != "" {
					if test.Field3 != "" {
						if test.Field4 != "" {
							err := SetField(&testing, test.Field3, test.Field4, v, test.Field14, prio)
							if err != nil {
								fmt.Println(err)
							}
						} else {
							err := SetField(&testing, test.Field2, test.Field3, v, test.Field14, prio)
							if err != nil {
								fmt.Println(err)
							}
						}
					} else {
						err := SetField(&testing, "", test.Field2, v, test.Field14, prio)
						if err != nil {
							fmt.Println(err)
						}
					}
				}
			}
		}
	}

	// Set some defaults in a config file
	var fixvalues map[string]string
	errun := json.Unmarshal(configuration.Getconfig(), &fixvalues)
	if errun != nil {
		fmt.Println("config was not set and could not be obtained - make sure the config is set at ~/.config/LS_reader.conf")
	}
	SetField(&testing, "", "GainRef_FlipRotate", fixvalues["Gainref_FlipRotate"], "", false)
	SetField(&testing, "", "CS", fixvalues["CS"], "mm", false)
	//
	mh["Instrument"] = testing
	// Filter out fields that are nil
	wut, err := json.Marshal(mh)
	if err != nil {
		return err
	}
	// this allows us to obtain nil values for types where Go usually doesnt allow them i.e. int
	var kek interface{}
	err = json.Unmarshal(wut, &kek)
	if err != nil {
		return err
	}
	cleaned := CleanMap(kek)
	out, err := json.MarshalIndent(cleaned, "", "   ")
	cwd, _ := os.Getwd()
	cut := strings.Split(cwd, string(os.PathSeparator))
	name := cut[len(cut)-1] + ".json"
	if err != nil {
		fmt.Print(err)
	} else {
		os.WriteFile(name, out, 0644)
		fmt.Println("Extracted data was written to: ", name)
	}
	return nil
}

// CleanMap removes nil values from a map of maps
func CleanMap(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		cleanedMap := make(map[string]interface{})
		for key, value := range v {
			cleanedValue := CleanMap(value)
			if cleanedValue != nil {
				cleanedMap[key] = cleanedValue
			}
		}
		if len(cleanedMap) == 0 {
			return nil
		}
		return cleanedMap
	default:
		if v == nil {
			return nil
		}
		return v
	}
}

// readCSVFile reads and parses a CSV file into a slice of CSVRecord structs
func readCSVFile(content embed.FS) ([]CSVRecord, error) {
	file, err := content.Open("conversion/conversions.csv")
	if err != nil {
		return nil, err
	}

	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var csvRecords []CSVRecord
	for _, record := range records {
		if len(record) < 13 {
			return nil, fmt.Errorf("invalid record: %v", record)
		}
		csvRecords = append(csvRecords, CSVRecord{
			Field1:  record[0],
			Field2:  record[1],
			Field3:  record[2],
			Field4:  record[3],
			Field5:  record[4],
			Field6:  record[5],
			Field7:  record[6],
			Field8:  record[7],
			Field9:  record[8],
			Field10: record[9],
			Field11: record[10],
			Field12: record[11],
			Field13: record[12],
			Field14: record[13],
			Field15: record[14],
		})
	}
	return csvRecords, nil
}
func unitcrunch(v string, test CSVRecord) (string, error) {
	check, err := strconv.ParseFloat(v, 64)
	factor, _ := strconv.ParseFloat(test.Field15, 64)
	if err != nil {
		return v, err
	}
	val := check * factor
	back := strconv.FormatFloat(val, 'f', 16, 64)

	return back, nil
}
