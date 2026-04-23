#!/usr/bin/env bash

SERVICES=("reids" "docker" "postgres")

for service in "${SERVICES[@]}" ;do
    systemctl  is-active --quiet "$service" \
        && echo "${service} is active"\
        || echo "${service} is down"
 done
