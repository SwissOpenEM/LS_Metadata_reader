# Metadata extractors
simple parser for the two common life-science EM metadata output formats, written in both go and python

## Requirements
Needs lxml for .xml parsing in python

## Comments
Runs on a directory containing raw files and their instrument written additional information files (.mdoc and .xml respectively), generates a dataset level .json file in that directory.

## Schema-Links 
Starting to link outputs to the OSCEM-Schemas https://github.com/osc-em/OSCEM_Schemas/tree/main/Instrument

## TODO

- Folder crawling
- might add another function for .json written by Athena for Tomo5/ merge of xml and mdoc
- decision on language
- extension of infered values based on schema needs




