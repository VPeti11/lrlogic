#!/bin/bash
rm *.svg
rm *.jpg
set -e
go build main.go
mv main lrlogic
#read -p "Press enter to run full test"
cp ./Tests/*.lrlogic ./
chmod +x runall.sh
./runall.sh
rm *.lrlogic
