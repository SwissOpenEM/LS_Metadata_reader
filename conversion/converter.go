package conversion

import (
	"LS_reader/configuration"
	"LS_reader/conversion/basetypes"
	oscem "LS_reader/conversion/fromlinkml"
	"embed"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/stoewer/go-strcase"
)

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

func Convert(jsonin []byte, content embed.FS, p1Flag string, p2Flag string) ([]byte, error) {

	csvRecords, err := readCSVFile(content)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading CSV file:", err)
		return nil, err
	}

	var jsonData map[string]string
	err = json.Unmarshal(jsonin, &jsonData)
	if err != nil {
		return nil, err
	}

	var testing oscem.Instrument
	var acq_testing oscem.AcquisitionTomo

	for k, v := range jsonData {
		for _, test := range csvRecords {
			if test.crunchfromxml != "" && test.fromxml == k {
				v, err = unitcrunch(v, test.crunchfromxml)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Unit crunching failed: ", err)
					continue
				}
			}
			if test.crunchfrommdoc != "" && (test.frommdoc == k || test.optionals_mdoc == k) {
				v, err = unitcrunch(v, test.crunchfrommdoc)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Unit crunching failed: ", err)
					continue
				}
			}
			prio := false
			if test.fromxml == k || test.frommdoc == k || test.optionals_mdoc == k || test.optionals_xml == k {
				if test.frommdoc == k || test.optionals_mdoc == k {
					prio = true
				}
				testing, acq_testing = untangle(test, prio, testing, acq_testing, v)
			}
		}
	}

	// Set some defaults in a config file
	var fixvalues map[string]string
	if p1Flag == "" && p2Flag == "" {
		config, err := configuration.Getconfig()
		if err != nil {
			_ = err
			//could enable the warning below but its not really required
			//fmt.Fprintln(os.Stderr, "Warning: config was not set or could not be obtained - make sure the config is set using LS_Metadata_reader --c or you are using the flags")
		} else {
			errun := json.Unmarshal(config, &fixvalues)
			if errun != nil {
				_ = errun
				//fmt.Fprintln(os.Stderr, "Warning: Config was found but unreadable")
			}
		}
	} else {
		fixvalues = make(map[string]string)
		fixvalues["CS"] = p1Flag
		fixvalues["Gainref_FlipRotate"] = p2Flag
	}

	SetField(&acq_testing, "", "GainrefFlipRotate", fixvalues["Gainref_FlipRotate"], "", false)
	SetField(&testing, "", "Cs", fixvalues["CS"], "mm", false)
	//
	mh := make(map[string]interface{})
	mh["instrument"] = testing
	mh["acquisition"] = acq_testing
	// Filter out fields that are nil
	wut, err := json.Marshal(mh)
	if err != nil {
		return nil, err
	}
	// this allows us to obtain nil values for types where Go usually doesnt allow them e.g. int
	var kek interface{}
	err = json.Unmarshal(wut, &kek)
	if err != nil {
		return nil, err
	}
	cleaned := CleanMap(kek)
	out, err := json.MarshalIndent(cleaned, "", "   ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Json generation failed ", err)
	}
	return out, nil
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

type csvextract struct {
	OSCEM          string
	fromxml        string
	frommdoc       string
	optionals_mdoc string
	units          string
	crunchfromxml  string
	crunchfrommdoc string
	optionals_xml  string
}

func readCSVFile(content embed.FS) ([]csvextract, error) {
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
	desiredColumns := []string{"OSCEM", "fromxml", "frommdoc", "optionals_mdoc", "units", "crunchfromxml", "crunchfrommdoc", "optionals_xml"}

	header := records[0]
	// correct csv bullshit
	for i, name := range header {
		header[i] = strings.TrimLeft(name, "\ufeff")
	}
	columnIndices := make(map[string]int)
	for idx, colName := range header {
		columnIndices[colName] = idx
	}
	// Check if desired columns exist
	for _, col := range desiredColumns {
		if _, ok := columnIndices[col]; !ok {
			log.Fatalf("Column %s not found in header", col)
		}
	}
	var bestextract []csvextract

	for _, row := range records[1:] {
		data := csvextract{
			OSCEM:          row[columnIndices["OSCEM"]],
			fromxml:        row[columnIndices["fromxml"]],
			frommdoc:       row[columnIndices["frommdoc"]],
			optionals_mdoc: row[columnIndices["optionals_mdoc"]],
			units:          row[columnIndices["units"]],
			crunchfromxml:  row[columnIndices["crunchfromxml"]],
			crunchfrommdoc: row[columnIndices["crunchfrommdoc"]],
			optionals_xml:  row[columnIndices["optionals_xml"]],
		}
		bestextract = append(bestextract, data)
	}
	return bestextract, nil
}

func unitcrunch(v string, fac string) (string, error) {
	check, err := strconv.ParseFloat(v, 64)
	factor, _ := strconv.ParseFloat(fac, 64)
	if err != nil {
		return v, err
	}
	val := check * factor
	back := strconv.FormatFloat(val, 'f', 16, 64)

	return back, nil
}

func untangle(coll csvextract, prio bool, testing oscem.Instrument, acq_testing oscem.AcquisitionTomo, v string) (oscem.Instrument, oscem.AcquisitionTomo) {
	// untangle the . seperation and send it against the OSCEM schema struct for field setting
	untang := strings.Split(coll.OSCEM, ".")
	length := len(untang)
	if length > 2 {
		err := SetField(&testing, strcase.UpperCamelCase(untang[length-2]), strcase.UpperCamelCase(untang[length-1]), v, coll.units, prio)
		if err != nil {
			err2 := SetField(&acq_testing, strcase.UpperCamelCase(untang[length-2]), strcase.UpperCamelCase(untang[length-1]), v, coll.units, prio)
			if err2 != nil {
				fmt.Println(err, err2)
			}
		}
	} else {
		err := SetField(&testing, "", strcase.UpperCamelCase(untang[length-1]), v, coll.units, prio)
		if err != nil {
			err2 := SetField(&acq_testing, "", strcase.UpperCamelCase(untang[length-1]), v, coll.units, prio)
			if err2 != nil {
				fmt.Println(err, err2)
			}
		}
	}
	return testing, acq_testing
}
