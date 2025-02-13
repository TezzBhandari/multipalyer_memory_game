#!/bin/bash

# Path to your Go executable
GO_PROGRAM="go run cmd/client/main.go"

# Number of instances to run
NUM_INSTANCES=100

# Loop to start multiple instances
for ((i=1; i<=NUM_INSTANCES; i++))
do
    echo "Starting instance $i..."
    $GO_PROGRAM &  # Run in background
    sleep 0.2      # Optional delay to avoid CPU spikes
done

echo "All $NUM_INSTANCES instances started."

