package utils

import (
	"strings"
	"time"
)

const layOut = "2006-01-02 15:04:05"
const layOutTime = "2006-01-02T15:04:05Z07:00"
const unit_G = "G"
const unit_f = 102
const prec = 2
const bitSize = 64

// Time2String  将时间戳转换为时间字符串
func Time2String(t int64) string {
	tm := time.Unix(t, 0)
	return tm.Format("2006-01-02 15:04:05")
}

// Str2Timestamp  将时间字符串转换为时间戳 （本地实际）
func Str2Timestamp(timeStr string) int64 {
	//loc, _ := time.LoadLocation("Local")      //获取时区
	theTime, _ := time.Parse(layOut, timeStr) //使用模板在对应时区转化为time.time类型
	sr := theTime.Unix()                      //转化为时间戳 类型是int64
	return sr
}

// Str2TimeBLayOutTime 将时间字符串转换为时间戳 （本地实际）
func Str2TimeBLayOutTime(timeStr string) int64 {
	var sr int64
	//loc, _ := time.LoadLocation("Local")      //获取时区
	if strings.Contains(timeStr, "T") {
		theTime, _ := time.Parse(layOutTime, timeStr) //使用模板在对应时区转化为time.time类型
		sr = theTime.Unix()                           //转化为时间戳 类型是int64
	} else {
		theTime, _ := time.Parse(layOut, timeStr) //使用模板在对应时区转化为time.time类型
		sr = theTime.Unix()                       //转化为时间戳 类型是int64
	}
	return sr
}

// StrRFC3339Time 将时间字符串转换为时间戳 （RFC3339）
func StrRFC3339Time(toFormatTime string) string {
	res_time, _ := time.Parse(layOutTime, toFormatTime)
	toFormatTime = res_time.Format(layOut)
	return toFormatTime
}
