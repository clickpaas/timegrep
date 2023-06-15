package timegrep

import (
	"bufio"
	"fmt"
	backscanner "github.com/clickpaas/dategrep/pkg/backscanner"
	"os"
	"strings"
	"time"
)

// eg: go run dategrep.go -start="2023-03-28 23:59:55" -end="2023-03-28 23:59:57" -file logfile.log
const layout = "2006-01-02 15:04:05"
const timeStringLength = len(layout)

// default java stack size and one line exception
const maxStackSize = 1024 + 1

func SearchLogfile(startTime string, endTime string, file string) {
	start, err := time.ParseInLocation(layout, startTime,time.Local)
	if err != nil {
		fmt.Println("Invalid start time format")
		//os.Exit(1)
		return
	}

	end, err := time.ParseInLocation(layout, endTime,time.Local)
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
	//f2, err := os.Stat(file)
	//
	//f2.ModTime()
	fileHandle, err := os.Open(file)

	fs, err := fileHandle.Stat()
	if err != nil {
		return
	}
	//local := time.FixedZone("CST", 8*3600)
	//loc, _ := time.LoadLocation("Asia/Shanghai")
	//time.Local = loc
	//finfo, _ := os.Stat(file)

	//linuxFileAttr := finfo.Sys().(*syscall.Stat_t)
	//mtime := time.Unix(linuxFileAttr.Mtimespec.Sec, 0)
	//fmt.Println("最后修改时间", finfo.ModTime())
	//fmt.Println("最后修改时间", start)
	//fs.ModTime().In(local)
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

	//fileSize := stat.Size()
	//lower := int64(0)
	//upper := fileSize
	startOffset := int64(0)

	lower := startOffset
	startPos, endPos := FindLastLineWithTimeString(fileHandle)
	upper := endPos
	if startPos != int64(0) {
		for lower <= upper {
			mid := (lower + upper) / 2
			fileHandle.Seek(mid, 0)
			scanner := bufio.NewScanner(fileHandle)
			canBreak, timestamp, scanSize, err := scanOneLineStartWithTime(scanner)
			if canBreak {
				startOffset = mid + scanSize
				break
			}
			if err != nil {
				continue
			}
			if timestamp.Before(start) {
				lower = mid + scanSize
			} else {
				startOffset = mid
				upper = mid - 1
			}
		}
	}
	fileHandle.Seek(startOffset, 0)
	scanner := bufio.NewScanner(fileHandle)
	hasFirstTimeString := false
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) >= timeStringLength {
			timestamp, err := time.ParseInLocation(layout, line[:timeStringLength],time.Local)
			if err != nil {
				// print exception line
				if hasFirstTimeString {
					fmt.Println(line)
				}
				continue
			}
			hasFirstTimeString = true
			if timestamp.After(end) {
				break
			}
		}
		//if patternRegexp.MatchString(line) {
		//	fmt.Println(line)
		//}
		if hasFirstTimeString {
			fmt.Println(line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file")
		//os.Exit(1)
		return
	}
}

func scanOneLineStartWithTime(scanner *bufio.Scanner) (bool, time.Time, int64, error) {
	canBreak := false
	var scanSize int64
	for i := 0; i < maxStackSize; i++ {
		if !scanner.Scan() {
			canBreak = true
			break
		}
		line := scanner.Text()
		lineLength := len(line)
		// +1 : scan one line include end character=\n
		scanSize += int64(lineLength) + 1
		if lineLength >= timeStringLength {
			timeString := line[:timeStringLength]
			timestamp, err := time.ParseInLocation(layout, timeString,time.Local)
			if err != nil {
				continue
			}
			return canBreak, timestamp, scanSize, err
		}
	}
	return true, time.Time{}, scanSize, nil
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

func FindLastLineWithTimeString(f *os.File) (int64, int64) {

	fi, err := f.Stat()
	if err != nil {
		panic(any(err))
	}
	//defer f.Close()

	pos := fi.Size()
	scanner := backscanner.New(f, pos)
	//scanSize:=int64(0)
	for i := 0; i < maxStackSize; i++ {
		line, pos, err := scanner.Line()
		if err != nil {
			fmt.Println("Error:", err)
			break
		}
		//if bytes.Contains(line, what) {
		//fmt.Printf("Found %q at line position: %d, line: %q\n", what, pos, line)
		//break
		//}
		lineLength := len(line)
		//scanSize += int64(lineLength) + 1
		if lineLength >= timeStringLength {
			timeString := line[:timeStringLength]
			_, err := time.ParseInLocation(layout, timeString,time.Local)
			if err != nil {
				continue
			}
			return pos, pos + int64(lineLength) + 1
		}
	}
	return 0, pos
}
