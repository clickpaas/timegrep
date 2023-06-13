package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

// eg: go run dategrep.go -start="2023-03-28 23:59:55" -end="2023-03-28 23:59:57" -file logfile.log
const layout = "2006-01-02 15:04:05"

func main() {
	var (
		startTime string
		endTime   string
		//pattern   string
		path string
	)

	flag.StringVar(&startTime, "start", "", "start time")
	flag.StringVar(&endTime, "end", "", "end time")
	//flag.StringVar(&pattern, "pattern", "", "pattern to match")
	flag.StringVar(&path, "path", "", "path to search")
	flag.Parse()

	if startTime == "" || endTime == "" || /*pattern == "" ||*/ path == "" {
		fmt.Println("Usage: timegrep -start START_TIME -end END_TIME -pattern PATTERN -path FILE")
		//os.Exit(1)
		return
	}

	if IsDir(path) {
		files, err := GetDirAllFilePaths(path)
		if err != nil {
			return
		}
		for i := range files {
			file := files[i]
			println("file:", file)
			searchLogfile(startTime, endTime, file)
		}
		return
	}
	if strings.HasSuffix(path, ".log") {
		searchLogfile(startTime, endTime, path)
	}
}

func searchLogfile(startTime string, endTime string, file string) {
	start, err := time.Parse(layout, startTime)
	if err != nil {
		fmt.Println("Invalid start time format")
		//os.Exit(1)
		return
	}

	end, err := time.Parse(layout, endTime)
	if err != nil {
		fmt.Println("Invalid end time format")
		//os.Exit(1)
		return
	}

	//patternRegexp, err := regexp.Compile(pattern)
	//if err != nil {
	//	fmt.Println("Invalid pattern format")
	//	os.Exit(1)
	//}

	fileHandle, err := os.Open(file)

	fs, err := fileHandle.Stat()
	if fs.ModTime().Before(start) {
		//os.Exit(1)
		return
	}

	if err != nil {
		fmt.Println("Error opening file")
		//os.Exit(1)
		return
	}
	defer fileHandle.Close()

	stat, err := fileHandle.Stat()
	if err != nil {
		fmt.Println("Error getting file size")
		//os.Exit(1)
		return
	}

	fileSize := stat.Size()
	lower := int64(0)
	upper := fileSize
	startOffset := int64(0)
	endOffset := int64(0)
	if endOffset > 0 {
		// skip warn
	}

	for lower <= upper {
		mid := (lower + upper) / 2
		fileHandle.Seek(mid, 0)
		scanner := bufio.NewScanner(fileHandle)
		//for tryTimes:=0;tryTimes<3;tryTimes++{
		//
		//}

		canBreak, timestamp, err := scanOneLineStartWithTime(scanner)
		if canBreak {
			break
		}
		if err != nil {
			continue
		}
		if timestamp.After(end) {
			upper = mid - 1
		} else {
			endOffset = mid
			lower = mid + 1
		}
	}

	lower = startOffset
	upper = fileSize

	for lower <= upper {
		mid := (lower + upper) / 2
		fileHandle.Seek(mid, 0)
		scanner := bufio.NewScanner(fileHandle)
		canBreak, timestamp, err := scanOneLineStartWithTime(scanner)
		if canBreak {
			break
		}
		if err != nil {
			continue
		}
		if timestamp.Before(start) {
			lower = mid + 1
		} else {
			startOffset = mid
			upper = mid - 1
		}
	}

	fileHandle.Seek(startOffset, 0)
	scanner := bufio.NewScanner(fileHandle)
	for scanner.Scan() {
		line := scanner.Text()
		timestamp, err := time.Parse(layout, line[:19])
		if err != nil {
			continue
		}
		if timestamp.After(end) {
			break
		}
		//if patternRegexp.MatchString(line) {
		//	fmt.Println(line)
		//}
		fmt.Println(line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file")
		//os.Exit(1)
		return
	}
}

func scanOneLineStartWithTime(scanner *bufio.Scanner) (bool, time.Time, error) {
	canBreak := false
	if !scanner.Scan() {
		canBreak = true
	}
	line := scanner.Text()
	timeString := ""
	if len(line) >= 19 {
		timeString = line[:19]
	}
	timestamp, err := time.Parse(layout, timeString)
	if err != nil {
		if !scanner.Scan() {
			canBreak = true
		}
		line := scanner.Text()
		timeString2 := ""
		if len(line) >= 19 {
			timeString2 = line[:19]
		}
		timestamp, err = time.Parse(layout, timeString2)
	}
	return canBreak, timestamp, err
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func GetDirAllFilePaths(dirname string) ([]string, error) {
	// Remove the trailing path separator if dirname has.
	dirname = strings.TrimSuffix(dirname, string(os.PathSeparator))

	infos, err := os.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(infos))
	for _, info := range infos {
		path := dirname + string(os.PathSeparator) + info.Name()
		if info.IsDir() {
			tmp, err := GetDirAllFilePaths(path)
			if err != nil {
				return nil, err
			}
			paths = append(paths, tmp...)
			continue
		}
		if strings.HasSuffix(info.Name(), ".log") {
			paths = append(paths, path)
		}
	}
	return paths, nil
}
