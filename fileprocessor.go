package main

import (
	"bufio"
	"errors"
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
	FileNotFoundErr   = FileProcessorErr("Unable to find file")
	CannotOpenFileErr = FileProcessorErr("Cannot open file")
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
// URLs with the highest count are returned in decending order.
// By default the result is limited to 10 URLs.
//
func main() {
	fmt.Print("Please enter filename: ")
	var filepath string
	fmt.Scanln(&filepath)

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

// Idea behind implementation
//
// The input file is processed line-by-line
// Each line is parsed and split into "url" as string and "count" as int64
// We keep a map of the urls with the most counts (default: maximum number of element in map is 10). count is the key, and the url is the value
// We also keep a sorted list of the largest counts (default: maximum number of elements in list is 10)
// Once all lines have been processed we take use the sorted list of keys and create a list of all URLs with the highest count
func (fp *FileProcessor) FindLargestEntriesInFile(resultSetSize int) ([]string, error) {

	// stores count as key and url as value. maximum number of items stored equals resultSetSize
	var resultMap = map[int64]string{}

	// stores counts
	countList := []int64{}

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
		url := fields[0]
		count, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			log.Fatal(err.Error())
		}

		// fill list and map until resultSetSize is reached
		if len(countList) < resultSetSize {
			countList = append(countList, count)
			resultMap[count] = url
			sort.Slice(countList, func(i, j int) bool { return countList[i] > countList[j] })
		} else {
			// check if smallest item in countList is smaller than new count
			smallestItemInCountList := countList[resultSetSize-1]
			if smallestItemInCountList < count {
				// remove smallest item from list and map
				delete(resultMap, smallestItemInCountList)
				countList = countList[:len(countList)-1]

				// add new item to list and map
				countList = append(countList, count)
				resultMap[count] = url

				// sort countList
				sort.Slice(countList, func(i, j int) bool { return countList[i] > countList[j] })
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
