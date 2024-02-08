### MDOC file parser for cryoEM data collections using SerialEM (or Tomo5)
import sys, os, re  
import json, codecs

if (((len(sys.argv) == 2) & os.path.isfile(sys.argv[1]) == True ) & ("mdoc" in (sys.argv[1])) ):
    file_in= sys.argv[1]
else:
    print("Please give a valid path to a mdoc file")
    exit()
    
def extract_key_value_pairs(fil):
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
                    key_value_pairs["Tilt_axis_angle"] = float( test_axis.strip())
                except:
                    ValueError
                try:
                    test_axis = re.split(r'[A-Za-z]', (line.split("=")[2]))[0]
                    key_value_pairs["Tilt_axis_angle"] = float(test_axis.strip())
                except:
                    ValueError
            # count number of tilts in tiltseries
            if "ZValue" in line:
                test_z = (line.split("=")[1]).split("]")[0]
                if int (test_z.strip() ) > numb:
                    numb = int( test_z.strip())
                key_value_pairs["NumberOfTilts"] = numb
    #remove useless last recorded value for items with min/max
    checklist = []    
    for k in key_value_pairs:
        if k+"_max" in key_value_pairs or k+"_min" in key_value_pairs:
            checklist.append(k)
    for k in checklist:
        del key_value_pairs[k]        
    return key_value_pairs


 # Write out
json.dump(extract_key_value_pairs(file_in), codecs.open(file_in.replace("mdoc", "json"), 'w', encoding='utf-8'), 
          separators=(',', ':'), 
          sort_keys=True, 
          indent=4) 
