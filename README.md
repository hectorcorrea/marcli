## What is this?
This is a toy project that I am working on to help me deal with MARC files and understand better how MARC files are stored on disk. This is also one of my first times using Go's file IO methods.


## MARC information
Understanding MARC: https://www.loc.gov/marc/umb/


## Sample data
Files under `./data/` are small MARC files that I use for testing.

* test_10.mrc has 10 MARC records
* test_1a.mrc is the first record of test_10.mrc
* test_1b.mrc is the second record of test_10.mrc


## Getting staertedTesting
Download the code and play with it:

```
git clone https://github.com/hectorcorrea/marcli.git
cd marcli
go build
./marcli data/test_1a.mrc  
```
