#!/bin/bash
set -e

if [ $# -ne 1 ]; then
    echo "Usage: $0 <file>"
    echo "Please pass one yaml file to use"
    exit 1
fi 

file=$1
if [ ! -f "$file" ]; then
    echo "File not found: ${file}"
    exit 1
fi

cat <<EOF
kind: ImageSetConfiguration
apiVersion: mirror.openshift.io/v2alpha1
mirror:
  additionalImages
EOF
while IFS= read -r line; do
    echo "  - name: $line"
done < <(grep quay.io "${file}" | cut -f2- -d: | awk '{ print $1 }')
