Benchmark using:
hyperfine --prepare 'sudo sync; echo 3 | sudo tee /proc/sys/vm/drop_caches' \
--warmup 10 \
--min-runs 100 \
"sudo $RUNTIME_EXECUTABLE create -b /tmp/roci a && sudo $RUNTIME_EXECUTABLE start a && sudo $RUNTIME_EXECUTABLE delete -f a"


## runc
- language: Go
- specification: Full OCI
```
  Time (mean ± σ):     102.4 ms ±  11.6 ms    [User: 10.6 ms, System: 12.6 ms]
  Range (min … max):    77.9 ms … 139.0 ms    100 runs
```

## crun
- language: C
- specification: Full OCI
```
  Time (mean ± σ):      48.8 ms ±   3.1 ms    [User: 10.4 ms, System: 12.7 ms]
  Range (min … max):    44.9 ms …  67.2 ms    100 runs
```

## roci
- language: Go
- specification: Minimal OCI
```
  Time (mean ± σ):      61.3 ms ±   1.0 ms    [User: 10.7 ms, System: 12.8 ms]
  Range (min … max):    57.3 ms …  65.3 ms    100 runs
```