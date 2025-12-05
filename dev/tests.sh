#!/usr/bin/env bash

# This script runs all project tests
# with coverage reports.

set -e

go test -coverprofile=coverage.txt -covermode=atomic ./...
