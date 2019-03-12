#!/bin/bash

if [ "$PUSHOVER_TOKEN" == "" ]; then
  echo "Missing pushover token"
  exit 1
fi

if [ "$PUSHOVER_USER" == "" ]; then
  echo "Missing pushover user"
  exit 1
fi

message=$(</dev/stdin)

if [ "$message" == "" ]; then
  echo "Missing input"
  exit 1
fi


title=$(echo "${message}" | head -n 1)
message=$(echo "${message}" | tail -n +2)

curl -s \
  --form-string "token=$PUSHOVER_TOKEN" \
  --form-string "user=$PUSHOVER_USER" \
  --form-string "title=$title" \
  --form-string "message=$message" \
  https://api.pushover.net/1/messages.json | jq '.status'