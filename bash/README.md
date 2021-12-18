# Sorting Exercise

Upon reading the assignment, my initial idea was; this can be done with just a few basic shell commands. Therefore
I provided two solutions; a bash script and a Java program. For details see below.

Which way to go depends here on the use-case (which is not stated in the problem description). As a simple util, the
shell script would likely be sufficient. Is it part of a bigger application / service, then of course the Java version
is more appropriate.


## Shell Script

The shell script is just some syntactic sugar around the three commands `sort`, `head` and `awk`. It provides some input
validation and opportunities for overriding configuration such as default number of results to be returned.

As stated in the problem description:

> Your solution should take into account extremely large files.

The `sort` command is not limited by memory. Instead, it creates temporary files for intermediate state. Since disk
space is 'cheap' and in general provides less of a restriction than RAM, for sake of simplicity I'm assuming this poses
no issue. The Java program does not have this limitation though.

### Usage

Run the script with:

```
$ ./top-hits.sh <file-name>
```

To see all options, run:

```
$ ./top-hits.sh -h
```

### Testing

The shell script was verified with [shellcheck](https://github.com/koalaman/shellcheck). Other testing done manually.
The `sort` command will by itself handle most invalid input, such as missing numeric values or non-numeric data.
