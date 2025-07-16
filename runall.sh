#!/bin/bash

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
