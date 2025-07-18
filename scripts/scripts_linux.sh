#!/bin/bash
cd ..
# Menu
echo "Select a function to run:"
echo "1) Cleanup"
echo "2) Full Test"
echo "3) Makerandom"
echo "4) Render all"
echo "5) Timed render"
echo "6) Compile LRLogic"
echo "7) Compile SVG2LR - Go"
read -p "Enter your choice (1/2/3/4/5/6/7): " choice

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
    # Capture start time
    start_time=$(date +%s)

    rm *.svg
    rm *.jpg
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

    # Capture end time and calculate elapsed time
    end_time=$(date +%s)
    elapsed_time=$((end_time - start_time))

    echo "Full Test complete."
    echo "Elapsed Time: $elapsed_time seconds."
    ;;
  3)
    echo "Running Makerandom..."
    python3 randomgen.py
    echo "Random file generation complete."
    ;;
  4)
    echo "Rendering all in directory"
    # Capture start time
    start_time=$(date +%s)

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

    # Capture end time and calculate elapsed time
    end_time=$(date +%s)
    elapsed_time=$((end_time - start_time))

    echo "All files processed."
    echo "Elapsed Time: $elapsed_time seconds."
    ;;

  5)
    echo "Enter the file path to render "
    read -rp "File Path: " file_path


    if [[ ! -f "$file_path" ]]; then
      echo "File not found. Exiting..."
      exit 1
    fi


    read -rp "Do you want to keep SVG file after conversion? (y/n): " keep_svg
    keep_svg=${keep_svg,,} # to lowercase

    # Capture start time
    start_time=$(date +%s)

    PROGRAM="./lrlogic"


    if [[ "$keep_svg" == "n" ]]; then
      $PROGRAM --file "$file_path" --nosvg --verbose
    else
      $PROGRAM --file "$file_path" --verbose
    fi

    # Capture end time and calculate elapsed time
    end_time=$(date +%s)
    elapsed_time=$((end_time - start_time))

    echo "Rendering complete."
    echo "Elapsed Time: $elapsed_time seconds."
    ;;
  
  6)
    # Capture start time
    start_time=$(date +%s)
    go build main.go
    mv main lrlogic
    # Capture end time and calculate elapsed time
    end_time=$(date +%s)
    elapsed_time=$((end_time - start_time))
    echo "Elapsed Time: $elapsed_time seconds."
    ;;
  7)
    # Capture start time
    start_time=$(date +%s)
    cd ./svg2lrlogic/Go
    go build main.go
    mv main svg2lr
    # Capture end time and calculate elapsed time
    end_time=$(date +%s)
    elapsed_time=$((end_time - start_time))
    echo "Elapsed Time: $elapsed_time seconds."
    ;;
  *)
    echo "Invalid choice. Exiting..."
    exit 1
    ;;
esac
