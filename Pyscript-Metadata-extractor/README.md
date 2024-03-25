# PyScript benchmark

Prototype of the Metadata extractor running in pyscript, derivative of https://github.com/SwissOpenEM/pyscript-fs-benchmark

## Target data

Any typical raw data folder generated either by SerialEM or EPU.
For my testing did discriminate files properly - so the presence of several GBs of image files does not negatively impact performance.

## Requirements

Needs Pyscript  https://pyscript.net/
On my Mac required me to change the requirements.txt from 
playwright==1.33.0 to playwright>=1.33.0 because of some issues with older greenlet versions.

## Running

Tested using a python venv that setup pyscript, start the server in your clones directory using

```sh
$ python -m http.server 
```
 and connect to localhost:8000 in chrome.
