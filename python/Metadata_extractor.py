import os, sys, re 
import json, codecs
from lxml import etree
from datetime import datetime

def read_in_files(directory):
    xml_contents = []
    mdoc_contents = []
    numberofitems = 0
    for file in os.listdir(directory):
        if not file.startswith('.'):
            if file.endswith(".xml"):
                file_path = os.path.join(directory, file)
                xml_contents.append(parse_xml(file_path))
                numberofitems +=1
            elif file.endswith(".mdoc"):
                file_path = os.path.join(directory, file)
                mdoc_contents.append(extract_mdoc(file_path)) 
                numberofitems +=1
    return mdoc_contents, xml_contents, numberofitems


### MDOC part
def extract_mdoc(fil):
    key_value_pairs = {}
    numb = 0
    # Define the pattern for extracting key-value pairs
    pattern = r'(\w+)\s*=\s*(.+)'
    with open(fil,"r" ) as f_in:
            for lnumber, line in enumerate(f_in):
                match = re.match(pattern, line)
                if match:
                    # Extract the key and value from the match
                    key = match.group(1) 
                    value = match.group(2).strip()
                    
                    #check for existence and update accordingly
                    if key in key_value_pairs:
                        ### if constant just keep
                        if key_value_pairs[key] == value:
                            continue
                        ### if multiple values exist (tomography only) - get min and max
                        if type(key_value_pairs[key]) == str:
                            if key+"_min" in key_value_pairs:
                                key_value_pairs[key+"_min"] = min(float(key_value_pairs[key+"_min"]), float(value))
                            if key+"_max" in key_value_pairs:
                                key_value_pairs[key+"_max"] = max(float(key_value_pairs[key+"_max"]), float(value))
                            else:
                                try:
                                    key_value_pairs[key+"_min"] = min(float(key_value_pairs[key]), float(value))
                                    key_value_pairs[key+"_max"] = max(float(key_value_pairs[key]), float(value))
                                except:
                                    TypeError
                            
                            # get tilt angle increment if applicable
                            if "TiltAngle" in key:
                                try:
                                    key_value_pairs["Tilt_increment"] = abs(float(key_value_pairs[key+"_max"]) - float(key_value_pairs[key+"_min"])) / key_value_pairs["NumberOfTilts"]
                                except:
                                    ValueError

                    else:         
                        # Add the key-value pair to the dictionary
                        key_value_pairs[key] = value
                # Extract tilt axis from header
                if "Tilt axis angle" in line or "TiltAxisAngle" in line:
                    try:
                        test_axis = (line.split("=")[2]).split(",")[0]
                        key_value_pairs["TiltAxisAngle"] = float( test_axis.strip())
                    except:
                        ValueError
                    try:
                        test_axis = re.split(r'[A-Za-z]', (line.split("=")[2]))[0]
                        key_value_pairs["TiltAxisAngle"] = float(test_axis.strip())
                    except:
                        ValueError
                # count number of tilts in tiltseries
                if "ZValue" in line:
                    test_z = (line.split("=")[1]).split("]")[0]
                    if int (test_z.strip() ) > numb:
                        numb = int( test_z.strip())
                    key_value_pairs["NumberOfTilts"] = numb

                 ##Inference based values:
                  ## Imaging Modes  
                if "MagIndex" in line: ## ADD PART TO DIFFERENTIATE WITH DARKFIELD AFTER SERIALEM UPDATE
                    if float(line.split("=")[-1] > 0):
                        key_value_pairs["ImagingMode"] = "Brightfield"
                if ("CameraLenght" in line) and ((line.split("=")[-1]) != "NaN"):
                    key_value_pairs["ImagingMode"] = "Diffraction"
                if "DarkField" in line: ### Works only for TFS scopes - see serial EM documentation
                    if float(line.split("=")[-1]== 1):
                        key_value_pairs["ImagingMode"] = "Darkfield"
                 ## Illumination modes - EMDB calls these only as "Flood Beam", "Spot Scan" and "Other" translation to be done later, feels slightly off anyways
                if "EMMode" in line: 
                    if float(line.split("=")[-1]== 0):
                        key_value_pairs["EMmode"] = "TEM"
                    elif float(line.split("=")[-1]== 1):
                        key_value_pairs["EMmode"] = "EFTEM"
                    elif float(line.split("=")[-1]== 2):    
                        key_value_pairs["EMmode"] = "STEM" 
                    elif float(line.split("=")[-1]== 3):
                        key_value_pairs["EMmode"] = "Diffraction"


        #remove useless last recorded value for items with min/max
    checklist = []    
    for k in key_value_pairs:
        if k+"_max" in key_value_pairs or k+"_min" in key_value_pairs:
            checklist.append(k)
    for k in checklist:
           del key_value_pairs[k]        
    return key_value_pairs


