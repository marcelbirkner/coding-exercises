package main

import (
	"testing"
)

func TestFileProcessing(t *testing.T) {
	filepath := "testdata/testdataset-1.txt"

	t.Run("get URLs for two largest values from "+filepath, func(t *testing.T) {
		const maxEntries = 2

		fp := FileProcessor{filepath}
		got, err := fp.FindLargestEntriesInFile(maxEntries)
		assertNoError(t, err)

		want := [maxEntries]string{"http://api.tech.com/item/122345", "http://api.tech.com/item/124345"}

		for i := 0; i < maxEntries; i++ {
			if got[i] != want[i] {
				t.Errorf("got %s want %s", got[i], want[i])
			}
		}
	})

	t.Run("get URLs for four largest values from "+filepath, func(t *testing.T) {
		const maxEntries = 4

		fp := FileProcessor{filepath}
		got, err := fp.FindLargestEntriesInFile(maxEntries)
		assertNoError(t, err)

		want := [maxEntries]string{"http://api.tech.com/item/122345", "http://api.tech.com/item/124345", "http://api.tech.com/item/125345", "http://api.tech.com/item/123345"}

		for i := 0; i < maxEntries; i++ {
			if got[i] != want[i] {
				t.Errorf("got %s want %s", got[i], want[i])
			}
		}
	})

}

func TestFileProcessingOfTestdataset2(t *testing.T) {
	filepath := "testdata/testdataset-2.txt"

	t.Run("get URLs for three largest values from "+filepath, func(t *testing.T) {
		const maxEntries = 3

		fp := FileProcessor{filepath}
		got, err := fp.FindLargestEntriesInFile(maxEntries)
		assertNoError(t, err)

		want := [maxEntries]string{"http://api.tech.com/item/760", "http://api.tech.com/item/210", "http://api.tech.com/item/100"}

		for i := 0; i < maxEntries; i++ {
			if got[i] != want[i] {
				t.Errorf("got %s want %s", got[i], want[i])
			}
		}
	})
}

func TestFileProcessingOfTestdataset3(t *testing.T) {
	filepath := "testdata/testdataset-3.txt"

	t.Run("handle different whitespaces for "+filepath, func(t *testing.T) {
		const maxEntries = 2

		fp := FileProcessor{filepath}
		got, err := fp.FindLargestEntriesInFile(maxEntries)
		assertNoError(t, err)

		want := [maxEntries]string{"http://api.tech.com/item/10", "http://api.tech.com/item/5"}

		for i := 0; i < maxEntries; i++ {
			if got[i] != want[i] {
				t.Errorf("got %s want %s", got[i], want[i])
			}
		}
	})

}

func TestFileProcessingOfEmptyFile(t *testing.T) {
	filepath := "testdata/empty-file.txt"

	t.Run("expect empty list for empty file "+filepath, func(t *testing.T) {
		const maxEntries = 10
		fp := FileProcessor{filepath}
		got, err := fp.FindLargestEntriesInFile(maxEntries)
		assertNoError(t, err)
		assertEqual(t, len(got), 0)
	})
}

func TestFileProcessingOfInvalidInputFile(t *testing.T) {
	filepath := "testdata/testdataset-4.txt"

	t.Run("ignore invalid lines in "+filepath, func(t *testing.T) {
		const maxEntries = 2
		fp := FileProcessor{filepath}
		got, err := fp.FindLargestEntriesInFile(maxEntries)
		assertNoError(t, err)

		want := [maxEntries]string{"http://api.tech.com/item/5", "http://api.tech.com/item/3"}

		for i := 0; i < maxEntries; i++ {
			if got[i] != want[i] {
				t.Errorf("got %s want %s", got[i], want[i])
			}
		}
	})
}

func BenchmarkFindLargestEntriesInFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		const maxEntries = 2
		filepath := "testdata/testdataset-1.txt"
		fp := FileProcessor{filepath}

		fp.FindLargestEntriesInFile(maxEntries)
	}
}

func assertNoError(t testing.TB, got error) {
	t.Helper()
	if got != nil {
		t.Fatal("got an error but didn't want one")
	}
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	t.Helper()
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}
