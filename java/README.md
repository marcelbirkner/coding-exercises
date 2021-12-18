# Sorting Exercise

Upon reading the assignment, my initial idea was; this can be done with just a few basic shell commands. Therefore
I provided two solutions; a bash script and a Java program. For details see below.

Which way to go depends here on the use-case (which is not stated in the problem description). As a simple util, the
shell script would likely be sufficient. Is it part of a bigger application / service, then of course the Java version
is more appropriate.

## Java Program

The Java program reads from the input file line by line, checking if there are higher hits than the ones given before.
A `build.gradle` file is provided, but this is only needed for easily running tests. Compiling and running the program
itself can also be done without Gradle. See both instructions below.

Everything was build and tested using JDK 11. JDK 11 is also the minimum because of the used `Scanner` class for reading
the filename from input.

### Compiling and running without Gradle

Compile using `javac`. Full command from inside the `top-hits` directory:

```
javac -d build/classes/java/main/ src/main/java/TopHits.java <input file>
```

Then run using:

```
java -cp build/classes/java/main/ TopHits <input-file>
```


### Compiling and running with Gradle

Compile by running Gradle (from inside the `top-hits` directory):

```
$ ./gradlew build
```

Then execute by:

```
java -jar build/libs/top-hits-1.0.0-SNAPSHOT.jar
```

Potentially provide the number of results to return by adding the arguments `-r <number>`
