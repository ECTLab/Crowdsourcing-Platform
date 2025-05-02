#!bin/bash

TOKEN=PASTE_YOUR_TOKEN_HERE
REPORT_ID=PASTE_THE_REPORT_ID_HERE

curl --request POST \
  --url http://localhost:8001/crowdsourcing/in-ride-report/${REPORT_ID}/confirm \
  --header 'Content-type: application/json' \
  --header "api-key: ${TOKEN}" \
  --data '{
	"type":"POLICE",
 "confirmed": false
}'