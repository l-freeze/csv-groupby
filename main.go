package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)
import "flag"

func main() {
	filePath := flag.String("file", "", "Path to the CSV file")
	columnNameOrIndex := flag.String("column", "", "Name of the column")
	delimiter := flag.String("delimiter", ",", "Delimiter used in the CSV file")
	header := flag.Bool("header", false, "Indicate if the CSV file has a header row")

	// 隠しオプション io.ReaderのバッファサイズKBを指定する(デフォルトは4KB)
	bufferSize := flag.Int("buffer", 0, "Read buffer size in KB")

	// 隠しオプション worker pool size
	//workerPoolSize := flag.Int("worker", runtime.NumCPU(), "Number of worker pool size")

	flag.Parse()

	if *filePath == "" || *columnNameOrIndex == "" {
		fmt.Println("Please provide both file and column name")
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
	var columnIndex = -1
	if *header {
		// Read the first row (header)
		headerRow, err := reader.Read()
		if err != nil {
			log.Fatalf("Error reading header: %v", err)
		}

		// Find the index of the specified column
		for i, col := range headerRow {
			if col == *columnNameOrIndex {
				columnIndex = i
				break
			}
		}

		if columnIndex == -1 {
			log.Fatalf("Column %s not found in header", *columnNameOrIndex)
		}

		log.Printf("Column %s found at index %d\n", *columnNameOrIndex, columnIndex)

	} else {
		if idx, err := strconv.Atoi(*columnNameOrIndex); err == nil {
			columnIndex = idx
			fmt.Printf("Column index provided: %d\n", columnIndex)
		} else {
			log.Fatalf("Invalid column index: %s", *columnNameOrIndex)
		}
	}

	// Count group by colum
	counts := make(map[string]int)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error reading record: %v", err)
		}

		if columnIndex >= len(record) {
			log.Fatalf("Column index %d out of range for record: %v", columnIndex, record)
		}

		value := record[columnIndex]
		counts[value]++
	}

	// Print the counts
	fmt.Printf("Group by %v counts\n\n", *columnNameOrIndex)
	for name, count := range counts {
		fmt.Printf("%s: %d\n", name, count)
	}

}
