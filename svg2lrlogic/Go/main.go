package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Line struct {
	X1, Y1, X2, Y2 int
	R, G, B        int
}

func main() {
	fileFlag := flag.String("file", "", "Path to the SVG file (required)")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	flag.Parse()

	if *fileFlag == "" {
		fmt.Println("Usage: svg2lrlogic --file input.svg [--verbose]")
		os.Exit(1)
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

