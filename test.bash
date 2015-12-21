#!/bin/bash

./clean.bash && go test -bench . -cpu 1,2,4,8