### XML part

#### Extract all Metadata stored in a Key-Value relation
####
def get_key_value(element, level=0, key_value_dict=None):
    
    # Check if the element is a terminal (leaf) element
    if len(element) == 0:
        
        # Check if the element contains "Key" in its tag
        if "Key" in element.tag:
            # Check if there is a sibling with "Value" in the tag
            sibling_value = element.xpath("following-sibling::*[contains(local-name(), 'Value')]")
            
            if sibling_value:
                # Get the text values of both the Key and Value elements
                key_text = element.text
                value_text = sibling_value[0].text
                
                # Check if neither text attribute is "None"
                if key_text is not None and value_text is not None:
                    key_value_dict[key_text] = value_text
    
    for child in element:
        get_key_value(child, level + 1, key_value_dict)

### Extract all Metadata stored in a Tag-Value relation 
###
def get_tag_value(element, tag_text_dict):
    tag = element.tag
    text = element.text
    
    if text is not None:
        # Check if the element's tag doesn't contain "Key" or "Value"
        if "Key" not in tag and (tag.split("}")[-1] != "Value"):
            parts = tag.split("}")
            key = parts[-1]
            parent_key = ((element.getparent()).tag).split("}")[-1] 
            
            # Check if the key or key+parent_key is already present
            combined_key = key if parent_key is None else f"{parent_key}+{key}"
            
            # Make two sets of values which would otherwise be indistinguishable useable
            if "numericValue" in combined_key:  
                grandparent_key = (((element.getparent()).getparent()).tag).split("}")[-1]
                combined_key = combined_key if grandparent_key is None else f"{grandparent_key}+{combined_key}"    

            if combined_key in tag_text_dict:
                # Append the text to the existing entry
                existing_text = tag_text_dict[combined_key]
                tag_text_dict[combined_key] = existing_text + ";" + text
            else:
                tag_text_dict[combined_key] = text

def search_and_process_tree(root):
    tag_text_dict = {}
    
    for element in root.iter():
        if len(element) == 0:  # Check if it's a terminal (leaf) element
            get_tag_value(element, tag_text_dict)
    
    return tag_text_dict

# Parse the XML data using lxml
def parse_xml(file_in):

    test = etree.parse(file_in)
    root = test.getroot()

    key_value_dict = {}
    # Search and process the XML for Tag-Value pairs
    tag_text_dict = search_and_process_tree(root)
    # Search and process the XML for Key-Value pairs
    get_key_value(root, key_value_dict=key_value_dict)
    ### Merge the results:
    merge_results = {**key_value_dict, **tag_text_dict}

    return merge_results


