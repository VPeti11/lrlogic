#!/bin/bash

# Create target folder if it doesn't exist
mkdir -p randomgen

read -p "How many .lrlogic files should I generate? " file_count
read -p "What's the maximum number of shapes per file? " max_shapes
read -p "Do you want to render the files? (yes/no) " render_choice
read -p "Delete .lrlogic files after rendering? (yes/no) " delete_choice

for ((i=1; i<=file_count; i++)); do
    shape_count=$((RANDOM % max_shapes + 1))
    curve=$((RANDOM % 11))
    filename="output_${i}.lrlogic"

    echo "Generating $filename with $shape_count shapes and curve $curve..."
    python3 randomgen.py --count 1 --width 1024 --height 768 --curve $curve --shapes $shape_count

    if [[ "$render_choice" == "yes" ]]; then
        echo "Rendering $filename..."
        ./lrlogic -file "$filename" -verbose -nosvg
    fi

    # Move generated files to randomgen folder
    mv "$filename" randomgen/
    if [[ -f "output_${i}.jpg" ]]; then
        mv "output_${i}.jpg" randomgen/
    fi

    if [[ "$delete_choice" == "yes" ]]; then
        echo "Deleting $filename from current directory (already moved)..."
        # Already moved, so nothing to remove from here
        :
    fi
done

echo "All done! Check your 'randomgen' folder for the results. 🎨"
