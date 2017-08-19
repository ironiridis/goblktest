# goblktest

## Overview
Tests a block device by writing an iterating SHA-512 hash over the entire device, and then reads back the pattern while comparing to the same computation run a second time.

This validates the entire storage pipeline in two important ways:

1. Guarantees that each read and write is addressing exactly one location, and doing so consistently. (eg No subtle device bug that is occasionally putting bytes in the wrong sector, or block scheduler/elevator bug that returns reads in the wrong order)
2. Tests that billions of pseudo-random bit patterns can reliably traverse media, caches, cables, connectors, and all the various components involved in transporting them.

Note that a failure that would be detected by this tool -- but not a more basic testing tool such as `badblocks -w` -- is a wildly pessimistic and unlikely one. However, it may be a failure that would be insanely difficult to detect otherwise, particularly on archival drives where data is mostly written and rarely read.

## Usage


## Contributing
Pull requests and issues are welcome.


