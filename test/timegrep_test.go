package main

import (
	"fmt"
	backscanner "github.com/clickpaas/dategrep/pkg/backscanner"
	timegrep "github.com/clickpaas/dategrep/pkg/timegrep"
	"log"
	"os"
	"syscall"
	"testing"
	"time"
)

func Test_logNoTime(t *testing.T) {
	println("begin-----")
	timegrep.SearchLogfile("2023-03-28 23:59:55", "2023-03-28 23:59:57", "logNoTime.log")
	println("end-----")
}
func Test_logSimple(t *testing.T) {
	println("begin-----")
	timegrep.SearchLogfile("2023-03-28 23:59:55", "2023-03-28 23:59:57", "logSimple.log")
	println("end-----")
}

func Test_logMiddleError(t *testing.T) {
	println("begin-----")
	timegrep.SearchLogfile("2023-03-28 23:59:55", "2023-03-28 23:59:59", "logMiddleError.log")
	println("end-----")
}

func Test_logEndError(t *testing.T) {
	println("begin-----")
	timegrep.SearchLogfile("2023-03-28 23:59:57", "2023-03-28 23:59:58", "logEndError.log")
	println("end-----")
}

func Test_logOneline(t *testing.T) {
	println("begin-----")
	timegrep.SearchLogfile("2023-03-28 23:59:19", "2023-03-28 23:59:21", "logOneLine.log")
	println("end-----")
}

func Test_FindLastLineWithTimeString(t *testing.T) {
	open, _ := os.Open("logEndError.log")
	startPos, endPos := timegrep.FindLastLineWithTimeString(open)
	print(startPos, endPos)
}

func Test_bigfile(t *testing.T) {
	timegrep.SearchLogfile("2023-06-15 14:30:31", "2023-06-15 14:31:31", "/Users/tingfeng/Downloads/app.log")
}

func Test_byteConvert(t *testing.T) {
	str := "abc123你好!"
	println(len(str))
	b := []byte(str)
	println(len(b))
}

func Test_tid2timestamp(t *testing.T) {
	tid, err := timegrep.ParseTid("1.25117452587.16800191604710001")
	println(err)
	println(tid.Format(timegrep.Layout))
}

func Test_filestat(t *testing.T) {
	//for _, arg := range os.Args[1:] {
	fileinfo, err := os.Stat("logNoTime.log")
	if err != nil {
		log.Fatal(err)
	}
	atime := fileinfo.Sys().(*syscall.Stat_t).Atimespec
	unix := time.Unix(atime.Sec, atime.Nsec)
	fmt.Println(unix)
	//}
}

func Test_backscanner_try(t *testing.T) {
	f, err := os.Open("logEndError.log")
	if err != nil {
		panic(any(err))
	}
	fi, err := f.Stat()
	if err != nil {
		panic(any(err))
	}
	defer f.Close()

	pos := fi.Size()
	//pos:=4470
	scanner := backscanner.New(f, pos)
	for {
		line, pos, err := scanner.LineBytes()
		if err != nil {
			fmt.Println("Error:", err)
			break
		}
		//if bytes.Contains(line, what) {
		fmt.Printf("Found at line position: %d, line: %q\n", pos, line)
		//break
		//}
	}
}
