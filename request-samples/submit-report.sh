#!bin/bash

TOKEN=PASTE_YOUR_TOKEN

curl --request POST \
  --url http://localhost:8001/crowdsourcing/in-ride-report \
  --header 'Content-type: application/json' \
  --header "api-key: ${TOKEN}" \
  --data '{
  "type": "POLICE",
  "engagement_location_time": {
    "latitude":     35.706861,
    "longitude": 51.336971,
    "time": 394
  },
  "submit_location_time": {
    "latitude":     35.706861,
    "longitude": 51.336971,
    "time": 530
  }
}'