#!/usr/bin/env bash
set -euo pipefail

httpbin_base="${HTTPBIN_BASE:-https://httpbin.org}"
http_timeout="${HTTP_TIMEOUT:-15}"
curl_args=(--silent --show-error --max-time "$http_timeout")

echo "GET /get"
curl "${curl_args[@]}" --verbose "${httpbin_base}/get"

echo
echo "POST /post"
curl "${curl_args[@]}" -X POST "${httpbin_base}/post" \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","age":30}'

echo
echo "PUT /put"
curl "${curl_args[@]}" -X PUT "${httpbin_base}/put" \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","age":31}'

echo
echo "DELETE /delete"
curl "${curl_args[@]}" -X DELETE "${httpbin_base}/delete"

echo
echo "HEAD /get"
curl "${curl_args[@]}" --head "${httpbin_base}/get"

echo
echo "GET /headers with X-Custom-Header"
curl "${curl_args[@]}" -H "X-Custom-Header: hello" "${httpbin_base}/headers"
