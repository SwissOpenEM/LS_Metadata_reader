# Metadata extractors
simple parser for the two common life-science EM metadata output formats, written in go

## Requirements
Benefits from setting an environmental variable "MPCPATH" to define the path of data acquistisions mirrored on the microscope PC in EPU. Will work regardless if pointed to the xmls/mdocs otherwise.

## Usage
Chose the appropriate binary from the [Releases](https://github.com/SwissOpenEM/LS_Metadata_reader/releases), then:
LS_reader_Version directory

## Comments
Runs on a directory containing raw files and their instrument written additional information files (.mdoc and .xml respectively), generates a dataset level .json file. In case of usage with EPU pointing to the top level directory is enough, it will search for the data folders and extract the info from there. Using --z you can also obtain a zip file of the xml files associated with your data collection. If the full metadata written out by a given software is required use the --f flag. 

## SerialEM
SerialEM properties examples are to be added to the existing properties files of your SerialEM installation (update values to reflect your instrument parameters). The two scripts are to be run after each image collection (the lowest tick mark on the SerialEM automization script selection) with the respective name indicating when to use which of the two. Otherwise SerialEM ouput will lack a few required fields for the schema. !!! Requires SerialEM 4.2.0 or newer !!!

## Schema-Links 
Output is conform to this schema https://github.com/osc-em/OSCEM_Schemas/tree/main/Instrument

## TODO
- find a way of assigning illumination modes for mdocs