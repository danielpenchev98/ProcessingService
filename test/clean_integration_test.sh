#!/usr/bin/env bash
set -e

docker stop $(docker ps -q --filter ancestor=integration-test )