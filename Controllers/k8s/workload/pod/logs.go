package pod

import (
	"sort"
	"strings"
)

// LINE_INDEX_NOT_FOUND is returned if requested line could not be found
var LINE_INDEX_NOT_FOUND = -1

// DefaultDisplayNumLogLines returns default number of lines in case of invalid request.
var DefaultDisplayNumLogLines = 100

// MaxLogLines is a number that will be certainly bigger than any number of logs. Here 2 billion logs is certainly much larger
// number of log lines than we can handle.
var MaxLogLines int = 2000000000

const (
	NewestTimestamp = "newest"
	OldestTimestamp = "oldest"
)

// Load logs from the beginning or the end of the log file.
// This matters only if the log file is too large to be loaded completely.
const (
	Beginning = "beginning"
	End       = "end"
)

// NewestLogLineId is the reference Id of the newest line.
var NewestLogLineId = LogLineId{
	LogTimestamp: NewestTimestamp,
}

// OldestLogLineId is the reference Id of the oldest line.
var OldestLogLineId = LogLineId{
	LogTimestamp: OldestTimestamp,
}

// DefaultSelection loads default log view selector that is used in case of invalid request
// Downloads newest DefaultDisplayNumLogLines lines.
var DefaultSelection = &Selection{
	OffsetFrom:      1 - DefaultDisplayNumLogLines,
	OffsetTo:        1,
	ReferencePoint:  NewestLogLineId,
	LogFilePosition: End,
}

// AllSelection returns all logs.
var AllSelection = &Selection{
	OffsetFrom:     -MaxLogLines,
	OffsetTo:       MaxLogLines,
	ReferencePoint: NewestLogLineId,
}

// LogDetails returns representation of log lines
type LogDetails struct {

	// Additional information of the logs e.g. container name, dates,...
	// 额外的日志信息，如容器名、日期等
	Info LogInfo `json:"info"`

	// Reference point to keep track of the position of all the logs
	// 通过ReferencePoint等信息，从所有日志中筛选日志
	Selection `json:"selection"`

	// Actual log lines of this page
	// 当前页的log（包含每行的时间戳和日志内容）
	LogLines `json:"logs"`
}

// LogInfo returns meta information about the selected log lines
type LogInfo struct {

	// Pod name.
	PodName string `json:"podName"`

	// The name of the container the logs are for.
	ContainerName string `json:"containerName"`

	// The name of the init container the logs are for.
	InitContainerName string `json:"initContainerName"`

	// Date of the first log line
	// 第一行日志的日期
	FromDate LogTimestamp `json:"fromDate"`

	// Date of the last log line
	// 最后一行日志的日期
	ToDate LogTimestamp `json:"toDate"`

	// Some log lines in the middle of the log file could not be loaded, because the log file is too large.
	// 日志截断
	Truncated bool `json:"truncated"`
}

// Selection of a slice of logs.
// It works just like normal slicing, but indices are referenced relatively to certain reference line.
// So for example if reference line has index n and we want to download first 10 elements in array we have to use
// from -n to -n+10. Setting ReferenceLogLineId the first line will result in standard slicing.
type Selection struct {
	// ReferencePoint is the ID of a line which should serve as a reference point for this selector.
	// You can set it to last or first line if needed. Setting to the first line will result in standard slicing.
	ReferencePoint LogLineId `json:"referencePoint"`
	// First index of the slice relatively to the reference line(this one will be included).
	// 相对于参考行的索引（在参考行左边为负值，在参考行右边为正值？？）
	OffsetFrom int `json:"offsetFrom"`
	// Last index of the slice relatively to the reference line (this one will not be included).
	// 相对于参考行的索引（在参考行左边为负值，在参考行右边为正值？？）
	OffsetTo int `json:"offsetTo"`
	// The log file is loaded either from the beginning or from the end. This matters only if the log file is too
	// large to be handled and must be truncated (to avoid oom)
	LogFilePosition string `json:"logFilePosition"`
}

