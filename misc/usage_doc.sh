#!/bin/bash

cd $(dirname $0)/..

EXEC=./usage_doc_tmp_build
go build -o $EXEC acg/main.go

echo "### Usage"
for cmd in \
  login \
  new \
  config \
  "config lang" \
  submit \
  session \
  open \
  show \
;do
  echo ""
  echo "#### ${cmd}"
  echo '```'
  $EXEC help $cmd
  echo '```'
done

rm $EXEC
