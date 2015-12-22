#!/bin/bash

go generate

echo "### Run test -- No color"
MYLOGCOLOR_DISABLED=1 go test || exit 1

echo "### Run test -- Color"
go test -run=TestLog_DebugColor || exit 1

echo "### Run bench"
./clean.bash && go test -run=NONE -bench=. -cpu 1,2,4,8
