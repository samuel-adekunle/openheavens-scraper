#!/bin/bash
set -e
TODAY=$(date +'%-d-%B-%Y' | tr '[:upper:]' '[:lower:]')
echo "Today is ${TODAY}"
if compgen -G "${TODAY}*.txt" >/dev/null; then
    echo "Already scraped for today"
else
    echo "Scraping for today"
    go run main.go post.go
fi
