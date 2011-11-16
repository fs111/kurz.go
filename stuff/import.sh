#!/bin/bash

HOW_MANY=100000

for site in $(bzcat top-1m.csv.bz2 | head -n $HOW_MANY | awk -F, '{print $2}'); do curl \
 -d"url=$site" http://localhost:9999/shorten/ ; done
