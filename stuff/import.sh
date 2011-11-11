#!/bin/bash

HOW_MANY=100000

for i in $(bzcat top-1m.csv.bz2 | head -n $HOW_MANY | awk -F, '{print $2}'); do curl \
    http://localhost:9999/store/$i 2>&1 > /dev/null; done
