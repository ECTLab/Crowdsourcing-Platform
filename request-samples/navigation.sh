#!bin/bash

TOKEN=PASTE_YOURE_TOKEN_HERE

curl --request POST \
  --url http://localhost:8080/navigation/get-route \
  --header 'Content-type: application/json' \
  --header "api-key: ${TOKEN}" \
  --data '{
  "vehicleType": "car",
  "origin": {
    "latitude":     35.707804,
    "longitude": 51.336780
  },
  "destination": {
    "latitude":     35.706735,
    "longitude": 51.336279
  }
}'