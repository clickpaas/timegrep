# performance grep VS timegrep
timegrep is 80+ times faster than grep with last several minutes scope
- skip grep old file that modified-date older than search time
- find log-time in log with binary search 
- search result include exception that following log-time
- tid include time , and timegrep will search time between (-15seconds-->+300seconds) 
# usage:
```
 ./timegrep.amd64 -start="2023-03-28 23:59:55" -end="2023-03-28 23:59:57" -path ./
 ```

```
 ./timegrep.amd64 -start="2023-03-28 23:59:55" -end="2023-03-28 23:59:57" -path ./test/a.log
```
```
 ./timegrep.amd64 -tid="1.25117452587.16800191604710001" -path ./test/logSimple.log
 ```
# build
- linux `GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o timegrep.amd64`
- arm64v8 `GOOS=linux GOARCH=arm64 go build -o timegrep.arm64`
- macos `go build`
# TODO:
0. support several log format
1. read first line to select a format layout
2. replace scanner.Text()  with scanner.Bytes()

#  grep VS timegrep (performance)
![test case](test/grepVStimegrep.png)

```
timer_start=`date "+%Y-%m-%d %H:%M:%S"`
find ./ -name "*.log" -exec grep '2023-06-15 16:26' {} \; &> /tmp/test1.log
timer_end=`date "+%Y-%m-%d %H:%M:%S"`
echo $(($(date +%s -d "${timer_end}") - $(date +%s -d "${timer_start}")))

timer_start=`date "+%Y-%m-%d %H:%M:%S"`
./timegrep.amd64 -start="2023-06-15 16:26:00" -end="2023-06-15 16:26:59" -path ./ &> /tmp/test2.log
timer_end=`date "+%Y-%m-%d %H:%M:%S"`
echo $(($(date +%s -d "${timer_end}") - $(date +%s -d "${timer_start}")))
```