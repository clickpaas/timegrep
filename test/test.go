package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
)




func main2() {
	file, err := os.Open("/Users/tingfeng/work/golang/src/github.com/carlvine500/dategrep/logfile.log")
	file.Seek(1470, 0)
	bufFile := bufio.NewReader(file)
	buf, err := bufFile.ReadBytes('\n')
	if err != nil && err != io.EOF {
		panic(err)
	}
	println(len(buf))
	fmt.Printf("\n----result line #%s#: ", string(buf))
}
const pattern = "^\\d{4}-\\d{1,2}-\\d{1,2} \\d{1,2}:\\d{1,2}:\\d{1,2}.\\d{1,3}"

func main() {
	file, err := os.Open("/Users/tingfeng/work/golang/src/github.com/carlvine500/dategrep/logfile.log")
	if err != nil {
		return
	}

	read, _ := ReverseReadLastTime(file)
	print(*read)
}

func ReverseReadLastTime(file *os.File) (*string, error) {
	re := regexp.MustCompile(pattern)
	lineNum:=300
	//打开文件
	defer file.Close()
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
				return &findString,nil;
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
	return nil,nil;
}