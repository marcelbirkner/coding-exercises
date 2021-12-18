package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	DefaultResultSetSize = 10
)

// Errors returned by FileProcessor.
const (
	FileNotFoundErr       = FileProcessorErr("Unable to find file")
	CannotOpenFileErr     = FileProcessorErr("Cannot open file")
	ParsingInvalidLongErr = FileProcessorErr("Invalid syntax. Expecting a long value")
)

type FileProcessorErr string

func (e FileProcessorErr) Error() string {
	return string(e)
}

type FileProcessor struct {
	filepath string
}

// Main program
//
// User has to provide absolute filepath via stdin.
//
// File is expected to have the following format:
// '<url value> <whitespace> <long value>'
//
// This file is parsed line-by-line to be able to handle large files.
// URLs with the highest count are returned in descending order.
// By default, the result is limited to 10 URLs.
//
func main() {

	var filepath string
	if len(os.Args) == 3 {
		flag.StringVar(&filepath, "filepath", "", "Full path to input file")
		flag.Parse()
	} else {
		fmt.Print("Please enter filename: ")
		fmt.Scanln(&filepath)
	}

	// fail fast: validate user input
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		log.Fatalf(FileNotFoundErr.Error())
	}

	fp := FileProcessor{filepath}
	urls, err := fp.FindLargestEntriesInFile(DefaultResultSetSize)
	if err != nil {
		log.Fatal(err)
	}

	// print urls on stdout
	for _, url := range urls {
		fmt.Fprintln(os.Stdout, url)
	}
}

// FindLargestEntriesInFile Idea behind implementation
//
// The input file is processed line-by-line
// Each line is parsed and split into "url" as string and "count" as int64
// We keep a map of the urls with the most counts (default: maximum number of element in map is 10). count is the key, and the url is the value
// We also keep a sorted list of the largest counts (default: maximum number of elements in list is 10)
// Once all lines have been processed we take use the sorted list of keys and create a list of all URLs with the highest count
func (fp *FileProcessor) FindLargestEntriesInFile(resultSetSize int) ([]string, error) {

	// Stores count as key and url as value. maximum number of items stored equals resultSetSize
	var resultMap = map[int64]string{}

	// Store highest counts in descending order
	var countListDescendingOrder []int64

	f, err := os.Open(fp.filepath)
	if err != nil {
		return nil, CannotOpenFileErr
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {

		// parse line delimited by whitespace
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) != 2 {
			continue
		}
		url := fields[0]
		count, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			return nil, ParsingInvalidLongErr
		}

		// fill list and map until resultSetSize is reached
		if len(countListDescendingOrder) < resultSetSize {
			// add new entry to result set
			countListDescendingOrder = fp.storeItemInResultSet(countListDescendingOrder, count, resultMap, url)
		} else {
			// check if smallest item in countListDescendingOrder is smaller than new count
			smallestItemInCountList := countListDescendingOrder[resultSetSize-1]
			if smallestItemInCountList < count {
				// remove the smallest item from list and map
				delete(resultMap, smallestItemInCountList)
				countListDescendingOrder = countListDescendingOrder[:len(countListDescendingOrder)-1]

				// add new entry to result set
				countListDescendingOrder = fp.storeItemInResultSet(countListDescendingOrder, count, resultMap, url)
			}
		}
	}

	// return if scanner encounters an error
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// get list of sorted keys from resultMap
	keys := make([]int64, 0, len(resultMap))
	for k := range resultMap {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] > keys[j] })

	// create result list with urls, based on sorted keys
	resultList := make([]string, 0, resultSetSize)
	for _, k := range keys {
		resultList = append(resultList, resultMap[k])
	}

	return resultList, nil
}

// Add item to countListDescendingOrder and resultMap
// Sort countListDescendingOrder afterwards in descending order
func (fp *FileProcessor) storeItemInResultSet(countListDescendingOrder []int64, count int64, resultMap map[int64]string, url string) []int64 {
	countListDescendingOrder = append(countListDescendingOrder, count)
	resultMap[count] = url

	// sort countListDescendingOrder
	sort.Slice(countListDescendingOrder, func(i, j int) bool { return countListDescendingOrder[i] > countListDescendingOrder[j] })
	return countListDescendingOrder
}
