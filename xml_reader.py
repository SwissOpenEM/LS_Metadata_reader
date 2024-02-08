### EXTRACT METADATA FROM EPU GENERATED XML FILES in (cryo)-electron microscopy
import os, sys
import json, codecs
from lxml import etree

if (((len(sys.argv) == 2) & os.path.isfile(sys.argv[1]) == True ) & ("xml" in (sys.argv[1])) ):
    file_in= sys.argv[1]
else:
    print("Please give a valid path to an EPU generated XML file")
    exit()

# Parse the XML data using lxml
test = etree.parse(file_in)
root = test.getroot()

key_value_dict = {}

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

# Search and process the XML for Tag-Value pairs
tag_text_dict = search_and_process_tree(root)
# Search and process the XML for Key-Value pairs
get_key_value(root, key_value_dict=key_value_dict)
### Merge the results:
merge_results = {**key_value_dict, **tag_text_dict}

# Write out 
json.dump(merge_results, codecs.open(file_in.replace("xml", "json"), 'w', encoding='utf-8'), 
          separators=(',', ':'), 
          sort_keys=True, 
          indent=4) 
