#!/bin/bash

cp templates/presets_template.go cmd/presets/presets.go

echo "// GetAll get all presets" >> cmd/presets/presets.go
echo "func GetAll() map[string]map[string]string {" >> cmd/presets/presets.go
echo "	var presets = make(map[string]map[string]string)" >> cmd/presets/presets.go

for folder in presets/*/; do
    if [ ! -d $folder ]; then
		continue 1
    fi

    echo "Found folder $folder"

	preset=$(basename $folder)

	echo "	presets[\"$preset\"] = map[string]string{" >> cmd/presets/presets.go

    for file in $folder/*; do
		fileName=$(basename $file)
		content=$(cat $file)
		echo "		\"$fileName\": \`$content\`," >> cmd/presets/presets.go
		echo "  Parsed file: $fileName"
	done

	echo "	}" >> cmd/presets/presets.go
done

echo "	return presets" >> cmd/presets/presets.go
echo "}" >> cmd/presets/presets.go

echo "Finished building cmd/presets/presets.go"
