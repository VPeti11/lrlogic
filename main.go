package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type Point struct {
	X, Y int
}

type ColoredLine struct {
	Start, End Point
	R, G, B    int
}

func main() {
	filepathFlag := flag.String("file", "", "Path to the .lrlogic file (required)")
	nojpg := flag.Bool("nojpg", false, "Do not generate JPG output")
	nosvg := flag.Bool("nosvg", false, "Delete SVG output after generating JPG")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	flag.Parse()

	if *filepathFlag == "" {
		fmt.Println("Usage: lrlogic --file filename.lrlogic [--nojpg] [--nosvg] [--verbose]")
		os.Exit(1)
	}

	file, err := os.Open(*filepathFlag)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Detect file version
	if !scanner.Scan() {
		log.Fatal("File is empty!")
	}
	header := scanner.Text()
	isV2 := false
	if header == "LRFILE VERSION 2" {
		isV2 = true
		if *verbose {
			fmt.Println("Detected LRFILE VERSION 2")
		}
	} else if header == "LRLOGIC FILE FORMAT V1" {
		isV2 = false
		if *verbose {
			fmt.Println("Detected LRLOGIC FILE FORMAT V1")
		}
	} else {
		log.Fatal("Invalid file header!")
		os.Exit(1)
	}

	// Defaults
	width := 640
	height := 480
	marginTop := 20
	marginBottom := 20
	fontSize := 16
	curveStrength := 5

	var paths []string
	var topText, bottomText string
	var topLine, bottomLine bool
	var coloredLines []ColoredLine

	// Fill mode logic:
	fillMode := true // default fill mode
	if isV2 {
		fillMode = false // for v2 start off with no fill until LRFILL ON is encountered
	}

	lineNum := 1 // counting from line after header
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if *verbose {
			fmt.Printf("Processing line %d: %s\n", lineNum, line)
		}

		if line == "LREXIT" {
			if *verbose {
				fmt.Println("Found LREXIT, stopping parse.")
			}
			break
		}

		if strings.HasPrefix(line, "LRRESDEFINEX") {
			parts := strings.Fields(line)
			if len(parts) == 2 {
				if val, err := strconv.Atoi(parts[1]); err == nil {
					width = val
					if *verbose {
						fmt.Printf("Set width to %d\n", width)
					}
				}
			}
			continue
		}

		if strings.HasPrefix(line, "LRRESDEFINEY") {
			parts := strings.Fields(line)
			if len(parts) == 2 {
				if val, err := strconv.Atoi(parts[1]); err == nil {
					height = val
					if *verbose {
						fmt.Printf("Set height to %d\n", height)
					}
				}
			}
			continue
		}

		if strings.HasPrefix(line, "LRMARGIN") {
			parts := strings.Fields(line)
			if len(parts) == 3 {
				if val, err := strconv.Atoi(parts[1]); err == nil {
					marginTop = val
					if *verbose {
						fmt.Printf("Set marginTop to %d\n", marginTop)
					}
				}
				if val, err := strconv.Atoi(parts[2]); err == nil {
					marginBottom = val
					if *verbose {
						fmt.Printf("Set marginBottom to %d\n", marginBottom)
					}
				}
			}
			continue
		}

		if strings.HasPrefix(line, "LRFONTSIZE") {
			parts := strings.Fields(line)
			if len(parts) == 2 {
				if val, err := strconv.Atoi(parts[1]); err == nil {
					fontSize = val
					if *verbose {
						fmt.Printf("Set fontSize to %d\n", fontSize)
					}
				}
			}
			continue
		}

		if strings.HasPrefix(line, "LRCURVE") {
			parts := strings.Fields(line)
			if len(parts) == 2 {
				if val, err := strconv.Atoi(parts[1]); err == nil {
					curveStrength = val
					if *verbose {
						fmt.Printf("Set curveStrength to %d\n", curveStrength)
					}
				}
			}
			continue
		}

		if strings.HasPrefix(line, "LRTXT.Top") {
			topText = extractText(line)
			topLine = true
			if *verbose {
				fmt.Printf("Set topText: %s\n", topText)
			}
			continue
		}

		if strings.HasPrefix(line, "LRTXT.Bottom") {
			bottomText = extractText(line)
			bottomLine = true
			if *verbose {
				fmt.Printf("Set bottomText: %s\n", bottomText)
			}
			continue
		}

		if strings.HasPrefix(line, "LRFILL") && isV2 {
			parts := strings.Fields(line)
			if len(parts) == 2 {
				val := strings.ToUpper(parts[1])
				if val == "ON" {
					fillMode = true
					if *verbose {
						fmt.Println("Fill mode enabled")
					}
				} else if val == "OFF" {
					fillMode = false
					if *verbose {
						fmt.Println("Fill mode disabled")
					}
				}
			}
			continue
		}

		// Handle circles and squares for v2
		if isV2 && strings.HasPrefix(line, "LRCIRCLE") {
			// Format: LRCIRCLE x,y,radius..r,g,b
			parts := strings.SplitN(line, " ", 2)
			if len(parts) < 2 {
				if *verbose {
					fmt.Println("Skipping malformed LRCIRCLE line")
				}
				continue
			}
			params := parts[1]
			colorR, colorG, colorB := 0, 0, 0
			if strings.Contains(params, "..") {
				subparts := strings.Split(params, "..")
				params = subparts[0]
				rgbParts := strings.Split(subparts[1], ",")
				if len(rgbParts) == 3 {
					colorR, _ = strconv.Atoi(rgbParts[0])
					colorG, _ = strconv.Atoi(rgbParts[1])
					colorB, _ = strconv.Atoi(rgbParts[2])
				}
			}
			vals := strings.Split(params, ",")
			if len(vals) != 3 {
				if *verbose {
					fmt.Println("Skipping malformed LRCIRCLE parameters")
				}
				continue
			}
			x, _ := strconv.Atoi(vals[0])
			y, _ := strconv.Atoi(vals[1])
			radius, _ := strconv.Atoi(vals[2])
			y = height - y // invert y

			fillAttr := "none"
			if fillMode {
				fillAttr = fmt.Sprintf("rgb(%d,%d,%d)", colorR, colorG, colorB)
			}
			strokeAttr := fmt.Sprintf("rgb(%d,%d,%d)", colorR, colorG, colorB)

			paths = append(paths,
				       fmt.Sprintf(`<circle cx="%d" cy="%d" r="%d" fill="%s" stroke="%s" stroke-width="2"/>`,
						   x, y, radius, fillAttr, strokeAttr))
			if *verbose {
				fmt.Printf("Added circle at (%d,%d) radius %d color rgb(%d,%d,%d) fillMode %v\n",
					   x, y, radius, colorR, colorG, colorB, fillMode)
			}
			continue
		}

		if isV2 && strings.HasPrefix(line, "LRSQUARE") {
			// Format: LRSQUARE x,y,size..r,g,b
			parts := strings.SplitN(line, " ", 2)
			if len(parts) < 2 {
				if *verbose {
					fmt.Println("Skipping malformed LRSQUARE line")
				}
				continue
			}
			params := parts[1]
			colorR, colorG, colorB := 0, 0, 0
			if strings.Contains(params, "..") {
				subparts := strings.Split(params, "..")
				params = subparts[0]
				rgbParts := strings.Split(subparts[1], ",")
				if len(rgbParts) == 3 {
					colorR, _ = strconv.Atoi(rgbParts[0])
					colorG, _ = strconv.Atoi(rgbParts[1])
					colorB, _ = strconv.Atoi(rgbParts[2])
				}
			}
			vals := strings.Split(params, ",")
			if len(vals) != 3 {
				if *verbose {
					fmt.Println("Skipping malformed LRSQUARE parameters")
				}
				continue
			}
			x, _ := strconv.Atoi(vals[0])
			y, _ := strconv.Atoi(vals[1])
			size, _ := strconv.Atoi(vals[2])
			y = height - y // invert y

			fillAttr := "none"
			if fillMode {
				fillAttr = fmt.Sprintf("rgb(%d,%d,%d)", colorR, colorG, colorB)
			}
			strokeAttr := fmt.Sprintf("rgb(%d,%d,%d)", colorR, colorG, colorB)

			paths = append(paths,
				       fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="%d" fill="%s" stroke="%s" stroke-width="2"/>`,
						   x, y-size, size, size, fillAttr, strokeAttr))
			if *verbose {
				fmt.Printf("Added square at (%d,%d) size %d color rgb(%d,%d,%d) fillMode %v\n",
					   x, y, size, colorR, colorG, colorB, fillMode)
			}
			continue
		}

		// The rest is line parsing like before:

		colorR, colorG, colorB := 0, 0, 0
		if strings.Contains(line, "..") {
			parts := strings.Split(line, "..")
			line = parts[0]
			rgbParts := strings.Split(parts[1], ",")
			if len(rgbParts) == 3 {
				colorR, _ = strconv.Atoi(rgbParts[0])
				colorG, _ = strconv.Atoi(rgbParts[1])
				colorB, _ = strconv.Atoi(rgbParts[2])
			}
		}

		parts := strings.Split(line, ",")
		if len(parts) != 4 {
			if *verbose {
				fmt.Printf("Skipping malformed line: %s\n", line)
			}
			continue
		}

		x1, _ := strconv.Atoi(parts[0])
		y1, _ := strconv.Atoi(parts[1])
		x2, _ := strconv.Atoi(parts[2])
		y2, _ := strconv.Atoi(parts[3])

		y1 = height - y1
		y2 = height - y2

		coloredLines = append(coloredLines, ColoredLine{
			Start: Point{x1, y1},
			End:   Point{x2, y2},
			R:     colorR,
			G:     colorG,
			B:     colorB,
		})
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	// Group and process lines
	groups := make(map[string][]ColoredLine)
	for _, line := range coloredLines {
		key := fmt.Sprintf("%d,%d,%d", line.R, line.G, line.B)
		groups[key] = append(groups[key], line)
	}

	for key, lines := range groups {
		if len(lines) < 4 {
			for _, l := range lines {
				paths = append(paths, curveLine(l.Start, l.End, curveStrength, l.R, l.G, l.B))
			}
			continue
		}

		used := make([]bool, len(lines))
		var chain []Point
		current := lines[0].Start
		chain = append(chain, current)
		used[0] = true
		found := 1

		for found < 4 {
			madeProgress := false
			for i, l := range lines {
				if used[i] {
					continue
				}
				if l.Start == current {
					chain = append(chain, l.End)
					current = l.End
					used[i] = true
					found++
					madeProgress = true
					break
				} else if l.End == current {
					chain = append(chain, l.Start)
					current = l.Start
					used[i] = true
					found++
					madeProgress = true
					break
				}
			}
			if !madeProgress {
				break
			}
		}

		// Close loop and check
		chain = append(chain, chain[0])
		if len(chain) == 5 && chain[0] == chain[4] {
			pointsStr := ""
			for i := 0; i < 4; i++ {
				pointsStr += fmt.Sprintf("%d,%d ", chain[i].X, chain[i].Y)
			}
			fillColor := fmt.Sprintf("rgb(%s)", key)
			if isV2 {
				if fillMode {
					paths = append(paths, fmt.Sprintf(
						`<polygon points="%s" fill="%s" stroke="black" stroke-width="1"/>`,
				       strings.TrimSpace(pointsStr), fillColor))
				} else {
					// Just stroke lines
					for _, l := range lines {
						paths = append(paths, curveLine(l.Start, l.End, curveStrength, l.R, l.G, l.B))
					}
				}
			} else {
				// V1 always fill
				paths = append(paths, fmt.Sprintf(
					`<polygon points="%s" fill="%s" stroke="black" stroke-width="1"/>`,
				      strings.TrimSpace(pointsStr), fillColor))
			}
		} else {
			for _, l := range lines {
				paths = append(paths, curveLine(l.Start, l.End, curveStrength, l.R, l.G, l.B))
			}
		}
	}

	baseName := strings.TrimSuffix(filepath.Base(*filepathFlag), filepath.Ext(*filepathFlag))
	svgName := baseName + ".svg"
	jpgName := baseName + ".jpg"

	// Write SVG
	output, err := os.Create(svgName)
	if err != nil {
		log.Fatalf("Failed to create %s: %v", svgName, err)
	}
	defer output.Close()

	fmt.Fprintf(output, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d">`+"\n", width, height)
	fmt.Fprintf(output, `<rect width="%d" height="%d" fill="white"/>`+"\n", width, height)

	if topText != "" {
		y := marginTop + fontSize
		if topLine {
			fmt.Fprintf(output, `<line x1="0" y1="%d" x2="%d" y2="%d" stroke="black" stroke-width="1"/>`+"\n", y+4, width, y+4)
		}
		fmt.Fprintf(output, `<text x="10" y="%d" font-size="%d" fill="black">%s</text>`+"\n", y, fontSize, topText)
	}

	if bottomText != "" {
		y := height - marginBottom
		if bottomLine {
			fmt.Fprintf(output, `<line x1="0" y1="%d" x2="%d" y2="%d" stroke="black" stroke-width="1"/>`+"\n", y-fontSize-4, width, y-fontSize-4)
		}
		fmt.Fprintf(output, `<text x="10" y="%d" font-size="%d" fill="black">%s</text>`+"\n", y, fontSize, bottomText)
	}

	for _, path := range paths {
		fmt.Fprintln(output, path)
	}

	fmt.Fprintln(output, `</svg>`)
	fmt.Printf("Generated %s successfully\n", svgName)

	// JPG conversion
	if !*nojpg {
		if checkCommand("rsvg-convert") {
			err = exec.Command("rsvg-convert", "-o", jpgName, svgName).Run()
		} else if checkCommand("convert") {
			err = exec.Command("convert", svgName, jpgName).Run()
		} else {
			fmt.Println("No rsvg-convert binary found!")
			os.Exit(1)
			return
		}

		if err != nil {
			log.Printf("JPG conversion failed with error: %v", err)
			os.Exit(1)
		} else {
			fmt.Printf("Generated %s successfully.\n", jpgName)
		}
	}

	if *nosvg {
		err := os.Remove(svgName)
		if err != nil {
			log.Printf("Failed to remove SVG file: %v", err)
			os.Exit(1)
		} else {
			fmt.Printf("Removed %s SVG file\n", svgName)
		}
	}
}



func extractText(line string) string {
	start := strings.Index(line, "'")
	end := strings.LastIndex(line, "'")
	if start == -1 || end == -1 || start == end {
		return ""
	}
	return line[start+1 : end]
}

func curveLine(start, end Point, strength, r, g, b int) string {
	mx := (start.X + end.X) / 2
	my := (start.Y + end.Y) / 2
	color := fmt.Sprintf("rgb(%d,%d,%d)", r, g, b)
	return fmt.Sprintf(`<path d="M %d %d Q %d %d %d %d" stroke="%s" fill="none" stroke-width="2"/>`,
			   start.X, start.Y, mx, my-strength, end.X, end.Y, color)
}

func checkCommand(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
