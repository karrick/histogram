# histogram

Display a histogram of the frequency of keys from input lines.

## Usage

This program reads from the files specified on the command line after
all the flags have been processed, or will read from standard input
when no files are specified. The following two invocations will have
the same effect:

    $ histogram < sample.txt
    $ histogram sample.txt

### Folding matching keys

By default this program only aggregates the count of keys when they
are adjacent to each other line after line. So a run of 10 keys "abc",
followed by a single key "def", followed by another "abc" will result
in the following output:

```
$ histogram
Value Count
  abc     2 *******************************************************************
  def     1 *********************************
  abc     1 *********************************
```

When provided the `--fold` flag, this program folds all matching keys
such that their key value is equal to the aggregate value of the
count:

```
$ histogram --fold
Value Count
  abc     3 *******************************************************************
  def     1 **********************
```

### Selecting a Field

By default this program parses each line into a token and strips
leading and trailing whitespace, creating a histogram assuming the
entire line is the token. When given the `-f N` command line option,
it will select and tokenize the Nth field. For example, `-f 1` creates
a histogram from the first field in the line.

    $ histogram --field 2

### Specifying a Delimiter

By default when this program is given a `-f N` command line option to
specify a particular field to use, it splits each line by
whitespace. When the `-d S` command line option is given, it uses the
provided string as the field delimiter. `S` may be a string of
multiple characters.

    $ histogram --field 2 --delimiter :

### Show Percentage

By default this program shows three columns of output. The value from
the input text, followed by the number of times that key was found in
a series, followed by a row of asterisk characters to show relative
number of times that key was found compared to the other keys.

When given the `-p, --percentage` flag, this program also shows a
numeric percentage after the count column and before the histogram
column.

```
$ histogram --percentage sample.txt
Value Count Percent
  abc     2   50.00 ************************************************************
  def     1   25.00 ******************************
  abc     1   25.00 ******************************
```

### Sort Ascending

By default this program displays the output in the order the keys were
encountered. When provided the `--ascending` flag, this program sorts
the output such that the keys with the lowest counts are displayed
first, and all successive lines will show a key with a count matching
or greater to the previous line's count.

```
$ histogram --ascending sample.txt
Value Count
  def     1 *********************************
  abc     1 *********************************
  abc     2 *******************************************************************
```

### Sort Descending

By default this program displays the output in the order the keys were
encountered. When provided the `--descending` flag, this program sorts
the output such that the keys with the largest counts are displayed
first, and all successive lines will show a key with a count matching
or less than the previous line's count.

```
$ histogram --descending sample.txt
Value Count
  abc     2 *******************************************************************
  def     1 *********************************
  abc     1 *********************************
```

### Histogram Width

By default this program scales the output such that the longest row
will consume no more than 80 characters. When provided the `-w,
--width` flag, this program scales the output such that the longest
row will consume no more than the specified number of characters. If
the specified number of characters is too narrow, this program exits
with an error.

```
$ histogram --width 20 sample.txt
Value Count
  abc     2 *******
  def     1 ***
  abc     1 ***
```

## Installation

If you don't have the Go programming language installed, then you'll
need to install a copy from
[https://golang.org/dl](https://golang.org/dl).

Once you have Go installed:

```
$ go get github.com/karrick/histogram
```
