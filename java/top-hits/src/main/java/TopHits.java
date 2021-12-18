import java.io.BufferedReader;
import java.io.File;
import java.io.FileReader;
import java.io.IOException;
import java.util.*;
import java.util.stream.Collectors;

public class TopHits {
    private static final String USAGE =
            "<command> [-r RESULTS] FILE\n" +
                    "where:\n" +
                    "\t-r  Number of results to print (Default 10 lines)";

    public static void main(String[] args) {
        // Hold a copy to parse arguments without modifying the original array
        String[] processedArgs = args;

        int numberOfResults = 10;
        String file = null;

        //
        // Do some basic assertions on the arguments here, so the actual function can assume parameters are valid
        //

        // Verify if the '-r' argument was provided and if it's valid.
        if (processedArgs[0].startsWith("-")) {
            try {
                numberOfResults = getNumberOfResultsArgument(processedArgs);
                // Remove processed arguments
                processedArgs = Arrays.copyOfRange(processedArgs, 2, processedArgs.length);
            } catch (IllegalArgumentException e) {
                System.err.println("Illegal argument(s); " + e.getMessage() + ". Invoke as:\n" + USAGE);
                System.exit(1);
            }
        }

        // We allow providing the file in two ways; either directly as parameter, otherwise request the path on 'stdin'
        if (processedArgs.length >= 1) {
            file = processedArgs[0];
        } else {
            Scanner scanner = new Scanner(System.in);
            //  prompt for the filename
            System.out.print("Enter the file for processing (submit with <enter>): ");
            // Get the actual value
            file = scanner.next();
        }

        if (!fileExistsAndIsReadable(file)) {
            System.err.println("No valid file provided, '" + file + "' does not exist or non-readable. Invoke as:\n" + USAGE);
            System.exit(1);
        }

        //
        // Retrieve the actual top results from the specified file. Then print the results.
        //
        List<String> results = getTopResults(new File(file), numberOfResults);
        System.out.println("Results:");
        for (String r : results) {
            System.out.println(r);
        }
    }

    protected static int getNumberOfResultsArgument(String... args) throws IllegalArgumentException {
        if (args[0].equals("-r")) {
            if (args.length >= 2) {
                try {
                    return Integer.parseInt(args[1]);
                } catch (Exception e) {
                    throw new IllegalArgumentException("value '" + args[1] + "' cannot be parsed as numeric value");
                }
            } else {
                throw new IllegalArgumentException("no value for '-r' argument");
            }
        } else {
            throw new IllegalArgumentException("unknown argument");
        }
    }

    protected static boolean fileExistsAndIsReadable(String filename) {
        if (filename == null) {
            return false;
        }
        final File file = new File(filename);
        return file.isFile() && file.canRead();
    }

    /**
     * Get the 'x' number of top results from the file. The File is expected to have lines in the format;
     *   <URL> <space> <long>
     * Invalid lines will be silently skipped.
     * @param file The file with lines to be processed
     * @param numberOfResults Number of results to be returned
     * @return List of sorted hits, the highest first. Or empty if the input file contained no valid lines
     */
    protected static List<String> getTopResults(File file, int numberOfResults) {
        // We cannot use a Map with e.g. the 'value' as key. While that would be faster for removing values, it won't allow to
        // contain duplicates
        List<ParsedLine> topResults = new ArrayList<>(numberOfResults);
        // Use a reference to the smallest result contained in the set, so if new results are smaller we know it doesn't need updating.
        ParsedLine smallestResult = null;

        // Assume default buffer size to be sufficient as we have relatively short lines
        try (BufferedReader br = new BufferedReader(new FileReader(file))) {
            String line;
            while ((line = br.readLine()) != null) {

                ParsedLine parsedLine = parseLine(line);
                if (parsedLine == null) {
                    continue;
                }

                // If the number of results isn't reached yet, fill the map
                if (topResults.size() < numberOfResults) {
                    topResults.add(parsedLine);
                    if (smallestResult == null || parsedLine.value < smallestResult.value) {
                        smallestResult = parsedLine;
                    }
                } else if (parsedLine.value > smallestResult.value) {
                    // We have x results already, but this value is higher than the previous ones and thus should go in the result
                    topResults.remove(smallestResult);
                    topResults.add(parsedLine);
                    // Update the "smallestResult" so its correct again
                    smallestResult = topResults.get(0);
                    for (ParsedLine l : topResults) {
                        if (l.value < smallestResult.value) {
                            smallestResult = l;
                        }
                    }
                } else {
                    // Nothing to do. Branch provided for clarity.
                    // The processed line has a value lower than all existing results and so can be skipped
                }
            }
        } catch (IOException e) {
            System.err.println("Unexpected error while reading from input file '" + file + "'");
            return Collections.emptyList();
        }

        // Get only the URLs but sort them before returning
        return topResults.stream()
                .sorted(Comparator.comparingLong(ParsedLine::getValue).reversed())
                .map(r -> r.url)
                .collect(Collectors.toList());
    }

    /**
     * Parses a single line in the format {@code <url> - <value>}. If the line is invalid, returns NULL
     * @param line Input line in valid format
     * @return {@link ParsedLine} instance or NULL if the input was invalid
     */
    protected static ParsedLine parseLine(String line) {
        if (line == null || line.length() == 0) {
            return null;
        }

        // This will handle ONLY spaces, but is internally optimized so that for simple characters it doesn't do a regex match.
        // If all whitespaces need to be covered, should use "Patter.compile("\\s").split(str)"
        final String[] split = line.trim().split(" ");

        if (split.length < 2) {
            return null;
        }

        String url = split[0];
        String value = null;
        // Because multiple spaces might lead to empty strings in the result, filter out intermediate empty results
        if (split.length > 2) {
            for (int i = 1; i < split.length; i++) {
                if (split[i] != null && split[i].length() > 0) {
                    value = split[i];
                    break;
                }
            }
        } else {
            value = split[1];
        }

        try {
            return new ParsedLine(url, Long.parseLong(value));
        } catch (Exception e) {
            // handle both the invalid (non-numeric) value or null if for some reason not found.
            return null;
        }
    }

    /**
     * Class to hold the combination of URL - value after parsing. So we can easily reference and sort the results
     */
    protected static class ParsedLine {
        private final String url;
        private final long value;

        protected ParsedLine(String url, long value) {
            this.url = url;
            this.value = value;
        }

        // Getter just needed for the stream sorting.
        public long getValue() {
            return value;
        }

        @Override
        public boolean equals(Object o) {
            if (this == o) return true;
            if (o == null || getClass() != o.getClass()) return false;

            ParsedLine that = (ParsedLine) o;

            if (value != that.value) return false;
            return url != null ? url.equals(that.url) : that.url == null;
        }

        @Override
        public int hashCode() {
            int result = url != null ? url.hashCode() : 0;
            result = 31 * result + (int) (value ^ (value >>> 32));
            return result;
        }
    }

}
