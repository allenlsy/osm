#!/bin/bash

set -aueo pipefail

pkgs=(
   pkg/injector
   pkg/certificate/providers/tresor
   pkg/envoy/ads
)

date_s=$(date '+%Y-%m-%d-%H-%M-%S')
bm_dir=bm_profiles/$date_s

output_file=$bm_dir/bm_output_$date_s

mkdir -p $bm_dir

for pkg in "${pkgs[@]}"; do
   echo "Benchmarking $pkg"
   name=$(echo ${pkg//\//_})-$date_s
   go test -benchmem -run=^$ -bench ^Benchmark github.com/openservicemesh/osm/$pkg -cpuprofile $bm_dir/$name.cpu.pprof -memprofile $bm_dir/$name.mem.pprof -trace $bm_dir/$name.trace >> $output_file
   go tool pprof -png -output $bm_dir/$name.cpu.png $bm_dir/$name.cpu.pprof
   go tool pprof -png -output $bm_dir/$name.mem.png $bm_dir/$name.mem.pprof
   go tool trace -png -output $bm_dir/$name.trace.png $bm_dir/$name.trace
done

rm bm_profiles/latest || true
ln -s $date_s bm_profiles/latest