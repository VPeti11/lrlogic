# LRLOGIC File Format Specification (Version 1)

The `.lrlogic` file format defines vector-based instructions for generating SVG drawings using the LRLogic Renderer. It supports text, curves, resolution settings, and automatic shape filling based on color.

---

## File Header

The **first line** of every `.lrlogic` file must be:

LRLOGIC FILE FORMAT V1

yaml
Copy
Edit

This identifies the format and must be present exactly as shown.

---

## Commands and Syntax

Each subsequent line can be one of the following instructions:

### 1. Canvas Resolution

#### Set the image dimensions (in pixels):

LRRESDEFINEX 'width'

LRRESDEFINEY 'height'

#### Example:

LRRESDEFINEX 1024

LRRESDEFINEY 768

#### Defaults: `640x480` if not set.

---

### 2. Margins

Adjust margins (in pixels) above the top text and below the bottom text:

LRMARGIN 'top' 'bottom'

#### Example:

LRMARGIN 40 60

Default: `20 20`

---

### 3. Font Size

#### Set the font size (in pixels) for text at top or bottom:

LRFONTSIZE <size>

#### Example:

LRFONTSIZE 24


#### Default: `16`

---

### 4. Curve Strength

Adjust the curvature of drawn lines. Higher = more curved:

LRCURVE 'strength'


Example:

LRCURVE 4


#### Default: `5`

---

### 5. Text Commands

#### Add text with optional dividers:

LRTXT.Top 'Your top text here'

LRTXT.Bottom 'Your bottom text here'


- Text must be enclosed in **single quotes**
- A horizontal line is automatically drawn below (top) or above (bottom) the text

#### Example:

LRTXT.Top 'Hello World'

LRTXT.Bottom 'Made with LRLogic'


---

### 6. Drawing Lines

Each line is specified using 4 coordinates:

x1,y1,x2,y2

Optionally, you can add color after the line using a double-dot (`..`) followed by an RGB triplet:

x1,y1,x2,y2..R,G,B

Example:
100,100,300,100..255,0,0

This draws a red curved line from (100,100) to (300,100).

#### Notes:
- Default color is black if not specified
- All coordinates must be integers
- Origin (0,0) is at the **top-right corner**
- +X goes right, +Y goes down

---

### 7. File End

Signal the end of input with:

    LREXIT

This must be the last line.

---

##  Shape Filling

If **multiple lines with the same RGB color** form a **closed shape**, the area is automatically filled with that color.

Requirements:
- Must form a closed loop (start = end)
- All edges must have the exact same RGB color

Works for:
- Triangles
- Squares
- Pentagons
- Any closed polygon

---


## Rules & Limitations

- Lines must be defined before `LREXIT`
- Only single-quoted text is supported
- Coordinates must be integers
- RGB values must be integers (0–255)
- Blank lines or comments are **not** supported
- All commands are **case-sensitive**

---

##  File Naming

- The program outputs:
  - `filename.svg` (vector image)
  - `filename.jpg` (if JPG is enabled)
- File name is based on the `.lrlogic` file name

---

## Tips

- Use precise coordinates to ensure polygons close cleanly
- Reuse RGB values for grouped shapes
- Use higher resolutions for more detailed drawings
- Use mild `LRCURVE` (3–5) for slight bending

---