// LogLineId uniquely identifies a line in logs - immune to log addition/deletion.
// 唯一标识日志中的一行 - 不受日志添加/删除的影响。
type LogLineId struct {
	// timestamp of this line.
	// 参考行的时间戳
	LogTimestamp `json:"timestamp"`
	// in case of timestamp duplicates (rather unlikely) it gives the index of the duplicate.
	// 如果时间戳重复（不太可能），它会给出重复项的索引。
	// For example if this LogTimestamp appears 3 times in the logs and the line is 1nd line with this timestamp,
	// 例如，如果此 LogTimestamp 在日志中出现 3 次并且该行是具有此时间戳的第一行，
	// then line num will be 1 or -3 (1st from beginning or 3rd from the end).
	// 那么LineNum将为 1 或 -3（从头开始第 1 行或从尾数第 3 行）。
	// If timestamp is unique then it will be simply 1 or -1 (first from the beginning or first from the end, both mean the same).
	// 如果 timestamp 是唯一的，那么LineNum将只是 1 或 -1（first from the beginning 或 first from the end，两者意思相同）
	LineNum int `json:"lineNum"`
}

// LogLines provides means of selecting log views. Problem with logs is that new logs are constantly added.
// Therefore the number of logs constantly changes and we cannot use normal indexing. For example
// if certain line has index N then it may not have index N anymore 1 second later as logs at the beginning of the list
// are being deleted. Therefore it is necessary to reference log indices relative to some line that we are certain will not be deleted.
// For example line in the middle of logs should have lifetime sufficiently long for the purposes of log visualisation. On average its lifetime
// is equal to half of the log retention time. Therefore line in the middle of logs would serve as a good reference point.
// LogLines allows to get ID of any line - this ID later allows to uniquely identify this line. Also it allows to get any
// slice of logs relatively to certain reference line ID.
type LogLines []LogLine

// LogLine is a single log line that split into timestamp and the actual content.
type LogLine struct {
	Timestamp LogTimestamp `json:"timestamp"`
	Content   string       `json:"content"`
}

// LogTimestamp is a timestamp that appears on the beginning of each log line.
type LogTimestamp string

// SelectLogs returns selected part of LogLines as required by logSelector, moreover it returns IDs of first and last
// of returned lines and the information of the resulting logView.
func (self LogLines) SelectLogs(logSelection *Selection) (LogLines, LogTimestamp, LogTimestamp, Selection, bool) {
	// 请求日志行数？
	requestedNumItems := logSelection.OffsetTo - logSelection.OffsetFrom
	// 通过ReferencePoint获取ReferenceLineIndex（参考日志在LogLines中的索引）
	referenceLineIndex := self.getLineIndex(&logSelection.ReferencePoint)

	if referenceLineIndex == LINE_INDEX_NOT_FOUND || requestedNumItems <= 0 || len(self) == 0 {
		// Requested reference line could not be found, probably it's already gone or requested no logs. Return no logs.
		// 如果ReferenceLineIndex索引为-1 或者 请求的日志项为0 或者 LogLines长度为0，返回空
		return LogLines{}, "", "", Selection{}, false
	}
	fromIndex := referenceLineIndex + logSelection.OffsetFrom // 起始索引（包含）
	toIndex := referenceLineIndex + logSelection.OffsetTo     // 结束索引（不包含）
	lastPage := false                                         // 是否时最后一页
	if requestedNumItems > len(self) {
		// 如果请求行数大于日志总行数
		fromIndex = 0
		toIndex = len(self)
		lastPage = true
	} else if toIndex > len(self) {
		// 如果结束索引大于日志总行数
		// fromIndex = fromIndex - toIndex + len(self) // 相当于fromIndex = len(self)-requestedNumItems
		fromIndex -= toIndex - len(self)
		toIndex = len(self)
		lastPage = logSelection.LogFilePosition == Beginning
	} else if fromIndex < 0 {
		toIndex += -fromIndex
		fromIndex = 0
		lastPage = logSelection.LogFilePosition == End
	}

	// set the middle of log array as a reference point, this part of array should not be affected by log deletion/addition.
	// 设置日志数组的中间位置作为新的参考点，并返回，供下次日志查询请求时做查询参数使用
	newSelection := Selection{
		ReferencePoint:  *self.createLogLineId(len(self) / 2),
		OffsetFrom:      fromIndex - len(self)/2,
		OffsetTo:        toIndex - len(self)/2,
		LogFilePosition: logSelection.LogFilePosition,
	}
	return self[fromIndex:toIndex], self[fromIndex].Timestamp, self[toIndex-1].Timestamp, newSelection, lastPage
}

