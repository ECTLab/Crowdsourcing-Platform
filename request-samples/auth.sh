#!bin/bash

curl --request POST \
  --url http://localhost:9000/generate-token \
  --header 'Content-Type: application/json' \
  --header 'User-Agent: insomnia/11.0.2' \
  --data '{
	"name":"demo"
}'