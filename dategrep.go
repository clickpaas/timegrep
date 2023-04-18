package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"regexp"
	"sync/atomic"
	"time"
)

type LogEntry struct {
	Timestamp time.Time
	Message   []byte
}

const pattern = "^\\d{4}-\\d{1,2}-\\d{1,2} \\d{1,2}:\\d{1,2}:\\d{1,2}.\\d{1,3}"
const min_log_size = 20

func main() {
	// Open the log file for reading
	//file, err := os.Open("logfile.log")
	file, err := os.Open("/Users/tingfeng/work/java/python-monitor/log/app-20230220.1.log")
	// 2023-03-28 23:59:55
	//var startTime time.Time = time.Date(2023, 3, 28, 23, 59, 57, 0, time.Local)
	var startTimeString = "2023-02-20 18:14:03";
	var endTimeString = "2023-02-20 18:14:04";
	//var startTimeString = "2023-03-28 23:59:50"
	//var endTimeString = "2023-08-28 23:59:57"
	var startTime, _ = time.ParseInLocation("2006-01-02 15:04:05", startTimeString, time.Local)
	var endTime, _ = time.ParseInLocation("2006-01-02 15:04:05", endTimeString, time.Local)
	endTime = endTime.Add(time.Millisecond * 999)
	re := regexp.MustCompile(pattern)
	bufFile := bufio.NewReader(file)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Get the file size
	fileInfo, err := file.Stat()
	if err != nil {
		panic(err)
	}

	// TODO 优化 起始时间定位行 结束时间定位行 ,其中找过的可以缓存

	startLogOffset := locateOffsetByTime(fileInfo, bufFile, file, re, startTime, false)
	println("=========")
	endLogOffset := locateOffsetByTime(fileInfo, bufFile, file, re, endTime, true)
	fmt.Printf("startLogOffset:%v,endLogOffset:%v\n", startLogOffset,endLogOffset)
	if startLogOffset >= endLogOffset {
		return
	}

	_, err = file.Seek(startLogOffset, 0)
	if err != nil && err != io.EOF {
		panic(err)
	}

	size := (endLogOffset - startLogOffset)
	var sum int64 = 0

	//for {
	//	var line []byte
	//	for {
	//		nextByte := make([]byte, 1)
	//		_, err := file.Read(nextByte)
	//		if err != nil && err != io.EOF {
	//			panic(err)
	//		}
	//		if nextByte[0] == '\n' || err == io.EOF {
	//			line = append(line, nextByte[0])
	//			sum = atomic.AddInt64(&sum, int64(len(line)))
	//			fmt.Printf("\n----result sum= %v, line %s: ", sum,string(line))
	//			break
	//		}
	//		line = append(line, nextByte[0])
	//	}
	//
	//	if sum > size {
	//		break
	//	}
	//}

	bufFile.Reset(file)
	for {
		buf, err := bufFile.ReadBytes('\n')
		if err != nil && err != io.EOF {
			panic(err)
		}

		fmt.Printf("%s", string(buf))
		atomic.AddInt64(&sum, int64(len(buf)))
		if sum >= size {
			break
		}
		//cur_offset,_:=file.Seek(0,os.SEEK_CUR)
		//if cur_offset>endLogOffset{
		//	break
		//}
	}

}

func locateOffsetByTime(fileInfo fs.FileInfo, bufFile *bufio.Reader, file *os.File, re *regexp.Regexp, searchTime time.Time, isEndTimeSearch bool) int64 {
	// Perform binary search until log entry is found or search range is exhausted

	fileSize := fileInfo.Size()
	// 文件修改时间超出范围
	if fileInfo.ModTime().Before(searchTime) {
		return fileSize
	}

	lastTimeString, err := ReverseReadLastTime(file)
	if err != nil {
		return -1
	}
	lastTime, err := time.ParseInLocation("2006-01-02 15:04:05.000", *lastTimeString, time.Local)
	// 判断最后带时间的一行是否超出时间范围,如果是结束
	if lastTime.Before(searchTime) {
		return fileSize
	}
	// Define the search range
	firstOffset := int64(0)
	lastOffset := fileSize - 1
	var logEntry LogEntry
	midpointOffset := lastOffset
	for firstOffset <= lastOffset - min_log_size {
		// Calculate midpoint offset
		midpointOffset = (firstOffset + lastOffset) / 2
		bufFile.Reset(file)
		// Seek to midpoint offset
		_, err := file.Seek(midpointOffset, 0)
		if err != nil {
			panic(err)
		}
		// Read the next full line, including the newline character
		line, timeStr := readNewLine(bufFile, re)
		// Parse the log entry timestamp

		logEntry.Timestamp, err = time.ParseInLocation("2006-01-02 15:04:05.000", timeStr, time.Local)
		if err != nil {
			panic(err)
		}
		logEntry.Message = line

		fmt.Printf("peek midpointOffset=%v, entry=%v\n", midpointOffset, string(logEntry.Message))
		if searchTime.Before(logEntry.Timestamp) {
			lastOffset = midpointOffset - 1
		} else if searchTime.After(logEntry.Timestamp) {
			firstOffset = midpointOffset + 1
		} else {
			fmt.Printf("Found log entry: %v\n", logEntry)
			return midpointOffset
		}
	}
	//line, timeStr := readNewLine(bufFile, re)
	//fmt.Printf("at last, Found log entry: %s,%s", timeStr, line)
	if isEndTimeSearch {
		atomic.AddInt64(&midpointOffset, int64(len(logEntry.Message)))
	}
	return midpointOffset
}

func readNewLine(bufFile *bufio.Reader, re *regexp.Regexp) ([]byte, string) {
	var line []byte
	var timeStr string
	for {
		buf, err := bufFile.ReadBytes('\n')
		if err != nil && err != io.EOF {
			panic(err)
		}

		if re.Match(buf) {
			line = buf
			timeStr = string(re.Find(line))
			//fmt.Printf("timeStr=#%s#buf = #%s#\n", timeStr, string(buf))
			break
		}
	}
	return line, timeStr
}

func ReverseReadLastTime(file *os.File) (*string, error) {
	re := regexp.MustCompile(pattern)
	lineNum := 300
	//打开文件
	//defer file.Close()
	//获取文件大小
	fs, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := fs.Size()

	var offset int64 = -1   //偏移量，初始化为-1，若为0则会读到EOF
	char := make([]byte, 1) //用于读取单个字节
	lineStr := ""           //存放一行的数据
	buff := make([]string, 0, 100)
	for (-offset) <= fileSize {
		//通过Seek函数从末尾移动游标然后每次读取一个字节
		file.Seek(offset, io.SeekEnd)
		_, err := file.Read(char)
		if err != nil {
			return nil, err
		}
		if char[0] == '\n' {
			offset--  //windows跳过'\r'
			lineNum-- //到此读取完一行
			findString := re.FindString(lineStr)
			if findString != "" {
				return &findString, nil
			}
			buff = append(buff, lineStr)
			lineStr = ""
			if lineNum == 0 {
				return &findString, nil
			}
		} else {
			lineStr = string(char) + lineStr
		}
		offset--
	}
	//buff = append(buff, lineStr)
	//return buff, nil
	return nil, nil
}
