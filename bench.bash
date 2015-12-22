#!/bin/bash

echo "### Run bench"
./clean.bash && go test -run=NONE -bench=. -cpu 1,2,4,8
