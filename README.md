## What is this?
MARC command line (marcli) This is a toy project that I am working on to help me deal with MARC files and understand better how [MARC](https://www.loc.gov/marc/umb/um01to06.html) files are stored on disk.

The basic idea is to create a program that I can run on our Linux servers to parse our very large MARC files. The goal is to be able to copy a single file to our servers to parse these MARC files, find records that match certain criteria and export them for further review.

This program is heavily inspired by amazing work that Terry Reese has done on [MarcEdit](http://marcedit.reeset.net/). This program is not a replacement for MarcEdit.

The code is written in Go. I've found Go a very interesting and powerful programming language. Go's standard library provides most of the functionality that I need to parse MARC files. Go also supports creating binaries that can be deployed to Mac/Linux/Windows as single executable files which I love because I can deploy my executable to our Linux servers without having to do a complicated installation (i.e. no JVM needed, or Ruby bundle).


## Sample of usage

Output MARC data to the console in a line delimited format:
```
./marcli -f data/test_1a.mrc
```

Extract MARC records on file that match a given string
```
./marcli -f data/test_10.mrc -x wildlife
```

Extracts MARC records on file that match a given string but outputs only field "650" (and the leader) for each record. [TODO: add support for multiple fields]
```
./marcli -f data/test_10.mrc -x wildlife -o 650
```


## Sample data
Files under `./data/` are small MARC files that I use for testing.

* test_10.mrc has 10 MARC records
* test_1a.mrc is the first record of test_10.mrc
* test_1b.mrc is the second record of test_10.mrc


## Getting started with the code
Download the code and play with it:

```
git clone https://github.com/hectorcorrea/marcli.git
cd marcli
go build
./marcli -f data/test_1a.mrc  
```


## MARC information
Understanding MARC: https://www.loc.gov/marc/umb/
