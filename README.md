# LS Metadata extractor
simple parser for the two common life-science EM metadata output formats (.xml from EPU and .mdoc from TOMO5 and SerialEM respectively), written in go

## Usage
Chose the appropriate binary from the [Releases](https://github.com/SwissOpenEM/LS_Metadata_reader/releases), then:
LS_Metadata_reader target_directory

For testing, try the associated [tutorial](https://github.com/SwissOpenEM/LS_Metadata_reader/tree/main/tutorial) folder; an example of how the output should look like is provided in the same folder (tutorial_correct.json).

## Comments
Runs on a directory containing raw files and their instrument written additional information files (.mdoc and .xml respectively), generates a dataset level .json file. In case of usage with EPU pointing to the top level directory is enough, it will search for the data folders and extract the info from there. Using --z you can also obtain a zip file of the xml files associated with your data collection. If you want all the metadata (dataset level, not all individual entries) written out by a given software use the --f flag, otherwise the output will be OSCEM conform. 

## SerialEM
SerialEM properties examples are to be added to the existing properties files of your SerialEM installation (update values to reflect your instrument parameters). The two scripts are to be run after each image collection (the lowest tick mark on the SerialEM automization script selection) with the respective name indicating when to use which of the two. Otherwise SerialEM ouput will lack a few required fields for the schema. 
**!!! Requires SerialEM 4.2.0 or newer !!!**

## For running with EPU directly
Benefits from setting the config (set using LS_Metadata_reader --c), or handing over the three instrument values directly with flags: <br>
--cs <br>
--gain_flip_rotate <br>
--epu <br>
The --cs (for the CS value of the instrument) and --gain_flip_rotate (for the orientation of the gain_reference relative to actual data) are unfortunately never provided in the metadata, and are both important for processing. It is therefore highly beneficial to set these two.
As for --epu, EPU writes its metadata files in a different directory than its actual data (TOMO5 also keeps some additional info that is processed by the LS_Metadata_reader there). It generates another set of folders, usually on the microscope controlling computer, that mirror its OffloadData folders in directory structure. Within them it stores some related information, among which are also the metadata xml files. If --epu is defined as a flag or in the config, the LS_Metadata_reader will directly grab those when the user points it at a OffloadData directory. <br>
NOTE: This requires you to mount the microscope computer directory for EPU on the machine you are running LS_Metadata_reader on as those are most likely NOT the same. The extractor will work regardless if pointed to the xmls/mdocs directly, this is just for convenience.

## Schema-Links 
Output is compatible to OSCEM schemas https://github.com/osc-em/OSCEM_Schemas/

Specific schema used to generate standard schema conform output (works for SPA and Tomography): https://github.com/osc-em/OSCEM_Schemas/blob/linkml_yaml/src/oscem_schemas/schema/oscem_schemas_tomo.yaml 
with LinkML gen-golang
