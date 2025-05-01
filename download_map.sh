#!/bin/bash

set -e

FILE_ID="1H6LVHWFIl6PLV9WkrKlB0cxXqKZcZeuH"
DEST_DIR="map-viewer"
DEST_NAME="iran.osm.pbf"

mkdir -p "$DEST_DIR"

echo "Downloading file from Google Drive..."
CONFIRM=$(curl -sc /tmp/gcookie "https://drive.google.com/uc?export=download&id=${FILE_ID}" | \
         grep -o 'confirm=[^&]*' | sed 's/confirm=//')
curl -Lb /tmp/gcookie "https://drive.usercontent.google.com/download?id=1H6LVHWFIl6PLV9WkrKlB0cxXqKZcZeuH&export=download&authuser=0&confirm=t&uuid=4d05e5d7-3fce-42a8-8b1e-93ad98be7a3a&at=APcmpozcRqIXMQIqW9y0wnaLAOuO:1746129325451" \
     -o "${DEST_DIR}/${DEST_NAME}"

echo "Download complete: ${DEST_DIR}/${DEST_NAME}"

