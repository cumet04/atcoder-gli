#!/bin/bash

cd $(dirname $0)/..

echo "### Usage"
for cmd in \
  login \
  new \
  lang \
  submit \
  session \
  open \
  show \
;do
  echo ""
  echo "#### ${cmd}"
  echo '```'
  go run main.go help $cmd
  echo '```'
done
