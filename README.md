## What is this?
MARC command line (marcli) is a tool to parse [MARC](https://www.loc.gov/marc/umb/um01to06.html) files from the command line.

The goal `marcli` is to allow a user to parse large MARC files in either MARC binary or MARC XML form and provide basic functionality to find specific records within a file.


## Installation
On the Mac the easiest way to install `marcli` is via Homebrew:

```
brew install marcli
```

Or by downloading the binary for your OS from the [releases tab](https://github.com/hectorcorrea/marcli/releases) and marking the downloaded file as an executable:

```
curl -L https://github.com/hectorcorrea/marcli/releases/download/1.0.1/marcli > marcli
chmod u+x marcli
```

Once installed you can just run it via:

```
./marcli -file yourfile.mrc
```


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

Extracts MARC records on file that contain the string "wildlife" but outputs only fields "LDR,001,040,245a,650" for each record.

```
./marcli -file data/test_10.mrc -match wildlife -fields LDR,010,040,245a,650
```

LDR means the leader of the MARC record.

A letter (or letters) after the field tag indicates to output only those subfields. For example "907xz" means output subfield "x" and "z" in field "907".

You can also use the `exclude` option to indicate fields to exclude from the output (notice that only full subfields are supported here, e.g. 970 is accepted but not 970a)

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


## Bugs, feedback, ideas?
If you find an issue parsing MARC files with `marcli` feel free to [submit an issue](https://github.com/hectorcorrea/marcli/issues) with details of the error, and if possible a sample file or contact me by email at hector@hectorcorrea.com


## More information
* Code4Lib 2021 lightning talk ([slides](https://docs.google.com/presentation/d/1hkLx5zNZCXal20vzP3Jg_nZy03qCsLishHeVTecnsY0/edit?usp=sharing) and [video](https://youtu.be/jLg7XreYS4M?t=186))
* Understanding MARC: https://www.loc.gov/marc/umb/
