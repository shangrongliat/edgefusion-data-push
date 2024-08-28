package common

import "time"

// Time2String  将时间戳转换为时间字符串
func Time2String(t int64) string {
	tm := time.Unix(t, 0)
	return tm.Format("2006-01-02 15:04:05")
}

// Time2String  将时间戳转换为时间字符串
func TimeUint2String(t uint64) string {
	tm := time.Unix(int64(t)/1000, 0)
	return tm.Format("2006-01-02 15:04:05")
}
