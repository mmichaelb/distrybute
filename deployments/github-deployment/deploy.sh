#!/bin/bash
set -e
if [ -z "$GITHUB_USERNAME" ] || [ -z "$GITHUB_ACCESS_TOKEN" ]; then
  echo "no github username (GITHUB_USERNAME) or github access token (GITHUB_ACCESS_TOKEN) set"
  exit 1
fi

function send_deployment_status() {
  local state
  state=$1
  local description
  description=$2
  curl -s -X POST -u "$GITHUB_USERNAME:$GITHUB_ACCESS_TOKEN" -H "Content-Type: application/json" -H "Accept: application/vnd.github.v3+json" \
    -d "{\"state\":\"$state\",\"description\":\"$description\"}" "$PAYLOAD_DEPLOYMENT_STATUSES_URL"
}

function check_exit_code() {
  local exit_code=$?
  if [ "$exit_code" != "0" ]; then
    send_deployment_status "error" "previous command exited with exit code $exit_code"
    exit $exit_code
  fi
}

PAYLOAD="$1"
PAYLOAD_ACTION=$(echo "$PAYLOAD" | jq -r '.action')
if [ "$PAYLOAD_ACTION" != 'created' ]; then
  printf "received %s as action in payload - stopping script\n" "$PAYLOAD_ACTION"
  exit 0
fi

PAYLOAD_DEPLOYMENT_TASK=$(echo "$PAYLOAD" | jq -r '.deployment.task')
if [ "$PAYLOAD_DEPLOYMENT_TASK" != 'deploy' ]; then
  printf "received %s as deployment.task in payload - stopping script\n" "$PAYLOAD_DEPLOYMENT_TASK"
  exit 0
fi

PAYLOAD_DEPLOYMENT_ID=$(echo "$PAYLOAD" | jq -r '.deployment.id')
PAYLOAD_DEPLOYMENT_SHA=$(echo "$PAYLOAD" | jq -r '.deployment.sha')
PAYLOAD_DEPLOYMENT_REF=$(echo "$PAYLOAD" | jq -r '.deployment.ref')
PAYLOAD_DEPLOYMENT_DESCRIPTION=$(echo "$PAYLOAD" | jq -r '.deployment.description')
PAYLOAD_REPOSITORY_FULL_NAME=$(echo "$PAYLOAD" | jq -r '.repository.full_name')
PAYLOAD_DEPLOYMENT_PAYLOAD=$(echo "$PAYLOAD" | jq -r '.deployment.payload')

printf "deploying incoming github deployment (repository=%s, id=%s, sha=%s, ref=%s, payload=%s, description=%s)\n" \
  "$PAYLOAD_REPOSITORY_FULL_NAME" "$PAYLOAD_DEPLOYMENT_ID" "$PAYLOAD_DEPLOYMENT_SHA" "$PAYLOAD_DEPLOYMENT_REF" "$PAYLOAD_DEPLOYMENT_PAYLOAD" "$PAYLOAD_DEPLOYMENT_DESCRIPTION"

PAYLOAD_DEPLOYMENT_STATUSES_URL=$(echo "$PAYLOAD" | jq -r '.deployment.statuses_url')
send_deployment_status 'in_progress' 'pulling docker image...'
docker-compose pull distrybute || check_exit_code
send_deployment_status 'in_progress' 'restarting docker container'
docker-compose restart distrybute || check_exit_code
send_deployment_status 'success' 'deployment done'
