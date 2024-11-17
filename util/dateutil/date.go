// description: keybox
//
// @author: xwc1125
// @date: 2020/8/21 0021
package dateutil

import (
	"fmt"
	"strings"
	"time"
)

// 返回毫秒
func CurrentTime() int64 {
	return time.Now().UnixNano() / 1e6
}

// 获取间隔时间错，传入的是毫秒
func GetDistanceTime(diffMS int64) string {
	isNegative := false
	if diffMS < 0 {
		isNegative = true
		diffMS = -diffMS
	}
	s := diffMS / 1000 // 秒
	m := s / 60        // 分钟
	h := m / 60        // 小时
	day := h / 24      // 天
	hour := h - 24*day
	min := m - h*60
	sec := s - m*60
	ms := diffMS - s*1000
	var buff strings.Builder
	if isNegative {
		buff.WriteString("-")
	}
	if day > 0 {
		buff.WriteString(fmt.Sprintf("%dd", day))
	}
	if hour > 0 {
		buff.WriteString(fmt.Sprintf("%dh", hour))
	}
	if min > 0 {
		buff.WriteString(fmt.Sprintf("%dm", min))
	}
	if sec > 0 {
		buff.WriteString(fmt.Sprintf("%ds", sec))
	}
	if ms > 0 {
		buff.WriteString(fmt.Sprintf("%dms", ms))
	}
	s2 := buff.String()
	if s2 == "" {
		return "0ms"
	}
	return s2
}

// 传入的是毫秒值
func GetDistanceTimeToCurrent(startTime int64) string {
	diff := CurrentTime() - startTime
	return GetDistanceTime(diff)
}
