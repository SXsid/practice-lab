#!/usr/bin/env bash
checkHealth() {
    URL=$1
    if [ -z "$URL" ];then
        echo "usage: ./health.sh <url of the site>"
        return 1
    fi

    CODE=$(curl -s -o /dev/null -w "%{http_code}" "$URL")
    if [ "$CODE" -ge 200 ] && [ "$CODE" -lt 500 ];then
        echo "Service is up"
    else
        echo "service is down got code $CODE"
        return 1
    fi
}

checkHealth "$@"
