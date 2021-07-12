#!/usr/bin/env bash

# Check required env variables
flag=0
if [[ -z "$SLACK_TOKEN" ]]; then
    flag=1
    missing_secret="SLACK_TOKEN"
fi

if [[ "$flag" -eq 1 ]]; then
    printf "[\e[0;31mERROR\e[0m] Secret \`$missing_secret\` is missing. Please add it to this action for proper execution.\nRefer https://github.com/giantswarm/update-slack-info for more information.\n"
    exit 1
fi

update-slack-info
