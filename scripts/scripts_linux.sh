#!/bin/bash
cd ..
# Menu
echo "Select a function to run:"
echo "1) Cleanup"
echo "2) Full Test"
echo "3) Makerandom"
echo "4) Render all"
read -p "Enter your choice (1/2/3/4): " choice

case $choice in
  1)
    echo "Running Cleanup..."
    rm *.jpg
    rm *.svg
    rm *.lrlogic
    echo "Cleanup complete."
    ;;
  2)
    echo "Running Full Test..."
    rm *.svg
    rm *.jpg
    set -e
    go build main.go
    mv main lrlogic
    cp ./Tests/*.lrlogic ./
PROGRAM="./lrlogic"

read -rp "Do you want to keep SVG files after conversion? (y/n): " keep_svg
keep_svg=${keep_svg,,} # to lowercase

rm *.svg
rm *.jpg

for file in *.lrlogic; do
  if [[ -f "$file" ]]; then
    echo "Processing $file..."
    if [[ "$keep_svg" == "n" ]]; then
      $PROGRAM --file "$file" --nosvg --verbose
    else
      $PROGRAM --file "$file" --verbose
    fi
    echo ""
  fi
done
    rm *.lrlogic
    echo "Full Test complete."
    ;;
  3)
    echo "Running Makerandom..."

        python3 randomgen.py

    echo "Random file generation complete."
    ;;
  4)
    echo "Rendering all in directory"
    PROGRAM="./lrlogic"

    read -rp "Keep SVG files after conversion? (y/n): " keep_svg
    keep_svg=${keep_svg,,} # to lowercase

    rm *.svg
    rm *.jpg

    for file in *.lrlogic; do
        if [[ -f "$file" ]]; then
            echo "Processing $file..."
            if [[ "$keep_svg" == "n" ]]; then
                $PROGRAM --file "$file" --nosvg --verbose
            else
                $PROGRAM --file "$file" --verbose
            fi
            echo ""
        fi
    done

    echo "All files processed."
    ;;
  *)
    echo "Invalid choice. Exiting..."
    exit 1
    ;;
esac
