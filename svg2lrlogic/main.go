package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	fileFlag   *string
	verbose    *bool
	pythonFlag *bool
	rgbStr     *string
)

type Line struct {
	X1, Y1, X2, Y2 int
	R, G, B        int
}

func main() {
	fileFlag = flag.String("file", "", "Path to the SVG file (required)")
	verbose = flag.Bool("verbose", false, "Enable verbose output")
	pythonFlag = flag.Bool("python", false, "Run the embedded Python version instead of Go")
	rgbStr = flag.String("rgb", "", "Override RGB color (e.g. \"255 128 64\")")

	flag.Parse()

	if *pythonFlag {
		runPythonVersion()
	}

	if *fileFlag == "" {
		fmt.Println("Usage: svg2lrlogic --file input.svg [--verbose] [--python] [--rgb R G B]")
		os.Exit(1)
	}

	var rgb [3]int
	if *rgbStr != "" {
		parts := strings.Fields(*rgbStr)
		if len(parts) != 3 {
			fmt.Println("Error: --rgb requires exactly 3 integer values")
			os.Exit(1)
		}
		for i, p := range parts {
			val, err := strconv.Atoi(p)
			if err != nil || val < 0 || val > 255 {
				fmt.Printf("Error: RGB values must be integers 0-255, got %q\n", p)
				os.Exit(1)
			}
			rgb[i] = val
		}
	} else {
		rgb = [3]int{0, 0, 0}
	}

	data, err := ioutil.ReadFile(*fileFlag)
	if err != nil {
		log.Fatalf("Failed to read SVG file: %v", err)
	}

	decoder := xml.NewDecoder(strings.NewReader(string(data)))
	var output []string
	output = append(output, "LRFILE VERSION 2")

	width, height := 640, 480
	fillState := ""
	lastFill := ""
	output = append(output, "LRMARGIN 20 20", "LRFONTSIZE 16", "LRCURVE 5")

	for {
		tok, err := decoder.Token()
		if err != nil {
			break
		}

		switch elem := tok.(type) {
		case xml.StartElement:
			switch elem.Name.Local {
			case "svg":
				for _, attr := range elem.Attr {
					if attr.Name.Local == "width" {
						width, _ = strconv.Atoi(attr.Value)
						output = append(output, fmt.Sprintf("LRRESDEFINEX %d", width))
					}
					if attr.Name.Local == "height" {
						height, _ = strconv.Atoi(attr.Value)
						output = append(output, fmt.Sprintf("LRRESDEFINEY %d", height))
					}
				}
			case "rect":
				x, y, w, h := 0, 0, 0, 0
				fill := "none"
				rCol, gCol, bCol := 0, 0, 0
				for _, attr := range elem.Attr {
					switch attr.Name.Local {
					case "x":
						x, _ = strconv.Atoi(attr.Value)
					case "y":
						y, _ = strconv.Atoi(attr.Value)
					case "width":
						w, _ = strconv.Atoi(attr.Value)
					case "height":
						h, _ = strconv.Atoi(attr.Value)
					case "fill":
						fill = attr.Value
					case "stroke":
						rCol, gCol, bCol = parseRGB(attr.Value)
					}
				}
				if w == width && h == height && strings.TrimSpace(fill) == "white" {
					continue // skip background
				}
				if w <= 0 || h <= 0 {
					continue
				}
				y = height - y - h
				fillState = "OFF"
				if fill != "none" {
					fillState = "ON"
				}
				if fillState != lastFill {
					output = append(output, "LRFILL "+fillState)
					lastFill = fillState
				}
				size := w
				output = append(output, fmt.Sprintf("LRSQUARE %d,%d,%d..%d,%d,%d", x, y, size, rCol, gCol, bCol))
			case "circle":
				x, y, r := 0, 0, 0
				fill := "none"
				rCol, gCol, bCol := 0, 0, 0
				for _, attr := range elem.Attr {
					switch attr.Name.Local {
					case "cx":
						x, _ = strconv.Atoi(attr.Value)
					case "cy":
						y, _ = strconv.Atoi(attr.Value)
					case "r":
						r, _ = strconv.Atoi(attr.Value)
					case "fill":
						fill = attr.Value
					case "stroke":
						rCol, gCol, bCol = parseRGB(attr.Value)
					}
				}
				if r <= 0 {
					continue
				}
				y = height - y
				fillState = "OFF"
				if fill != "none" {
					fillState = "ON"
				}
				if fillState != lastFill {
					output = append(output, "LRFILL "+fillState)
					lastFill = fillState
				}
				output = append(output, fmt.Sprintf("LRCIRCLE %d,%d,%d..%d,%d,%d", x, y, r, rCol, gCol, bCol))
			case "line":
				x1, y1, x2, y2 := 0, 0, 0, 0
				rCol, gCol, bCol := 0, 0, 0
				for _, attr := range elem.Attr {
					switch attr.Name.Local {
					case "x1":
						x1, _ = strconv.Atoi(attr.Value)
					case "y1":
						y1, _ = strconv.Atoi(attr.Value)
					case "x2":
						x2, _ = strconv.Atoi(attr.Value)
					case "y2":
						y2, _ = strconv.Atoi(attr.Value)
					case "stroke":
						rCol, gCol, bCol = parseRGB(attr.Value)
					}
				}
				y1 = height - y1
				y2 = height - y2
				if lastFill != "OFF" {
					output = append(output, "LRFILL OFF")
					lastFill = "OFF"
				}
				output = append(output, fmt.Sprintf("%d,%d,%d,%d..%d,%d,%d", x1, y1, x2, y2, rCol, gCol, bCol))
			case "path":
				var d string
				rCol, gCol, bCol := 0, 0, 0
				for _, attr := range elem.Attr {
					if attr.Name.Local == "d" {
						d = attr.Value
					}
					if attr.Name.Local == "stroke" {
						rCol, gCol, bCol = parseRGB(attr.Value)
					}
				}
				tokens := strings.Fields(d)
				var x1, y1 int
				i := 0
				for i < len(tokens) {
					switch tokens[i] {
					case "M":
						if i+2 < len(tokens) {
							x1, _ = strconv.Atoi(tokens[i+1])
							y1, _ = strconv.Atoi(tokens[i+2])
							i += 3
						} else {
							i++
						}
					case "Q":
						if i+4 < len(tokens) {
							x2, _ := strconv.Atoi(tokens[i+3])
							y2, _ := strconv.Atoi(tokens[i+4])
							y1i := height - y1
							y2i := height - y2
							if lastFill != "OFF" {
								output = append(output, "LRFILL OFF")
								lastFill = "OFF"
							}
							output = append(output, fmt.Sprintf("%d,%d,%d,%d..%d,%d,%d", x1, y1i, x2, y2i, rCol, gCol, bCol))
							if *verbose {
								fmt.Printf("Parsed path segment: (%d,%d) to (%d,%d) rgb(%d,%d,%d)\n", x1, y1i, x2, y2i, rCol, gCol, bCol)
							}
							x1 = x2
							y1 = y2
							i += 5
						} else {
							i++
						}
					default:
						i++
					}
				}
			case "polygon":
				var points string
				rCol, gCol, bCol := 0, 0, 0
				fill := "none"
				for _, attr := range elem.Attr {
					if attr.Name.Local == "points" {
						points = attr.Value
					}
					if attr.Name.Local == "fill" {
						fill = attr.Value
						rCol, gCol, bCol = parseRGB(fill)
					}
				}
				fillState = "OFF"
				if fill != "none" {
					fillState = "ON"
				}
				if fillState != lastFill {
					output = append(output, "LRFILL "+fillState)
					lastFill = fillState
				}
				pts := strings.Fields(points)
				for i := 0; i < len(pts)-1; i++ {
					p1 := strings.Split(pts[i], ",")
					p2 := strings.Split(pts[i+1], ",")
					x1, _ := strconv.Atoi(p1[0])
					y1, _ := strconv.Atoi(p1[1])
					x2, _ := strconv.Atoi(p2[0])
					y2, _ := strconv.Atoi(p2[1])
					y1 = height - y1
					y2 = height - y2
					output = append(output, fmt.Sprintf("%d,%d,%d,%d..%d,%d,%d", x1, y1, x2, y2, rCol, gCol, bCol))
				}
			case "text":
				var y int
				var content string
				for _, attr := range elem.Attr {
					if attr.Name.Local == "y" {
						y, _ = strconv.Atoi(attr.Value)
					}
				}
				tok, _ := decoder.Token()
				if charData, ok := tok.(xml.CharData); ok {
					content = string(charData)
				}
				if y < height/2 {
					output = append(output, fmt.Sprintf("LRTXT.Top '%s'", content))
				} else {
					output = append(output, fmt.Sprintf("LRTXT.Bottom '%s'", content))
				}
				if *verbose {
					fmt.Printf("Parsed text: '%s' at y=%d\n", content, y)
				}
			}
		}
	}

	output = append(output, "LREXIT")

	outFile := strings.TrimSuffix(filepath.Base(*fileFlag), filepath.Ext(*fileFlag)) + ".lrlogic"
	err = ioutil.WriteFile(outFile, []byte(strings.Join(output, "\n")), 0644)
	if err != nil {
		log.Fatalf("Failed to write output file: %v", err)
	}

	fmt.Printf("Generated %s successfully\n", outFile)
}

