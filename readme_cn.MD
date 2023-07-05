# 性能对比 grep VS timegrep
文件最后几分钟发生的日志检索速度, timegrep比grep快80倍以上
## timegrep 按如下步骤执行
- 搜索时间范围在文件最后修改时间之后,则直接跳过,这样可以跳过大量的历史文件
- 通过二分法定位日志时间行,大文件速度也非常快 
- 日志时间行后的异常也输出
- 如果使用tid搜索则把tid中的日志时间提取出来放大一些搜索范围,再找出tid的行,放大日志时间范围 (-15seconds-->+300seconds) 
# 使用方法:
```
Usage1: timegrep -s startTime -e endTime -p pathOrFile
Usage2: timegrep -t tid -p pathOrFile

```
# 示例:
```
 ./timegrep.amd64 -s "2023-03-28 23:59:55" -e "2023-03-28 23:59:57" -p ./
 ./timegrep.amd64 -s "2023-03-28 23:59:55" -e "2023-03-28 23:59:57" -p ./test/a.log
 ./timegrep.amd64 -t 1.25117452587.16800191604710001 -p ./test/logSimple.log
```
![指定时间范围搜索日志](test/searchWithTimeRange.png)

![搜索结果包含异常信息](test/searchWithException.png)

# 编译
- linux `GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o timegrep.amd64`
- arm64v8 `GOOS=linux GOARCH=arm64 go build -o timegrep.arm64`
- macos `go build`
# 待办:
0. 支持多种日志格式
1. 抽查第一行来确定日志格式,后续同一文件使用同种日志格式
2. 代码中使用了字符串,可替换成bytes加速 ,替换 scanner.Text()  为 scanner.Bytes()

#  (性能对比) grep VS timegrep 
![test case](test/grepVStimegrep.png)

```
timer_start=`date "+%Y-%m-%d %H:%M:%S"`
find ./ -name "*.log" -exec grep '2023-06-15 16:26' {} \; &> /tmp/test1.log
timer_end=`date "+%Y-%m-%d %H:%M:%S"`
echo $(($(date +%s -d "${timer_end}") - $(date +%s -d "${timer_start}")))

timer_start=`date "+%Y-%m-%d %H:%M:%S"`
./timegrep.amd64 -s="2023-06-15 16:26:00" -e="2023-06-15 16:26:59" -p ./ &> /tmp/test2.log
timer_end=`date "+%Y-%m-%d %H:%M:%S"`
echo $(($(date +%s -d "${timer_end}") - $(date +%s -d "${timer_start}")))
```