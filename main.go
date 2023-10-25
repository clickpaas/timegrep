package main

import (
	"flag"
	"fmt"
	"github.com/clickpaas/timegrep/pkg/timegrep"
	"strings"
	"time"
)

func main() {
	var (
		start string
		end   string
		//pattern   string
		path string
		tid  string
	)

	flag.StringVar(&start, "s", "", "start time")
	flag.StringVar(&end, "e", "", "end time")
	//flag.StringVar(&pattern, "pattern", "", "pattern to match")
	flag.StringVar(&path, "p", "", "path to search")
	flag.StringVar(&tid, "t", "", "tid to search")
	flag.Parse()

	if start == "" && end == "" && /*pattern == "" ||*/ path == "" && tid == "" {
		fmt.Println("Usage1: timegrep -s startTime -e endTime -p pathOrFile")
		fmt.Println("Usage2: timegrep -t tid -p pathOrFile")
		fmt.Println("Usage3: timegrep -t tid1,tid2,tid3 -p pathOrFile")
		//os.Exit(1)
		return
	}

	if tid != "" && path != "" {
		if !strings.Contains(tid, ".") {
			tid = strings.ReplaceAll(tid, "d", ".")
		}
		if strings.Contains(tid, ",") {
			minTime, maxTime := timegrep.ParseTidArr(tid)
			startTime := minTime.Add(-15 * time.Second)
			endTime := maxTime.Add(300 * time.Second)
			searchPathOrDirectory(path, startTime, endTime)
		} else {
			tidTime, err := timegrep.ParseTid(tid)
			if err != nil {
				fmt.Println("fail to parse tid,maybe it's not contains timestamp,tid=%s,err=%v", tid, err)
				return
			}
			startTime := tidTime.Add(-15 * time.Second)
			endTime := tidTime.Add(300 * time.Second)
			searchPathOrDirectory(path, startTime, endTime)
		}
	} else if start != "" && end != "" && path != "" {
		startTime, err := time.ParseInLocation(timegrep.Layout, start, time.Local)
		if err != nil {
			fmt.Println("Invalid startTime time format")
			//os.Exit(1)
			return
		}

		endTime, err := time.ParseInLocation(timegrep.Layout, end, time.Local)
		if err != nil {
			fmt.Println("Invalid endTime time format")
			//os.Exit(1)
			return
		}
		searchPathOrDirectory(path, startTime, endTime)
	}

}

func searchPathOrDirectory(path string, startTime time.Time, endTime time.Time) {
	if timegrep.IsDir(path) {
		files, err := timegrep.GetDirAllFilePaths(path)
		if err != nil {
			return
		}
		for i := range files {
			file := files[i]
			println("grepPath:", file)
			timegrep.SearchLogfile2(startTime, endTime, file)
		}
		return
	}
	if strings.HasSuffix(path, ".log") {
		timegrep.SearchLogfile2(startTime, endTime, path)
	}
}