func installSvgPathAUR() error {
	// Check if pacman exists (assume Arch if yes)
	if _, err := exec.LookPath("pacman"); err != nil {
		return fmt.Errorf("pacman not found, not Arch Linux")
	}

	// Check if base-devel group installed
	fmt.Println("Checking if base-devel group is installed...")
	checkCmd := exec.Command("pacman", "-Qi", "base-devel")
	if err := checkCmd.Run(); err != nil {
		fmt.Println("base-devel not found, installing...")
		installCmd := exec.Command("sudo", "pacman", "-S", "--needed", "--noconfirm", "base-devel")
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr
		if err := installCmd.Run(); err != nil {
			return fmt.Errorf("failed to install base-devel: %v", err)
		}
	} else {
		fmt.Println("base-devel group already installed.")
	}

	// Create temp dir for cloning
	tmpDir, err := os.MkdirTemp("", "aur-python-svgpath")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Clone AUR repo
	fmt.Println("Cloning python-svg.path AUR repo...")
	gitClone := exec.Command("git", "clone", "https://aur.archlinux.org/python-svg.path.git", tmpDir)
	gitClone.Stdout = os.Stdout
	gitClone.Stderr = os.Stderr
	if err := gitClone.Run(); err != nil {
		return fmt.Errorf("git clone failed: %v", err)
	}

	// Run makepkg -si inside cloned repo
	fmt.Println("Running makepkg -si inside AUR repo...")
	makepkg := exec.Command("makepkg", "-si", "--noconfirm")
	makepkg.Dir = tmpDir
	makepkg.Stdout = os.Stdout
	makepkg.Stderr = os.Stderr
	if err := makepkg.Run(); err != nil {
		return fmt.Errorf("makepkg failed: %v", err)
	}

	fmt.Println("python-svg.path installed successfully via AUR.")
	return nil
}

