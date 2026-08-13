#!/usr/bin/env bash
set -euo pipefail

httpbin_base="${HTTPBIN_BASE:-https://httpbin.org}"

echo "GET /get"
curl --silent --show-error --verbose "${httpbin_base}/get"

echo
echo "POST /post"
curl --silent --show-error -X POST "${httpbin_base}/post" \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","age":30}'

echo
echo "PUT /put"
curl --silent --show-error -X PUT "${httpbin_base}/put" \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","age":31}'

echo
echo "DELETE /delete"
curl --silent --show-error -X DELETE "${httpbin_base}/delete"

echo
echo "HEAD /get"
curl --silent --show-error --head "${httpbin_base}/get"

echo
echo "GET /headers with X-Custom-Header"
curl --silent --show-error -H "X-Custom-Header: hello" "${httpbin_base}/headers"
