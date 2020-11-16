## What is this?
MARC command line (marcli) is a toy project that I am working on to help me deal with MARC files and understand better how [MARC](https://www.loc.gov/marc/umb/um01to06.html) files are stored on disk.

The basic idea is to create a program that I can run on our Linux servers to parse our very large MARC files. The goal is to be able to copy a single file to our servers to parse these MARC files, find records that match certain criteria and export them for further review.

The code is written in Go. I've found Go a very interesting and powerful programming language. Go's standard library provides most of the functionality that I need to parse MARC files. Go also supports creating binaries that can be deployed to Mac/Linux/Windows as single executable files which I love because I can deploy my executable to our Linux servers without having to do a complicated installation (i.e. no JVM needed, or Ruby bundle).


## Sample of usage

Output MARC data to the console in a line delimited format:
```
./marcli -file data/test_1a.mrc
```

If the file extension is `.xml` the file is expected to be a MARC XML file, otherwise MARC binary is assumed.

Extract MARC records on file that contain the string "wildlife"
```
./marcli -file data/test_10.mrc -match wildlife
```

Extracts MARC records on file that contain the string "wildlife" but outputs
only fields "LDR,001,040,245a,650" for each record.

```
./marcli -file data/test_10.mrc -match wildlife -fields LDR,010,040,245a,650
```

LDR means the leader of the MARC record.

A letter (or letters) after the field tag indicates to output only those
subfields. For example "907xz" means output subfield "x" and "z" in
field "907".

You can also filter based on the presence of certain fields in the MARC record (regardless of their value), for example the following will only output records that have a MARC 110 field:

```
./marcli -file data/test_10.mrc -hasFields 110
```

The program supports a `format` parameter to output to other formats other than MARC line delimited (MRK) such as MARC XML, JSON, or MARC binary. Notice that not all the features are available in all the formats yet.

You can also pass `start` and `count` parameters to output only a range of MARC records.


## Sample data
Files under `./data/` are small MARC files that I use for testing.

* test_10.mrc has 10 MARC records (MARC binary)
* test_1a.mrc is the first record of test_10.mrc  (MARC binary)
* test_1b.mrc is the second record of test_10.mrc  (MARC binary)
* test_10.xml same as test_10.mrc but in MARC XML.


## Getting started with the code
Download the code and play with it:

```
cd ~/src
git clone https://github.com/hectorcorrea/marcli.git
cd marcli/cmd/marcli
go build
./marcli -file ~/src/marcli/data/test_1a.mrc
```

## Code Structure

* `cmd/marcli` contains the code for the command line interface.
* `pkg/marc` contains the code to parse MARC files.

## Getting started (without the source code)
If you don't care about the source code, you can download the binary file appropriated for your operating system from the [releases tab](https://github.com/hectorcorrea/marcli/releases) and run it.

The basic syntax is:

```
./marcli -file yourfile.mrc
```


## Warning
I've only tested this program with a few internal MARC files and, although they are pretty big (400MB), I have no idea how well this program works with MARC files in the wild. Please keep this in mind if you download and play with it. And feel free to let me know if you run into any issues.


## MARC information
Understanding MARC: https://www.loc.gov/marc/umb/
