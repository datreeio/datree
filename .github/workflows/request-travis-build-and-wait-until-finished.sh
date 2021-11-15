#!/bin/bash

GH_ORGANIZATION_NAME="datreeio"
GH_REPOSITORY_NAME="datree"
GH_SLUG=${GH_ORGANIZATION_NAME}%2F${GH_REPOSITORY_NAME}

echo "Requesting a travis build..."

request_build_response=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -H "Travis-API-Version: 3" \
  -H "Authorization: token ${TRAVIS_API_TOKEN}" \
  -d "$REQUEST_BUILD_BODY" \
  "https://api.travis-ci.com/repo/${GH_SLUG}/requests")

REQUEST_ID=$(echo $request_build_response | jq '.request.id')

echo "Travis request id: ${REQUEST_ID}"

is_travis_build_already_created=false

while true; do
  sleep 10
  echo "Polling travis build status..."

  response=$(curl -s -X GET \
    -H "Content-Type: application/json" \
    -H "Accept: application/json" \
    -H "Travis-API-Version: 3" \
    -H "Authorization: token ${TRAVIS_API_TOKEN}" \
    "https://api.travis-ci.com/repo/${GH_SLUG}/request/${REQUEST_ID}")

  build_status=$(echo $response | jq '.builds[0].state')
  build_request_state=$(echo $response | jq '.state')

  if [ $build_request_state = "\"pending\"" ]; then
    echo "Build request pending..."
  elif [ $build_status = "\"passed\"" ]; then
    echo "Travis build passed!"
    exit 0
  elif [ $build_status = "\"failed\"" ]; then
    echo "Travis build failed!"
    exit 1
  elif [ $build_status = "\"canceled\"" ]; then
    echo "Travis build canceled!"
    exit 1
  elif [ $build_status = "\"created\"" ]; then
    if [ $is_travis_build_already_created = true ]; then
      echo "Travis build queued..."
    else
      echo "Travis build created!"
    fi
    is_travis_build_already_created=true
  elif [ $build_status = "\"started\"" ]; then
    echo "Travis build in progress..."
  else
    echo "Travis build unknown state!"
    echo "build_status:${build_status}"
    echo "REQUEST_ID:${REQUEST_ID}"
    echo "response:${response}"
    exit 1
  fi
done
