import sys, re, math, argparse
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

if __name__ == "__main__":
    main()
