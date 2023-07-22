#!/bin/bash
set -e
TODAY=$(date +'%-d-%B-%Y' | tr '[:upper:]' '[:lower:]')
if compgen -G "./posts/${TODAY}*.txt" >/dev/null; then
    echo "Already scraped for today, ${TODAY}"
else
    echo "Scraping for today, ${TODAY}"
    go run main.go post.go -date=${TODAY} -output=./posts/${TODAY}.txt
    echo "Scraped for today, ${TODAY}"
fi
