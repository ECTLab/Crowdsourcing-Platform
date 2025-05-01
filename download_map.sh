#!/bin/bash

set -e

FILE_ID="1H6LVHWFIl6PLV9WkrKlB0cxXqKZcZeuH"
DEST_DIR="map-viewer"
DEST_NAME="iran.osm.pbf"

mkdir -p "$DEST_DIR"

echo "Downloading file from Google Drive..."
CONFIRM=$(curl -sc /tmp/gcookie "https://drive.google.com/uc?export=download&id=${FILE_ID}" | \
         grep -o 'confirm=[^&]*' | sed 's/confirm=//')
curl -Lb /tmp/gcookie "https://drive.google.com/uc?export=download&confirm=${CONFIRM}&id=${FILE_ID}" \
     -o "${DEST_DIR}/${DEST_NAME}"

echo "Download complete: ${DEST_DIR}/${DEST_NAME}"

