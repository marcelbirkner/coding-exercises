package main

import (
	"testing"
)

func TestFileProcessingForInitialDataset(t *testing.T) {
	filepath := "testdata/initial-dataset.txt"

	t.Run("get URLs for two largest values from "+filepath, func(t *testing.T) {
		const maxEntries = 2

		fp := FileProcessor{filepath}
		got, err := fp.FindLargestEntriesInFile(maxEntries)
		assertNoError(t, err)

		want := [maxEntries]string{
			"http://api.tech.com/item/122345",
			"http://api.tech.com/item/124345",
		}

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

		want := [maxEntries]string{
			"http://api.tech.com/item/122345",
			"http://api.tech.com/item/124345",
			"http://api.tech.com/item/125345",
			"http://api.tech.com/item/123345",
		}

		for i := 0; i < maxEntries; i++ {
			if got[i] != want[i] {
				t.Errorf("got %s want %s", got[i], want[i])
			}
		}
	})

}

func TestFileProcessingOfLargerTestdata(t *testing.T) {
	filepath := "testdata/100-entries.txt"

	t.Run("get URLs for three largest values from "+filepath, func(t *testing.T) {
		const maxEntries = 3

		fp := FileProcessor{filepath}
		got, err := fp.FindLargestEntriesInFile(maxEntries)
		assertNoError(t, err)

		want := [maxEntries]string{
			"http://api.tech.com/item/760",
			"http://api.tech.com/item/210",
			"http://api.tech.com/item/100",
		}

		for i := 0; i < maxEntries; i++ {
			if got[i] != want[i] {
				t.Errorf("got %s want %s", got[i], want[i])
			}
		}
	})
}

func TestFileProcessingForDifferentWhitespaces(t *testing.T) {
	filepath := "testdata/different-whitespaces.txt"

	t.Run("handle different valid whitespaces for "+filepath, func(t *testing.T) {
		const maxEntries = 2

		fp := FileProcessor{filepath}
		got, err := fp.FindLargestEntriesInFile(maxEntries)
		assertNoError(t, err)

		want := [maxEntries]string{
			"http://api.tech.com/item/10",
			"http://api.tech.com/item/5",
		}

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

func TestFileProcessingOfInvalidLines(t *testing.T) {
	filepath := "testdata/empty-lines.txt"

	t.Run("ignore invalid lines in "+filepath, func(t *testing.T) {
		const maxEntries = 2
		fp := FileProcessor{filepath}
		got, err := fp.FindLargestEntriesInFile(maxEntries)
		assertNoError(t, err)

		want := [maxEntries]string{
			"http://api.tech.com/item/99",
			"http://api.tech.com/item/10",
		}

		for i := 0; i < maxEntries; i++ {
			if got[i] != want[i] {
				t.Errorf("got %s want %s", got[i], want[i])
			}
		}
	})
}

func TestFileProcessingOfInvalidLongValue(t *testing.T) {
	filepath := "testdata/invalid-long-value.txt"

	t.Run("expect parsing error "+filepath, func(t *testing.T) {
		const maxEntries = 2
		fp := FileProcessor{filepath}
		_, err := fp.FindLargestEntriesInFile(maxEntries)
		assertError(t, err, ParsingInvalidLongErr)
	})
}

func BenchmarkFindLargestEntriesInFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		const maxEntries = 2
		filepath := "testdata/initial-dataset.txt"
		fp := FileProcessor{filepath}

		fp.FindLargestEntriesInFile(maxEntries)
	}
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	t.Helper()
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

func assertNoError(t testing.TB, got error) {
	t.Helper()
	if got != nil {
		t.Fatal("got an error but didn't want one")
	}
}

func assertError(t testing.TB, got, want error) {
	t.Helper()

	if got != want {
		t.Errorf("got error %q want %q", got, want)
	}
}