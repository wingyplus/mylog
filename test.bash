#!/bin/bash

go generate

echo "### Run test -- No color"
MYLOGCOLOR_DISABLED=1 go test || exit 1

echo "### Run test -- Color"
go test -run=TestLog_Colors || exit 1
