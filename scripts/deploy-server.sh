#!/bin/bash
set -e
curl -X POST ${RENDER_DEPLOY_HOOK} && echo "Render deploy hook sent"