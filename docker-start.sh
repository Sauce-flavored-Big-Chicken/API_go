#!/bin/sh
set -eu

/app/server &
api_pid=$!

nginx -g "daemon off;" &
nginx_pid=$!

term_handler() {
    kill -TERM "$api_pid" "$nginx_pid" 2>/dev/null || true
    wait "$api_pid" 2>/dev/null || true
    wait "$nginx_pid" 2>/dev/null || true
    exit 0
}

trap term_handler INT TERM

wait "$nginx_pid"
kill -TERM "$api_pid" 2>/dev/null || true
wait "$api_pid" 2>/dev/null || true
