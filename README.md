# Run main program

```
go run fileprocessor.go
```

# Run Tests 

```
/usr/local/bin/go test
PASS
ok  	fileprocessing	0.094s
```

With code coverage

```
Running tool: /usr/local/bin/go test -timeout 30s -coverprofile=/var/folders/4w/frj676b93pv__6tl9b14hmn40000gn/T/vscode-govAHXMf/go-code-cover fileprocessing

ok  	fileprocessing	0.468s	coverage: 70.0% of statements
```

# Run benchmark

```
Running tool: /usr/local/bin/go test -benchmem -run=^$ -bench ^BenchmarkFindLargestEntriesInFile$ fileprocessing

goos: darwin
goarch: amd64
pkg: fileprocessing
cpu: Intel(R) Core(TM) i7-8850H CPU @ 2.60GHz
BenchmarkFindLargestEntriesInFile-12    	   34191	     35788 ns/op	    4952 B/op	      27 allocs/op
PASS
ok  	fileprocessing	2.125s
```

# Handle Edge Cases

[x] validate user input
[x] wrong filepath
[x] empty file
[x] handle different whitespaces

# Create test files

```
rm generatedfile.txt
for i in {1..10000000}
do
 echo "http://api.tech.com/item/${i}   ${i}" >> generatedfile.txt
done
```