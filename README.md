# Metadata extractors
simple parser for the two common life-science EM metadata output formats, written in go - python version outdated for now

## Requirements
Needs lxml for .xml parsing in python

## Comments
Runs on a directory containing raw files and their instrument written additional information files (.mdoc and .xml respectively), generates a dataset level .json file. In case of usage with EPU pointing to the top level directory is enough, it will search for the data folders and extract the info from there. Using --z you can also obtain a zip file of the xml files associated with your data collection.

## Schema-Links 
Starting to link outputs to the OSCEM-Schemas https://github.com/osc-em/OSCEM_Schemas/tree/main/Instrument

## TODO
- find a way of assigning illumination modes for mdocs
- more inference based values if applicable and needed 
- folder crawling for serialEM; problem: non standardardized directory architecture/naming