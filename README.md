# histogram

Display a sorted histogram of the frequency of the input lines.

## Usage

```
histogram < sample.txt
```

### Columns

This program does not columnize the output; rather it uses a single
space character as the delimiter between the frequency count and the
input lines. If prettified output is desired, one can pipe output to
another utility to do that.  One such example is
[my columnize utility](https://github.com/karrick/columnize).

```
histogram < sample.txt | columnize
```

### Reverse the Sort for Ascending Order

By default this program sorts the histogram in descending order, from
the most frequent to the least frequent token. When provided the `-r`
command line option, it will reverse the sort and the output in
ascending order.

```
histogram -r < sample.txt | columnize
```

### Selecting a Field

By default this program parses each line into a token and strips
leading and trailing whitespace, creating a histogram assuming the
entire line is the token. When given the `-f N` command line option,
it will select and tokenize the Nth field. For example, `-f 1` creates
a histogram from the first field in the line.

```
histogram -f 2 < sample.txt | columnize
```

### Specifying a Delimiter

By default when this program is given a `-f N` command line option to
specify a particular field to use, it splits each line by
whitespace. When the `-d S` command line option is given, it uses the
provided string as the field delimiter. `S` may be a string of
multiple characters.

```
histogram -f 2 -d : < sample.txt | columnize
```

## Installation

If you don't have the Go programming language installed, then you'll
need to install a copy from
[https://golang.org/dl](https://golang.org/dl).

Once you have Go installed:

```
go get github.com/karrick/histogram
```
