import os
import sys
import random
import argparse
import subprocess
import time
from concurrent.futures import ProcessPoolExecutor, as_completed

def detect_os():
    return 'lrlogic.exe' if os.name == 'nt' else './lrlogic'

def worker(i, width, height, maxshapes, shape_weights, render, cleanup):
    filename = f"output_{i}.lrlogic"
    shape_types = ['line', 'circle', 'square']
    shape_count = random.randint(1, maxshapes)
    curve = random.randint(0, 10)
    fill_mode = False

    with open(filename, 'w') as f:
        f.write("LRFILE VERSION 2\n")
        f.write(f"LRRESDEFINEX {width}\n")
        f.write(f"LRRESDEFINEY {height}\n")
        f.write(f"LRCURVE {curve}\n")
        for _ in range(shape_count):
            if random.random() < 0.2:
                fill_mode = not fill_mode
                f.write(f"LRFILL {'ON' if fill_mode else 'OFF'}\n")
            shape = random.choices(shape_types, weights=shape_weights)[0]
            r, g, b = [random.randint(0, 255) for _ in range(3)]
            if shape == 'line':
                x1, y1 = random.randint(0, width), random.randint(0, height)
                x2, y2 = random.randint(0, width), random.randint(0, height)
                f.write(f"{x1},{y1},{x2},{y2}..{r},{g},{b}\n")
            elif shape == 'circle':
                x, y = random.randint(0, width), random.randint(0, height)
                radius = random.randint(10, min(width, height) // 4)
                f.write(f"LRCIRCLE {x},{y},{radius}..{r},{g},{b}\n")
            elif shape == 'square':
                x, y = random.randint(0, width), random.randint(0, height)
                size = random.randint(10, min(width, height) // 3)
                f.write(f"LRSQUARE {x},{y},{size}..{r},{g},{b}\n")
        f.write("LREXIT\n")

    if render:
        lrlogic_exec = detect_os()
        if not os.path.exists(lrlogic_exec):
            print(f"Warning: {lrlogic_exec} not found. Skipping rendering.")
            return
        subprocess.run([lrlogic_exec, '-file', filename, '-verbose', '-nosvg'], stdout=subprocess.DEVNULL)

    os.rename(filename, os.path.join("randomgen", filename))
    jpgname = f"output_{i}.jpg"
    if os.path.exists(jpgname):
        os.rename(jpgname, os.path.join("randomgen", jpgname))

    if cleanup:
        logic_path = os.path.join("randomgen", filename)
        if os.path.exists(logic_path):
            os.remove(logic_path)

def main():
    os.makedirs("randomgen", exist_ok=True)

    parser = argparse.ArgumentParser()
    parser.add_argument('--count', type=int, default=100, help="Number of files to generate (default: 100)")
    parser.add_argument('--maxshapes', type=int, default=100, help="Maximum number of shapes per file (default: 100)")
    parser.add_argument('--width', type=int)
    parser.add_argument('--height', type=int)
    parser.add_argument('--render', action='store_true', default=True, help="Render files (default: enabled)")
    parser.add_argument('--cleanup', action='store_true')
    parser.add_argument('--ratio', type=str, help="Shape ratio line,circle,square (default: 5,2,3)")
    args = parser.parse_args()

    if args.count == 100:
        args.count = int(input(f"How many files to generate? (default: {args.count}): ").strip() or args.count)

    if args.maxshapes == 100:
        args.maxshapes = int(input(f"Max shapes per file? (default: {args.maxshapes}): ").strip() or args.maxshapes)

    if args.width is None:
        args.width = int(input("Canvas width (default 1024): ").strip() or 1024)

    if args.height is None:
        args.height = int(input("Canvas height (default 768): ").strip() or 768)

    if '--render' not in sys.argv:
        render_input = input("Render files? (y/n, default: y): ").strip().lower()
        args.render = not (render_input == 'n' or render_input == 'no')

    if '--cleanup' not in sys.argv:
        cleanup_input = input("Delete .lrlogic after rendering? (y/n): ").strip().lower()
        args.cleanup = cleanup_input in ['y', 'yes']

    if args.ratio is None:
        ratio_input = input("Shape ratio line,circle,square (default: 5,2,3): ").strip()
        args.ratio = ratio_input or "5,2,3"

    print(f"Using shape ratio: {args.ratio}")

    shape_weights = [5, 2, 3]
    parts = args.ratio.split(',')
    if len(parts) == 3:
        try:
            shape_weights = [int(p) for p in parts]
        except ValueError:
            pass  # fallback to default weights

    start_time = time.time()

    with ProcessPoolExecutor() as executor:
        tasks = [
            executor.submit(worker, i+1, args.width, args.height, args.maxshapes, shape_weights, args.render, args.cleanup)
            for i in range(args.count)
        ]
        for _ in as_completed(tasks):
            pass

    elapsed_time = time.time() - start_time
    print(f"Done! Time elapsed: {elapsed_time:.2f} seconds")

if __name__ == "__main__":
    main()