### Full dataset:
def check_all_files(Listofimage_metadata):


    dataset_dict = {}
    for data in Listofimage_metadata:       
                    # Check if the key exists in the dict
                    for key in data:
                        if not key in dataset_dict and key != "DateTime":
                            dataset_dict[key]=data[key]
                        if not key in dataset_dict and key == "DateTime":
                            try:
                                time = datetime.strptime(data[key], "%y-%b-%d %H:%M:%S")
                            except:
                                ValueError
                            try:
                                time = datetime.strptime(data[key], "%Y-%b-%d %H:%M:%S")
                            except:
                                ValueError
                            try:
                                time = datetime.strptime(data[key], "%d-%b-%Y %H:%M:%S")
                            except:
                                ValueError
                            dataset_dict[key]= time.isoformat(timespec='minutes')
                        
                        if (key in dataset_dict and dataset_dict[key] != data[key]): 
                            # Big effort due to non standartized time formats used - can still easily break on YMD vs DMY order etc.
                            if key == "DateTime":
                                try:
                                    test_time = datetime.strptime(data[key], "%y-%b-%d %H:%M:%S")
                                    dataset_dict[key+"_min"] = (min(datetime.fromisoformat(dataset_dict[key]), test_time)).isoformat(timespec='minutes')
                                    dataset_dict[key+"_max"] = (max(datetime.fromisoformat(dataset_dict[key]), test_time)).isoformat(timespec='minutes')
                                except:
                                    ValueError
                                    #wrong time format, try again
                                try:
                                    test_time = datetime.strptime(data[key], "%Y-%b-%d %H:%M:%S")
                                    dataset_dict[key+"_min"] = (min(datetime.fromisoformat(dataset_dict[key]), test_time)).isoformat(timespec='minutes')
                                    dataset_dict[key+"_max"] = (max(datetime.fromisoformat(dataset_dict[key]), test_time)).isoformat(timespec='minutes')
                                except:
                                    ValueError
                                try:
                                    test_time = datetime.strptime(data[key], "%d-%b-%Y %H:%M:%S")
                                    dataset_dict[key+"_min"] = (min(datetime.fromisoformat(dataset_dict[key]), test_time)).isoformat(timespec='minutes')
                                    dataset_dict[key+"_max"] = (max(datetime.fromisoformat(dataset_dict[key]), test_time)).isoformat(timespec='minutes')
                                except:
                                    ValueError

                            elif key+"_min" in dataset_dict:
                                dataset_dict[key+"_min"] = min(float(dataset_dict[key+"_min"]), float(data[key]))
                            elif key+"_max" in dataset_dict:
                                dataset_dict[key+"_max"] = max(float(dataset_dict[key+"_max"]), float(data[key]))
                            elif "_min" in key:
                                dataset_dict[key] = min(float(dataset_dict[key]), float(data[key]))
                            elif "_max" in key:
                                dataset_dict[key] = max(float(dataset_dict[key]), float(data[key]))
                            else:
                                try:
                                    dataset_dict[key+"_min"] = min(float(dataset_dict[key]), float(data[key]))
                                    dataset_dict[key+"_max"] = max(float(dataset_dict[key]), float(data[key]))
                                except:
                                    TypeError
                           
    checklist = []    
    for k in dataset_dict:
        if k+"_max" in dataset_dict or k+"_min" in dataset_dict:
            checklist.append(k)
    for k in checklist:
        del dataset_dict[k]       
    return dataset_dict


if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python Metadata_extractor.py <directory>")
    # Check if the provided directory exists
    if not os.path.exists(sys.argv[1]):
        print(f"Error: Directory '{sys.argv[1]}' does not exist.")
        exit
    # Check if the provided directory is a directory
    if not os.path.isdir(sys.argv[1]):
        print(f"Error: '{sys.argv[1]}' is not a directory.")
        exit
    else:
        directory = sys.argv[1]
        target=read_in_files(directory)
        individual_images = []
        # Proceed depending on mdoc or xml:
        if target[0] != []:
            individual_images = target[0]
            results=check_all_files(individual_images)
            results["Number_of_Movies"] = target[2]
        elif target[1] != []:
            individual_images = target[1]
            results=check_all_files(individual_images)
            results["Number_of_Movies"] = target[2]
        else:
            print("Import failed")
            exit
        #generate Metadata filename based on directory name, add timestamp
        name= os.path.abspath(sys.argv[1])
        name_out = name.split("/")[-1] + "_" + (datetime.now()).isoformat(timespec='seconds') + '.json'
        print("Dataset level metadata file created as " + name_out )
        json.dump(results, codecs.open(name_out, 'w', encoding='utf-8'), 
            separators=(',', ':'), 
            sort_keys=True, 
            indent=4) 