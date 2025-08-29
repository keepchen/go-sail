#!/bin/bash
set -e
mkdir -p bench_results

# CPU profile
go test -bench=. -benchmem -benchtime=1s \
    -cpuprofile bench_results/cpu.out

# Memory profile
go test -bench=. -benchmem -benchtime=1s \
    -memprofile bench_results/mem.out

# Trace
go test -bench=. -benchmem -benchtime=1s \
    -trace bench_results/trace.out

echo "✅ Benchmark done. Results saved in bench_results/"
echo "Run 'go tool pprof bench_results/cpu.out' for CPU profile"
echo "Run 'go tool pprof bench_results/mem.out' for Memory profile"
echo "Run 'go tool trace bench_results/trace.out' for Trace analysis"
