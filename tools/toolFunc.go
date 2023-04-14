package tools

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"
)

// IsStructEmpty 判断结构体是否为空
func IsStructEmpty(x, y interface{}) bool {
	return reflect.DeepEqual(x, y)
}

// TimeStrToUTCTime 时间字符串转换为UTC时间
func TimeStrToUTCTime(layout string, s string) (time.Time, error) {
	t, err := time.Parse(layout, s)
	return t.UTC(), err
}

// UTCTimeToTimeStr 将时间按指定格式输出为字符串
func UTCTimeToTimeStr(t time.Time, layout string) string {
	return t.UTC().Format(layout)
}

// DeltaTime 判断获取的时间相对于当前时间是多久以前
func DeltaTime(firstTime, lastTime time.Time) string {
	// 时间减法,获取时间差值
	deltaTime := lastTime.UTC().Sub(firstTime.UTC())
	// 将差值转换为小时
	//h := deltaTime.Hours()
	// 将差值转换为分钟
	//m := deltaTime.Minutes()
	// 将差值转换为秒钟
	deltaSeconds := deltaTime.Seconds()
	//if (h / 24) >= 1 {
	//	result = fmt.Sprintf("%d天前", int64(h/24))
	//} else if (m / 60) >= 1 {
	//	result = fmt.Sprintf("%d小时前", int64(m/60))
	//} else if (s / 60) >= 1 {
	//	result = fmt.Sprintf("%d分钟前", int64(s/60))
	//} else if s > 0 {
	//	result = fmt.Sprintf("%d秒前", int64(s))
	//} else {
	//	result = "Wrong Time"
	//}
	hour := math.Trunc(deltaSeconds / 3600)
	minute := math.Trunc((deltaSeconds - hour*3600) / 60)
	second := deltaSeconds - hour*3600 - minute*60

	var day float64 = 0
	if hour >= 24 {
		day = math.Trunc(hour / 24)
		hour = hour - day*24
	}

	var deltaTimeString string
	if day == 0 {
		if hour == 0 {
			if minute == 0 {
				deltaTimeString = fmt.Sprintf("%ds", int(second))
			} else {
				deltaTimeString = fmt.Sprintf("%dm%ds", int(minute), int(second))
			}
		} else {
			deltaTimeString = fmt.Sprintf("%dh%dm%ds", int(hour), int(minute), int(second))
		}
	} else {
		deltaTimeString = fmt.Sprintf("%dd%dh%dm%ds", int(day), int(hour), int(minute), int(second))
	}
	return deltaTimeString
}

// ConvertInt64ToInt 转换int64为int
func ConvertInt64ToInt(i int64) int {
	s := strconv.FormatInt(i, 10)
	m, _ := strconv.Atoi(s)
	return m
}

// ConvertMapToStruct convert map to struct
func ConvertMapToStruct(m map[string]interface{}, i interface{}) (err error) {
	// convert map to json bytes
	jsonBytes, err := json.Marshal(m)

	// convert json bytes to struct
	err = json.Unmarshal(jsonBytes, &i)

	return err
}
