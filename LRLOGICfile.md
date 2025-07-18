# LRLOGIC File Format Specification


## LRLOGIC FILE FORMAT V2

The first line must be:

```
LRFILE VERSION 2
```

---

### Commands and Syntax

* `LRFILL ON` / `LRFILL OFF`
  Enables or disables shape filling for polygons, circles, and squares. Polygons are filled only if `LRFILL ON` (default OFF in V2).

* `LRCIRCLE x,y,radius..r,g,b`
  Draws a circle centered at `(x, y)` with the specified `radius`. Color is given by RGB values. Coordinates use the bottom-left origin. Filled only if `LRFILL ON`.

* `LRSQUARE x,y,size..r,g,b`
  Draws a square with bottom-left corner `(x, y)` and side length `size`. Color as above. Filled only if `LRFILL ON`.

* Behavior changes:

  * `LRFILL` controls fill behavior (default OFF).
  * Polygons respect the `LRFILL` flag (unlike V1 where polygons are always filled).
  * New commands: `LRCIRCLE` and `LRSQUARE`.
  * Coordinates use bottom-left origin.

* Backward compatibility:

  * Files with header `LRLOGIC FILE FORMAT V1` are parsed with V1 rules.

---

## LRLOGIC FILE FORMAT V1

The first line must be:

```
LRLOGIC FILE FORMAT V1
```

---

### Commands and Syntax

1. **Canvas Resolution**

```
LRRESDEFINEX width
LRRESDEFINEY height
```

Defaults to `640x480` if not set.

2. **Margins**

```
LRMARGIN top bottom
```

Defaults: `20 20`

3. **Font Size**

```
LRFONTSIZE size
```

Default: `16`

4. **Curve Strength**

```
LRCURVE strength
```

Default: `5`

5. **Text Commands**

```
LRTXT.Top 'text'
LRTXT.Bottom 'text'
```

* Text must be enclosed in single quotes.
* Horizontal lines are automatically drawn below (top) or above (bottom) text.

6. **Drawing Lines**

```
x1,y1,x2,y2..R,G,B
```

* Default color is black.
* All coordinates must be integers.
* Origin (0,0) is at top-left corner.
* +X goes right, +Y goes down.

7. **File End**

```
LREXIT
```

Must be the last line.

---

### Shape Filling

* Multiple lines with the same RGB color forming a closed shape (start = end) are automatically filled.
* Works for triangles, squares, pentagons, and any closed polygon.

---

### Rules & Limitations

* Lines must be defined before `LREXIT`.
* Only single-quoted text supported.
* Coordinates and RGB values are integers.
* Blank lines or comments not supported.
* Commands are case-sensitive.

---

## File Naming

* Outputs:

  * `filename.svg` (vector image)
  * `filename.jpg` (if JPG enabled)
* Filename based on input `.lrlogic` file.

---

## Tips

* Use precise coordinates to close polygons cleanly.
* Reuse RGB values for grouped shapes.
* Use higher resolution for detailed drawings.
* Use mild `LRCURVE` values (3–5) for slight bending.
