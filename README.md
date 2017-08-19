# goblktest

## Overview
Tests a block device by writing an iterating SHA-512 hash over the entire device, and then reads back the pattern while comparing to the same computation run a second time.

This validates the entire storage pipeline in two important ways:

1. Guarantees that each read and write is addressing exactly one location, and doing so consistently. (eg No subtle device bug that is occasionally putting bytes in the wrong sector, or block scheduler/elevator bug that returns reads in the wrong order)
2. Tests that billions of pseudo-random bit patterns can reliably traverse media, caches, cables, connectors, and all the various components involved in transporting them.

Note that a failure that would be detected by this tool -- but not a more basic testing tool such as `badblocks -w` -- is a wildly pessimistic and unlikely one. However, it may be a failure that would be insanely difficult to detect otherwise, particularly on archival drives where data is mostly written and rarely read.

## Usage

`goblktest --open /dev/sdz --bs 4096`

* `--open` (required)
* * Specifies block device to open for testing.
* `--bs` (optional, but recommended)
* * Specifies block size for each write. Larger values will tend to perform better, but will make failure regions larger.
* * Note that bs must be an even multiple of the hash size (presently SHA-512, 64 bytes) and must evenly divide the medium size
* * (TODO for future: gather failed block regions and re-test at hash size)
* `--seed` (optional)
* * If running the test on the same medium multiple times, use a different seed value to force a unique hash pattern each time.
* `--start` (optional)
* * It's possible to skip some initial blocks on the device, for example to preserve an MBR. 
* `--checkonly` (optional)
* * If the medium has passed a test previously, it's possible to re-test the read portion. For success, use the same `--bs`, `--seed`, and `--start` arguments, if provided before.

## Contributing
Pull requests and issues are welcome.


