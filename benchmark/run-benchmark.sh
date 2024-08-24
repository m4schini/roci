#!/bin/bash

# Set the RUNTIME_EXECUTABLE variable to the first argument or default to "roci"
RUNTIME_EXECUTABLE="${1:-roci}"

# Run the hyperfine command with the RUNTIME_EXECUTABLE variable
hyperfine --prepare 'sudo sync; echo 3 | sudo tee /proc/sys/vm/drop_caches' \
          --warmup 10 \
          --min-runs 100 \
          --show-output \
          "sudo $RUNTIME_EXECUTABLE create -b /tmp/roci a && sudo $RUNTIME_EXECUTABLE start a && sudo $RUNTIME_EXECUTABLE delete -f a"
