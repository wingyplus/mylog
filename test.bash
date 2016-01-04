#!/bin/bash

echo "### Run test -- No color"
MYLOGCOLOR_DISABLED=1 godep go test || exit 1

echo "### Run test -- Color"
godep go test -run=TestLog_Colors || exit 1