func parseRGB(rgb string) (int, int, int) {
	rgb = strings.TrimPrefix(rgb, "rgb(")
	rgb = strings.TrimSuffix(rgb, ")")
	parts := strings.Split(rgb, ",")
	if len(parts) != 3 {
		return 0, 0, 0
	}
	r, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
	g, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
	b, _ := strconv.Atoi(strings.TrimSpace(parts[2]))
	return r, g, b
}

func convertTransformedPathsToLines(paths []xml.StartElement, transform string, height int, verbose bool) []string {
	var output []string
	scaleX, scaleY := 1.0, 1.0
	translateX, translateY := 0.0, 0.0

	// Parse transform
	if strings.Contains(transform, "translate(") {
		start := strings.Index(transform, "translate(") + len("translate(")
		end := strings.Index(transform[start:], ")")
		if end != -1 {
			inner := transform[start : start+end]
			parts := strings.FieldsFunc(inner, func(r rune) bool {
				return r == ',' || r == ' '
			})
			if len(parts) >= 1 {
				translateX, _ = strconv.ParseFloat(parts[0], 64)
			}
			if len(parts) >= 2 {
				translateY, _ = strconv.ParseFloat(parts[1], 64)
			}
		}
	}
	if strings.Contains(transform, "scale(") {
		start := strings.Index(transform, "scale(") + len("scale(")
		end := strings.Index(transform[start:], ")")
		if end != -1 {
			inner := transform[start : start+end]
			parts := strings.FieldsFunc(inner, func(r rune) bool {
				return r == ',' || r == ' '
			})
			if len(parts) >= 1 {
				scaleX, _ = strconv.ParseFloat(parts[0], 64)
			}
			if len(parts) >= 2 {
				scaleY, _ = strconv.ParseFloat(parts[1], 64)
			} else {
				scaleY = scaleX
			}
		}
	}

	// Process each path
	for _, elem := range paths {
		var d string
		rCol, gCol, bCol := 0, 0, 0
		for _, attr := range elem.Attr {
			if attr.Name.Local == "d" {
				d = attr.Value
			}
			if attr.Name.Local == "stroke" {
				rCol, gCol, bCol = parseRGB(attr.Value)
			}
		}

		tokens := strings.Fields(d)
		var x1, y1 int
		i := 0
		for i < len(tokens) {
			switch tokens[i] {
			case "M":
				if i+2 < len(tokens) {
					x1f, _ := strconv.ParseFloat(tokens[i+1], 64)
					y1f, _ := strconv.ParseFloat(tokens[i+2], 64)
					x1 = int(x1f*scaleX + translateX)
					y1 = int(y1f*scaleY + translateY)
					i += 3
				} else {
					i++
				}
			case "Q":
				if i+4 < len(tokens) {
					x2f, _ := strconv.ParseFloat(tokens[i+3], 64)
					y2f, _ := strconv.ParseFloat(tokens[i+4], 64)
					x2 := int(x2f*scaleX + translateX)
					y2 := int(y2f*scaleY + translateY)
					y1i := height - y1
					y2i := height - y2
					output = append(output, fmt.Sprintf("%d,%d,%d,%d..%d,%d,%d", x1, y1i, x2, y2i, rCol, gCol, bCol))
					if verbose {
						fmt.Printf("Transformed path line: (%d,%d) to (%d,%d) rgb(%d,%d,%d)\n", x1, y1i, x2, y2i, rCol, gCol, bCol)
					}
					x1 = x2
					y1 = y2
					i += 5
				} else {
					i++
				}
			default:
				i++
			}
		}
	}

	return output
}

