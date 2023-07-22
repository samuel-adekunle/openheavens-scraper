#!/bin/bash
set -e
if [ -f .deploy ]; then
    echo "Deploying..."
    rm .deploy
    curl -X POST ${RENDER_DEPLOY_HOOK} && echo "Render deploy hook sent"
fi