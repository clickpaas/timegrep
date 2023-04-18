package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

var previousOffset int64 = 0

func main() {
	flag.Parse()

	//filename := flag.Arg(0)
	filename := "a.log"

	//delay := time.Tick(2 * time.Second)

	//for _ = range delay {
		readLastLine(filename)
	//}

}

func readLastLine(filename string) {

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	// we need to calculate the size of the last line for file.ReadAt(offset) to work

	// NOTE : not a very effective solution as we need to read
	// the entire file at least for 1 pass :(

	lastLineSize := 0

	for {
		line, _, err := reader.ReadLine()

		if err == io.EOF {
			break
		}

		lastLineSize = len(line)
	}

	fileInfo, err := os.Stat(filename)

	// make a buffer size according to the lastLineSize
	buffer := make([]byte, lastLineSize)

	// +1 to compensate for the initial 0 byte of the line
	// otherwise, the initial character of the line will be missing

	// instead of reading the whole file into memory, we just read from certain offset

	offset := fileInfo.Size() - int64(lastLineSize+1)
	numRead, err := file.ReadAt(buffer, offset)

	if previousOffset != offset {

		// print out last line content
		buffer = buffer[:numRead]
		fmt.Printf("%s \n", buffer)

		previousOffset = offset
	}

}