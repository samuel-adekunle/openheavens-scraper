#!/bin/bash
set -e
TODAY=$(date +'%-d-%B-%Y' | tr '[:upper:]' '[:lower:]')
go run main.go -date=${TODAY}