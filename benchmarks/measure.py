#!/usr/bin/env python3
"""Run the aggregator and collect child peak RSS and elapsed wall time."""

import argparse
import resource
import subprocess
import sys
import time


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--binary", required=True)
    parser.add_argument("--input", required=True)
    parser.add_argument("--output", required=True)
    parser.add_argument("--stdout", required=True)
    parser.add_argument("--stderr", required=True)
    parser.add_argument("--cpuprofile")
    args = parser.parse_args()

    command = [args.binary, "--input", args.input, "--output", args.output]
    if args.cpuprofile:
        command.extend(["--cpuprofile", args.cpuprofile])
    start = time.perf_counter()
    with open(args.stdout, "wb") as stdout, open(args.stderr, "wb") as stderr:
        process = subprocess.Popen(command, stdout=stdout, stderr=stderr)
        return_code = process.wait()
    usage = resource.getrusage(resource.RUSAGE_CHILDREN)
    elapsed = time.perf_counter() - start
    peak_rss_bytes = int(usage.ru_maxrss)
    if sys.platform != "darwin":
        peak_rss_bytes *= 1024

    print(f"command={' '.join(command)}")
    print(f"wall_seconds={elapsed:.3f}")
    print(f"peak_rss_bytes={peak_rss_bytes}")
    print("peak_rss_source=getrusage(RUSAGE_CHILDREN).ru_maxrss")
    print(f"user_cpu_seconds={usage.ru_utime:.3f}")
    print(f"system_cpu_seconds={usage.ru_stime:.3f}")
    print(f"exit_code={return_code}")
    return return_code


if __name__ == "__main__":
    raise SystemExit(main())
