package common

import "time"

func Time2String(t uint64) string {
	tm := time.Unix(int64(t), 0)
	return tm.Format("2006-01-02 15:04:05")
}

func TimeInt2String(t int64) string {
	tm := time.Unix(t, 0)
	return tm.Format("2006-01-02 15:04:05")
}