// getLineIndex returns the index of the line (referenced from beginning of log array) with provided logLineId.
func (self LogLines) getLineIndex(logLineId *LogLineId) int {
	if logLineId == nil || logLineId.LogTimestamp == NewestTimestamp || len(self) == 0 || logLineId.LogTimestamp == "" {
		// if no line id provided return index of last item.
		// 返回LogLines的最后一项索引
		return len(self) - 1
	} else if logLineId.LogTimestamp == OldestTimestamp {
		// 返回LogLines的第一项索引
		return 0
	}

	// 时间戳不是（NewestTimestamp或OldestTimestamp）时，需要计算索引
	logTimestamp := logLineId.LogTimestamp

	// 初始化matchingStartedAt(满足logTimestamp条件的LogLines索引值)
	matchingStartedAt := 0
	// 遍历LogLines，当满足时间戳时，返回该日志项索引，作为日志筛选的起始索引（matchingStartedAt就是匿名函数返回true时的i值）
	matchingStartedAt = sort.Search(len(self), func(i int) bool {
		return self[i].Timestamp >= logTimestamp
	})

	// 返回logLines中，恰好满足logTimestamp条件的日志行数（应该是？？）
	linesMatched := 0
	// 匹配logTimestamp相同的行？？？？？logTimestamp居然不是唯一的？？？？？？？？？？
	// 目的可能是找出时间戳相同的项（这就是某个地方说的 几乎不可能出现的情况），匹配最后一项？？？？
	// 如果没有匹配到self[matchingStartedAt].Timestamp == logTimestamp，linesMatched为0，最终函数返回Line_index_NOT_FOUND=-1
	if matchingStartedAt < len(self) && self[matchingStartedAt].Timestamp == logTimestamp { // match found
		for (matchingStartedAt+linesMatched) < len(self) && self[matchingStartedAt+linesMatched].Timestamp == logTimestamp {
			linesMatched += 1
		}
	}

	// !!!!如果同一时间戳有多条日志（几乎不可能）可以将这些日志视为小的数组
	// 从左往右，LineNum从1开始取值；从右往左LineNum从-1开始取值
	// LineNum的值可以根据自己的需要设置取第几个值，函数最终会根据LineNum的取值，返回该日志在整个日志中的索引
	var offset int
	if logLineId.LineNum < 0 {
		offset = linesMatched + logLineId.LineNum
	} else {
		offset = logLineId.LineNum - 1
	}

	if 0 <= offset && offset < linesMatched {
		return matchingStartedAt + offset
	}
	// (上个if的else)没有匹配到符合时间戳的日志，返回-1
	return LINE_INDEX_NOT_FOUND
}

// CreateLogLineId returns ID of the line with provided lineIndex.
func (self LogLines) createLogLineId(lineIndex int) *LogLineId {
	logTimestamp := self[lineIndex].Timestamp
	// determine whether to use negative or positive indexing
	// check whether last line has the same index as requested line. If so, we can only use positive referencing
	// as more lines may appear at the end.
	// negative referencing is preferred as higher indices disappear later.
	var step int
	if self[len(self)-1].Timestamp == logTimestamp {
		// use positive referencing
		step = 1
	} else {
		step = -1
	}
	offset := step
	for ; 0 <= lineIndex-offset && lineIndex-offset < len(self); offset += step {
		if !(self[lineIndex-offset].Timestamp == logTimestamp) {
			break
		}
	}
	return &LogLineId{
		LogTimestamp: logTimestamp,
		LineNum:      offset,
	}
}

// ToLogLines converts rawLogs (string) to LogLines. Proper log lines start with a timestamp which is chopped off.
// In error cases the server returns a message without a timestamp
func ToLogLines(rawLogs string) LogLines {
	logLines := LogLines{}
	for _, line := range strings.Split(rawLogs, "\n") {
		if line != "" {
			startsWithDate := ('0' <= line[0] && line[0] <= '9') //2017-...
			idx := strings.Index(line, " ")
			if idx > 0 && startsWithDate {
				timestamp := LogTimestamp(line[0:idx])
				content := line[idx+1:]
				logLines = append(logLines, LogLine{Timestamp: timestamp, Content: content})
			} else {
				logLines = append(logLines, LogLine{Timestamp: "0", Content: line})
			}
		}
	}
	return logLines
}
