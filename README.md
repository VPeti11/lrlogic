# LRLogic Renderer

**LRLogic Renderer** is a cli tool written in Go that reads custom `.lrlogic` files describing vector-based drawings and text annotations, then generates SVG vector images with optional JPG conversion.

It’s designed for "easy" scripting and batch processing, allowing detailed control over resolution, colors, curved line styling, margins, and fonts.

Dont ask why this exist. Its simple logic to me but i cant imagine this being used by anyone else than me

---

## Features

- Reads `.lrlogic` files in a "simple" text format
- Supports drawing colored curved lines and automatically fills closed polygons
- Configurable canvas resolution, margins, font sizes, and curve curvature
- Supports adding top and bottom text annotations with dividing lines
- Produces SVG vector-based output
- Optional JPG conversion
- Command line options for scripting and batch processing

---

## Installation

Make sure you have Go installed

Clone the repository or copy `main.go`, then build:

```
go build -o lrlogic main.go
```

## File format
LRLogic uses .lrlogic files. You can read more [here](LRLOGICfile.md)

For SVG conversion LRLogic uses W3 SVG2000 that you can read about [here](https://www.w3.org/2000/svg)


## Usage
### Run the program with a .lrlogic file:


Copy or Edit an lrlogic file then run the program
```
./lrlogic --file filename.lrlogic
```
This generates:

filename.svg — the vector image output

filename.jpg — a JPG version (requires rsvg-convert)

### Command-line Flags
    Flag	    Description	                    
    --file	    Path to .lrlogic input file	(required)
    --nojpg	    Skip generating JPG output	
    --nosvg	    Delete the SVG after JPG generation	
    --verbose   Verbose mode                            

### Example
Create a .lrlogic file (Version 1 syntax) e.g. square.lrlogic:

    LRLOGIC FILE FORMAT V1
    LRTXT.Top 'Simple Red Square'
    100,100,300,100..255,0,0
    300,100,300,300..255,0,0
    300,300,100,300..255,0,0
    100,300,100,100..255,0,0
    LREXIT
Run: 
```
./lrlogic --file square.lrlogic
```
Output files generated:

square.svg

square.jpg

## SVG2LR Helper 

See [SVG2LR README](svg2lrlogic/README.md) for more details.

This Go program converts basic shapes and path elements from an SVG file into .lrlogic format. It supports rectangles, circles, lines, polygons, simple paths (M/Q), text placement, and transform handling (translate, scale).

Vertical Y-coordinates are flipped to match .lrlogic top-down orientation,

Curved path elements are simplified to straight lines using Q command segments only,

Fills are toggled with LRFILL ON/OFF when appropriate


| Flag        | Type   | Description                       |
| ----------- | ------ | --------------------------------- |
| `--file`    | string | Path to the SVG file (required)   |
| `--verbose` | flag   | Enable detailed output to console |

| SVG Element | Conversion             |
| ----------- | ---------------------- |
| `<rect>`    | `LRSQUARE`             |
| `<circle>`  | `LRCIRCLE`             |
| `<line>`    | raw line format        |
| `<path>`    | simplified `M` and `Q` |
| `<polygon>` | line segments          |
| `<text>`    | `LRTXT.Top` / `Bottom` |

Transforms Supported:

translate(x, y),

scale(x[, y]),

These are applied to path elements using transform="...".

### Color Parsing
Accepts stroke="rgb(R,G,B)" or fill="rgb(R,G,B)"

If both are present, stroke takes priority

If fill="white" matches canvas size, it's treated as a background and skipped



## Usage of helper Python scripts
### randomgen.py
Output files are saved into a folder named randomgen. If --cleanup is enabled, only .jpg files remain.


This script generates .lrlogic files with randomized shapes and optionally renders them into .jpg images using an external lrlogic binary. It supports parallel generation for faster performance.

| Argument      | Type | Default   | Description                                     |
| ------------- | ---- | --------- | ----------------------------------------------- |
| `--count`     | int  | `100`     | Number of `.lrlogic` files to generate          |
| `--maxshapes` | int  | `100`     | Max number of shapes per file                   |
| `--width`     | int  | `None`    | Width of the canvas (prompted if not provided)  |
| `--height`    | int  | `None`    | Height of the canvas (prompted if not provided) |
| `--render`    | flag | `True`    | Render files to `.jpg` using `./lrlogic`        |
| `--cleanup`   | flag | `False`   | Delete `.lrlogic` files after rendering         |
| `--ratio`     | str  | `"5,2,3"` | Shape ratio for line,circle,square generation   |

#### Example:
    python3 randomgen.py --count 50 --maxshapes 30 --width 800 --height 600 --render --cleanup --ratio 4,3,3
Generates 50 .lrlogic files,

Each with up to 30 shapes,

On an 800x600 canvas,

Renders them to .jpg,

Deletes the .lrlogic files after rendering,

Uses a shape distribution of 4 line : 3 circle : 3 square

### transformsvg2lr.py
This Python script converts vector path data from an SVG file into .lrlogic format. It flattens Bezier curves into straight line segments and supports color customization and SVG transformations. 

Want to change curve accuracy and file size? Increase or decrease the segments=20 value in the code.

| Argument      | Type             | Default    | Description                                       |
| ------------- | ---------------- | ---------- | ------------------------------------------------- |
| `svg_file`    | str| *required* | Path to the SVG file to convert                   |
| `--verbose`   | flag             | `False`    | Show debug information during conversion          |
| `--rgb R G B` | 3 ints           | `0 0 0`    | Override default black stroke with RGB color. I usually use neon yellow or dark blue (255, 255, 51 and 0, 0, 139) |

If no arguments are provided, the script will prompt for:

SVG file path,

Whether to enable verbose output,

Optional RGB override



#### Transform Support
This script supports the following SVG transform types:

translate(x, y)

scale(x, y)

rotate(angle, [cx, cy])

These are applied correctly to path coordinates during conversion.




## Scripts
LRLogic comes with scripts to help with using the software. These scripts work on linux and windows. Note that the linux scipts are ported to windows and not the other way around. 
### Script functions:
    Cleanup: cleans image files from the root of the repo
    Full test: Compiler LRLogic and renders preincluded testfiles
    Makerandom: Runs a Python helper script to generate random .lrlogic files and renders them
    Render all: Render all .lrlogic files a the project root directory
    Timed render: Render a specified file with a time elapsed message at the end
    Compile LRLogic: Compiles the LRLogic Go file in the repos root directory
    Compile SVG2LR: Compiles the LRLogic SVG2LR helper program (Go versiom)

## Compatibility
This software is designed to run on both linux and windows but the software is meant to be linux first. Thanks to the Go compiler im able to ship windows executables with the same code. As for the scripts like i mentioned it before they are ported from linux to windows. If you want maximum compatibility use linux or WSL

## Dependencies
LRLogic is a Go application so dependencies and the Go Runtime is bundled into the executable file. Note that LRLogic calls the external program  `rsvg-convert` to render files in JPG. This is can be skipped by running the program with the `--nojpg` option. And the SVG2LR (Python version) requires the `svg.path` library. This can be installed using PIP (Or in the case of Arch Linux it can be installed using the AUR). The "scripts" folder includes scripts for both windows and linux to install dependencies with the option to install the go compiler. The linux version works with apt,dnf and pacman. The script can automatically determine the package manager and install required packages. The windows version uses Chocolatey to install dependencies (including Python). As a bonus if Chocolatey is not installed it will install it for you.

## License
This software is licensed under the GNU General Public License (V3). This means:

GPLv3 is a free software license that ensures users' freedom to run, study, share, and modify the software.
Any redistributed versions must also be licensed under GPLv3, preserving these freedoms.
You must make the source code available when distributing the program or any derivatives.
You may not impose further restrictions that conflict with the license.
The license also includes protections against tivoization, software patents, and anti-circumvention laws.

| Feature                    | Description                                                                |
| -------------------------- | -------------------------------------------------------------------------- |
| Freedom to Use             | You can run the program for any purpose                                    |
| Freedom to Modify          | You can study the source code and make changes                             |
| Freedom to Share           | You can distribute original or modified versions under the same license    |
| Source Code Disclosure     | When distributing, you must provide or offer access to the source code     |
| No Tivoization / DRM Locks | You cannot use hardware restrictions to block modified software            |
| Patent Protection          | Users are protected from patent claims that would restrict software rights |

As an extra the `.lrlogic` files included in this project are also licensed under the GPL. But the art you create (any format) can be licensed under any license or you can choose to copyright it even. But Creative Commons licenses are preffered.

You can read the license [here](LICENSE.md)

The documentation and README files are licensed under the GFDL that you can read [here](fdl.md)

# Contributing
If you want to make the project better check out the docs [here](CONTRIBUTING.md)
