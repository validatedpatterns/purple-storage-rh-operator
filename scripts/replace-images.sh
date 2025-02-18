#!/bin/bash
set -e

if [ $# -lt 1 ]; then
    echo "Usage: $0 <file1> ... <fileN>"
    echo "Please pass at least one yaml file to use"
    exit 1
fi 


# Define the replacement list as an associative array
declare -A replacements=(
    ["cp.icr.io/cp/spectrum/scale/data-access"]="quay.io/rhsysdeseng/cp/spectrum/scale/data-access"
    ["cp.icr.io/cp/spectrum/scale/data-management"]="quay.io/rhsysdeseng/cp/spectrum/scale/data-management"
    ["cp.icr.io/cp/spectrum/scale"]="quay.io/rhsysdeseng/cp/spectrum/scale"
    ["cp.icr.io/cp/spectrum/scale/csi"]="quay.io/rhsysdeseng/cp/spectrum/scale/csi"
)

# Process each file passed as argument
for file in "$@"; do
    if [ -f "$file" ]; then
        echo "Processing file: $file"

        # Loop through the replacement list and perform replacements
        for search in "${!replacements[@]}"; do
            replace="${replacements[$search]}"
            sed -i "s|$search|$replace|g" "$file"
        done

        echo "Replacements completed for: $file"
    else
        echo "File not found: $file"
    fi
done

echo "All files processed."
