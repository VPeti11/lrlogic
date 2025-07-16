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
    Flag	                Description	                    Default
    --file	Path to .lrlogic input file	(required)
    --nojpg	Skip generating JPG output	false
    --nosvg	Delete the SVG after JPG generation	false


### Example
Create a .lrlogic file e.g. square.lrlogic:

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

### Batch Processing
Use the included  script to process all .lrlogic files:

