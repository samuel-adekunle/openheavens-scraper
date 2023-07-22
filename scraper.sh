#!/bin/bash
set -e
go version
echo $(date)
TODAY=$(date +'%-d-%B-%Y' | tr '[:upper:]' '[:lower:]')
echo "Today is ${TODAY}"
if compgen -G "./posts/${TODAY}*.txt" >/dev/null; then
    echo "Already scraped for today"
else
    echo "Scraping for today"
    go run main.go post.go
    curl -X POST ${RENDER_DEPLOY_HOOK}
fi
