#!/bin/bash

# MAP APPLICATIONS
declare -A apps=(["api"]="api" ["driver"]="cli" ["worker"]="worker")

for app in "${!apps[@]}"; do
    echo "Building - ${apps[$app]}"
    go build -o output/bin/$app ./cmd/${apps[$app]}
done
