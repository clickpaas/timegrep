package main

import (
	"flag"
	"fmt"
	"github.com/clickpaas/dategrep/pkg/timegrep"
	"strings"
)

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
	if timegrep.IsDir(path) {
		files, err := timegrep.GetDirAllFilePaths(path)
		if err != nil {
			return
		}
		for i := range files {
			file := files[i]
			println("grepPath:", file)
			timegrep.SearchLogfile(startTime, endTime, file)
		}
		return
	}
	if strings.HasSuffix(path, ".log") {
		timegrep.SearchLogfile(startTime, endTime, path)
	}
}
