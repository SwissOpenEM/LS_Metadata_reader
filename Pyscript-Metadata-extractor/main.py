"""Pyscript entrypoint
"""
from pathlib import Path
from pyscript import document, window, when, display
from Metadata_extractor_pyscript import run
# from generate_testdata import generate
import pyodide_js
import os
from pathlib import Path


@when("click", "#fileBtn")
async def click_handler(event):
    """
    Event handlers get an event object representing the activity that raised
    them.
    """

    dirHandle = await window.showDirectoryPicker();
    display(f"Reading {dirHandle.name}")

    # print(f"dirHandle: {dir(dirHandle)}")
    if await dirHandle.queryPermission(mode="readwrite") != "granted":
        if await dirHandle.requestPermission(mode="readwrite") != "granted":
            raise Error("Unable to read and write directory")
  
  
    nativefs = await pyodide_js.mountNativeFS("/upload", dirHandle)
    root = Path("/upload")
    
    msg = run(root)

    display(msg)

# def generate(event):
#     output_div = document.querySelector("#output")

#     root = Path("data")
#     root.mkdir(exist_ok=True)
#     with Path(root, "1.txt").open('wb') as file:
#         file.write(b"hello world")

#     generate(root, 2, 10)

#     output = f'<p>Root: {root.absolute()}</p>'

#     output += f'<p>Files: {", ".join(str(f) for f in root.iterdir())}</p>'
#     # return "foo"
#     # filenum = 0
#     # files = sorted(f for f in root.iterdir() if f.is_file())
#     # if filenum > 0:
#     #     files = files[:filenum]

#     # t = timeit.Timer(lambda: read_files(files))
#     # elapsed, size = t.timeit(number=1)
#     # return f"Read {size} bytes from {len(files)} files in {elapsed}s = {metricunit(size/elapsed)}B/s"

#     output_div.innerHTML = output