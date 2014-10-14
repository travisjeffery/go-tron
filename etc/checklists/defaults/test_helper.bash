#!/bin/env sh

echo_url_for_description() {
  local test_description="$1"
  local param="$(echo $test_description | sed -e 's/[^a-zA-Z0-9 ]//g' | sed -e 's/[^a-zA-Z0-9]/-/g' | tr '[:upper:]' '[:lower:]')"
  echo "https://github.com/travisjeffery/tron/wiki/Defaults-Checklist#$param"
}
