#!/usr/bin/env bash
set -e

if [ -z "${PORT}" ]; then
    echo "PORT env variable was not set to a specific variable. Using the default one - 8080"
    export PORT=8080
else
    sed -i "s/ENV PORT [0-9]*/ENV PORT ${PORT}/g" ../Dockerfile
fi

docker build -t integration-test ./..

docker run -d -p${PORT}:${PORT} integration-test