#!/usr/bin/env python3
"""Report platform information available without elevated privileges."""

import os
import platform


print(f"processor={platform.processor() or 'unknown'}")
print(f"logical_cpus={os.cpu_count()}")
try:
    physical_memory = os.sysconf("SC_PHYS_PAGES") * os.sysconf("SC_PAGE_SIZE")
except (AttributeError, OSError, ValueError):
    physical_memory = "unknown"
print(f"physical_memory_bytes={physical_memory}")
