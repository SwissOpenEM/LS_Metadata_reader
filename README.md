# LS Metadata extractor
simple parser for the two common life-science EM metadata output formats, written in go

## Usage
Chose the appropriate binary from the [Releases](https://github.com/SwissOpenEM/LS_Metadata_reader/releases), then:
LS_reader_Version target_directory

For testing, try the associated tutorial folder; an example of how the output should look like is provided in the same folder (tutorial_correct.json). For first time use disregard the warnings about config/flags those are for use directly with EPU or the OpenEM Ingestor.

## Comments
Runs on a directory containing raw files and their instrument written additional information files (.mdoc and .xml respectively), generates a dataset level .json file. In case of usage with EPU pointing to the top level directory is enough, it will search for the data folders and extract the info from there. Using --z you can also obtain a zip file of the xml files associated with your data collection. If you want all the metadata (dataset level, not all individual entries) written out by a given software use the --f flag, otherwise the output will be OSCEM conform. 

## SerialEM
SerialEM properties examples are to be added to the existing properties files of your SerialEM installation (update values to reflect your instrument parameters). The two scripts are to be run after each image collection (the lowest tick mark on the SerialEM automization script selection) with the respective name indicating when to use which of the two. Otherwise SerialEM ouput will lack a few required fields for the schema. 
**!!! Requires SerialEM 4.2.0 or newer !!!**

## For running with EPU directly
Benefits from setting the "MPCPATH" variable in the .config file (set using LS_reader_version --c) to define the path of data acquistisions mirrored on the microscope PC in EPU. Will work regardless if pointed to the xmls/mdocs otherwise.

## Schema-Links 
Output is compatible to OSCEM schemas https://github.com/osc-em/OSCEM_Schemas/

Specific schema used to generate standard schema conform output (works for SPA and Tomography): https://github.com/osc-em/OSCEM_Schemas/blob/linkml_yaml/src/oscem_schemas/schema/oscem_schemas_tomo.yaml 
with LinkML gen-golang
