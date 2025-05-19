package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)
import "flag"

type GroupingIndex struct {
	Index    int
	JsonPath string
}

func main() {
	filePath := flag.String("file", "", "Path to the CSV file")
	columnNamesOrIndexes := flag.String("column", "", "ColumnName or columnIndex. Header is required if ColumnName is used.")
	delimiter := flag.String("delimiter", ",", "Delimiter used in the CSV file")
	header := flag.Bool("header", false, "Indicate if the CSV file has a header row. Header is required if ColumnName is used.")

	// io.ReaderのバッファサイズKBを指定する(デフォルトは4KB)
	bufferSize := flag.Int("buffer", 0, "Read buffer size in KB")

	// worker pool size
	workerPoolSize := flag.Int("worker", 1, fmt.Sprintf("Number of worker pool size(max:%d)", runtime.NumCPU()))

	flag.Parse()

	if *workerPoolSize > runtime.NumCPU() {
		fmt.Printf("Worker pool size is limited to %d\n", runtime.NumCPU())
		*workerPoolSize = runtime.NumCPU()
		os.Exit(0)
	}

	if *filePath == "" || *columnNamesOrIndexes == "" {
		fmt.Println("Please provide both -file and -column")
		os.Exit(1)
	}

	// Open the CSV file
	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	var ioReader io.Reader = file
	if *bufferSize > 0 {
		ioReader = bufio.NewReaderSize(file, *bufferSize*1024)
		if *bufferSize > 1024 {
			log.Printf("Using buffer size: %d MB\n", *bufferSize/1024)
		} else {
			log.Printf("Using buffer size: %d KB\n", *bufferSize)
		}
	}

	reader := csv.NewReader(ioReader)
	reader.Comma = rune((*delimiter)[0]) // Set the delimiter

	// Satisfy the column index
	var groupingIndexes []GroupingIndex
	if *header {
		// Read the first row (header)
		headerRow, err := reader.Read()
		if err != nil {
			log.Fatalf("Error reading header: %v", err)
		}

		// Find the index of the specified column
		splitColumnNames := strings.Split(*columnNamesOrIndexes, ",")
		var columnNames []string
		for _, col := range splitColumnNames {
			if col != "" {
				columnNames = append(columnNames, col)
			}
		}

		for _, columnName := range columnNames {
			parts := strings.SplitN(columnName, "#", 2)
			var jsonPath string
			groupingColumn := parts[0]
			if len(parts) == 2 {
				jsonPath = parts[1]
			}

			for i, col := range headerRow {
				if col == groupingColumn {
					groupingIndexes = append(groupingIndexes, GroupingIndex{Index: i, JsonPath: jsonPath})
					break
				}
			}
		}

		if groupingIndexes == nil {
			log.Fatalf("Column %s not found in header", *columnNamesOrIndexes)
		}

		log.Printf("Column %s found at index %#v\n", *columnNamesOrIndexes, groupingIndexes)

	} else {
		splitColumnIndexes := strings.Split(*columnNamesOrIndexes, ",")
		for _, columnIndex := range splitColumnIndexes {

			parts := strings.SplitN(columnIndex, "#", 2)
			idxString := parts[0]
			if idx, err := strconv.Atoi(idxString); err == nil {
				var jsonPath string

				if len(parts) == 2 {
					jsonPath = parts[1]
				}

				groupingIndexes = append(groupingIndexes, GroupingIndex{Index: idx, JsonPath: jsonPath})

			} else {
				log.Fatalf("Invalid column index: %s", *columnNamesOrIndexes)
			}

		}
	}

	if len(groupingIndexes) == 0 {
		log.Fatalf(" Unspecified column")
	}

	// worker pool
	var wg sync.WaitGroup
	jobs := make(chan []string, 1000) // CSVの行を格納するチャネル
	results := make(chan map[string]int, *workerPoolSize)

	for i := 0; i < *workerPoolSize; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			log.Printf("Worker %d started\n", workerID)

			subCounts := make(map[string]int)
			for job := range jobs {
				for _, columnIndex := range groupingIndexes {
					if columnIndex.Index < 0 || columnIndex.Index >= len(job) {
						log.Fatalf("Column index %d out of range for record: %v", columnIndex, job)
					}
				}

				var countKeys []string
				for _, columnIndex := range groupingIndexes {
					if columnIndex.JsonPath != "" {
						cellValue := job[columnIndex.Index]
						json := gjson.Get(cellValue, columnIndex.JsonPath)
						countKey := columnIndex.JsonPath + "=" + json.String()
						countKeys = append(countKeys, countKey)
					} else {
						countKeys = append(countKeys, job[columnIndex.Index])
					}
				}
				subCounts[strings.Join(countKeys, ",")]++

			}
			//for name, count := range subCounts {
			//	fmt.Printf("[worker:%d]%s: %d\n", workerID, name, count)
			//}

			results <- subCounts
			log.Printf("Worker %d finished\n", workerID)
		}(i)
	}

	// Count group by colum
	go func() {
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error reading record: %v", err)
			}
			jobs <- record
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()
	log.Printf("All workers finished\n")

	counts := make(map[string]int)
	for count := range results {
		for name, c := range count {
			counts[name] += c
		}
	}

	// Print the counts
	for name, count := range counts {
		fmt.Printf("%s: %d\n", name, count)
	}
	fmt.Printf("Group by %v counts\n\n", *columnNamesOrIndexes)

}