func runPythonVersion() {
	tmpDir, err := os.MkdirTemp("", "pytemp")
	if err != nil {
		log.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	files := map[string]string{
		"main.py": `import sys, re, math, argparse
import xml.etree.ElementTree as ET
from svg.path import parse_path, Move, Line, CubicBezier, QuadraticBezier, Arc

def parse_size(size_str):
    try:
        return float(re.findall(r"[\d.]+", size_str)[0])
    except:
        return 0

def parse_transform(s):
    """Parse a SVG transform string into a list of (name, args)."""
    funcs = []
    for name, args_str in re.findall(r'(\w+)\(([^)]*)\)', s):
        parts = [float(v) for v in re.split(r'[,\s]+', args_str.strip()) if v]
        funcs.append((name, parts))
    return funcs

def apply_transforms(x, y, transforms):
    """Apply translate/scale/rotate transforms in order to point (x,y)."""
    x0, y0 = x, y
    for name, args in transforms:
        if name == 'translate':
            dx = args[0]
            dy = args[1] if len(args) > 1 else 0
            x0 += dx; y0 += dy
        elif name == 'scale':
            sx = args[0]
            sy = args[1] if len(args) > 1 else sx
            x0 *= sx; y0 *= sy
        elif name == 'rotate':
            a = math.radians(args[0])
            cx = args[1] if len(args) > 1 else 0
            cy = args[2] if len(args) > 2 else 0
            x0 -= cx; y0 -= cy
            xr = x0*math.cos(a) - y0*math.sin(a)
            yr = x0*math.sin(a) + y0*math.cos(a)
            x0, y0 = xr + cx, yr + cy
    return x0, y0

def approximate_path(path, transforms, segments=20):
    """Return a list of straight line segments [(x1,y1,x2,y2),...] from an SVG path."""
    out = []
    pos = complex(0, 0)
    for seg in path:
        if isinstance(seg, Move):
            pos = seg.end
        elif isinstance(seg, Line):
            out.append((pos, seg.end))
            pos = seg.end
        elif isinstance(seg, (CubicBezier, QuadraticBezier, Arc)):
            for i in range(segments):
                t0 = i / segments
                t1 = (i + 1) / segments
                out.append((seg.point(t0), seg.point(t1)))
            pos = seg.end
        else:
            out.append((pos, seg.end))
            pos = seg.end

    lines = []
    for p0, p1 in out:
        x1, y1 = apply_transforms(p0.real, p0.imag, transforms)
        x2, y2 = apply_transforms(p1.real, p1.imag, transforms)
        lines.append((x1, y1, x2, y2))
    return lines

def hex_to_rgb(h):
    h = h.lstrip('#')
    if len(h) == 3:
        h = ''.join([c*2 for c in h])
    return tuple(int(h[i:i+2], 16) for i in (0, 2, 4))

def interactive_input():
    svg_file = input("Enter path to SVG file: ").strip()
    verbose = input("Enable verbose output? (y/n): ").strip().lower() == 'y'
    rgb_input = input("Override default RGB color? (R,G,B or leave blank for 0,0,0): ").strip()
    rgb = tuple(map(int, rgb_input.split(','))) if rgb_input else (0, 0, 0)
    return svg_file, verbose, rgb

def main():
    parser = argparse.ArgumentParser(description="Convert SVG to .lrlogic format")
    parser.add_argument('svg_file', nargs='?', help="Path to SVG input file")
    parser.add_argument('--verbose', action='store_true', help="Enable verbose output")
    parser.add_argument('--rgb', nargs=3, type=int, metavar=('R', 'G', 'B'), help="Override RGB color (default 0,0,0)")

    args = parser.parse_args()

    if not args.svg_file:
        args.svg_file, args.verbose, rgb_override = interactive_input()
    else:
        rgb_override = tuple(args.rgb) if args.rgb else (0, 0, 0)

    svg_file = args.svg_file
    verbose = args.verbose

    tree = ET.parse(svg_file)
    root = tree.getroot()

    ns = ''
    if root.tag.startswith('{'):
        ns = root.tag[: root.tag.index('}') + 1]

	# Start measuring time
    start_time = time.time()

    output = ["LRFILE VERSION 2", "LRFILL OFF"]

    width = root.get('width')
    height = root.get('height')
    w = parse_size(width) if width else 0
    h = parse_size(height) if height else 0

    output.append(f"LRRESDEFINEX {int(round(w))}")
    output.append(f"LRRESDEFINEY {int(round(h))}")

    all_lines = []

    for g in root.findall(f'.//{ns}g'):
        tstr = g.get('transform', '').strip()
        transforms = parse_transform(tstr) if tstr else []

        for pe in g.findall(f'.//{ns}path'):
            d = pe.get('d')
            if not d:
                continue

            color = pe.get('stroke') or pe.get('fill') or '#000000'
            if color.lower() == 'none':
                color = pe.get('fill') or '#000000'

            r, g_col, b = hex_to_rgb(color)
            if color == '#000000':  # apply override only if default black
                r, g_col, b = rgb_override

            path = parse_path(d)
            lines = approximate_path(path, transforms, segments=20)
            for x1, y1, x2, y2 in lines:
                all_lines.append((x1, y1, x2, y2, r, g_col, b))

    xs = [x for l in all_lines for x in (l[0], l[2])]
    ys = [y for l in all_lines for y in (l[1], l[3])]
    min_x, min_y = min(xs), min(ys)
    shift_x, shift_y = -min_x, -min_y

    if verbose:
        print(f"Original width x height: {w} x {h}")
        print(f"Shift values: X={shift_x}, Y={shift_y}")
        print(f"Lines detected: {len(all_lines)}")

    for x1, y1, x2, y2, r, g_col, b in all_lines:
        nx1 = int(round(x1 + shift_x))
        ny1 = int(round(h - (y1 + shift_y)))  # Vertical flip
        nx2 = int(round(x2 + shift_x))
        ny2 = int(round(h - (y2 + shift_y)))  # Vertical flip
        output.append(f"{nx1},{ny1},{nx2},{ny2}..{r},{g_col},{b}")

    output.append("LREXIT")

    out_file = svg_file.rsplit('.', 1)[0] + '.lrlogic'
    with open(out_file, 'w') as f:
        f.write('\n'.join(output))

    print("Generated", out_file)
    # Measure the time elapsed
    end_time = time.time()
    elapsed_time = end_time - start_time
    print(f"Done! Time elapsed: {elapsed_time:.2f} seconds")

if __name__ == "__main__":
    main()
`,
	}

	var mainPyPath string

	for name, content := range files {
		fullPath := filepath.Join(tmpDir, name)
		if err := os.WriteFile(fullPath, []byte(content), 0755); err != nil {
			log.Fatalf("Failed to write %s: %v", name, err)
		}
		if name == "main.py" {
			mainPyPath = fullPath
		}
	}

	args := []string{mainPyPath, *fileFlag}

	if *rgbStr != "" {
		args = append(args, "--rgb", *rgbStr)
	}
	if *verbose {
		args = append(args, "--verbose")
	}

	cmd := exec.Command("python3", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		log.Fatalf("Python script failed: %v", err)
		if !confirmInstall() {
			fmt.Println("Dependency installation cancelled. Exiting.")
			os.Exit(1)
		}

		if _, err := exec.LookPath("pacman"); err == nil {
			if err := installSvgPathAUR(); err != nil {
				fmt.Println("AUR install failed:", err)
				fmt.Println("Falling back to pip install...")
				pipInstall()
			}
		} else {
			pipInstall()
		}

	}

	os.Exit(0)
}

func pipInstall() {
	cmd := exec.Command("python", "-m", "pip", "install", "-y", "svg.path")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to install svg.path via pip: %v", err)
	}
}

func confirmInstall() bool {
	fmt.Print("Do you want to install the svg.path dependency? (y/n): ")
	var response string
	fmt.Scanln(&response)
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}
